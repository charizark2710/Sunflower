
from base_model.model import ReplayMemory, ActorModel, CriticModel
from utils.constant import DEVICE
from unix_socket import send_msg, recv_msg
from base_model.env import TaskEvalEnv
from utils.common import load_model
from utils.constant import SAVE_DIR

import torch.nn.functional as F
import numpy as np
import torch

import torch
import torch.nn.functional as F
import numpy as np
class ActorCritic:
    def __init__(self,
                 num_extra_feats: int,
                 conn,
                 gamma=0.99,
                 tau=0.005,
                 actor_lr=1e-4,
                 critic_lr=1e-3,
                 batch_size=5,
                 memory_capacity=1000):

        self.device = DEVICE
        actor, critic = load_model(num_extra_feats)
        self.actor = actor.to(self.device)
        self.critic = critic.to(self.device)
        self.target_actor = ActorModel(num_extra_feats).to(self.device)
        self.target_critic = CriticModel().to(self.device)

        self.target_actor.load_state_dict(self.actor.state_dict())
        self.target_critic.load_state_dict(self.critic.state_dict())

        self.conn = conn
        self.gamma = gamma
        self.tau = tau
        self.batch_size = batch_size

        self.actor_optimizer = torch.optim.Adam(self.actor.parameters(), lr=actor_lr)
        self.critic_optimizer = torch.optim.Adam(self.critic.parameters(), lr=critic_lr)

        self.memory = ReplayMemory(memory_capacity)
        self.env = TaskEvalEnv()

    def updateConn(self, conn):
        self.conn = conn

    # -------------------------------------------------------------------------
    # --- ACTOR PREDICTION / ACTION SELECTION
    # -------------------------------------------------------------------------
    @torch.no_grad()
    def select_action(self, code_tokens, attention_mask, state, ast_metric, master_pred):
        self.actor.eval()
        dist, pred_state = self.actor(code_tokens, attention_mask, state, ast_metric, master_pred)
        action = torch.argmax(dist.probs, dim=-1).item()
        return action, pred_state

    # -------------------------------------------------------------------------
    # --- TRAINING LOOP
    # -------------------------------------------------------------------------
    def train(self, episodes):
        rewards_per_episode = np.zeros(len(episodes))

        for ep_idx, episode in enumerate(episodes):
            ast_metric = episode["ast_metric"]
            code_tokens = episode["code_tokens"]
            attention_mask = episode["attention_mask"]
            master_pred = episode["master_pred"]
            state = torch.tensor(self.env.reset(), dtype=torch.float, device=self.device)
            episode_reward = 0.0
            done = False
            certainty = torch.tensor(0.0, device=self.device)
            inner_step = 0

            while not done:
                inner_step += 1
                if inner_step >= 5:
                    done = True

                action, pred_state = self.select_action(
                    code_tokens, attention_mask,
                    {"cpu": state[1], "memory": state[0]}, ast_metric, master_pred
                )
                
                pred_ic = pred_state[:, 0:1]         # blended_ic
                pred_cycle = pred_state[:, 1:2]      # blended_cycle
                certainty = pred_state[:, 2:3]       # certainty

                head_loss = 0.0
                result = None
                # execute
                if action == 1:
                    send_msg(self.conn, {"isSkip": False, "confidence": float(certainty.item()), "done": done})
                    result = recv_msg(self.conn)
                    if result is not None:
                        head_loss = self.actor.train_head(result["ic"], result["cycle"])
                    else:
                        head_loss = 10.0
                else:
                    send_msg(self.conn, {"isSkip": True, "done": done})

                # update env
                self.env.update_context(result, (certainty.item(), pred_ic.item(), pred_cycle.item()))
                next_state, reward, terminated = self.env.step(action)
                next_state = torch.FloatTensor(next_state).to(self.device)

                self.memory.append({
                    'state': state,
                    'action': action,
                    'reward': reward,
                    'next_state': next_state,
                    'head_loss': head_loss,
                    'certainty': certainty.item(),
                    "code_tokens": code_tokens,
                    "attention_mask": attention_mask,
                    "ast_metric": ast_metric,
                    "actual_cycle": result["cycle"] if result else None,
                    "actual_ic": result["ic"] if result else None,
                    "master_pred": master_pred
                })

                state = next_state
                episode_reward += reward

            rewards_per_episode[ep_idx] = episode_reward

        return self.update_model()

    # -------------------------------------------------------------------------
    # --- MODEL UPDATE
    # -------------------------------------------------------------------------
    def update_model(self):
        batch = self.memory.sample(self.batch_size)
        if batch is None:
            return
        total_actor_loss = 0.0
        total_critic_loss = 0.0

        for exp in batch:
            state = exp['state'].unsqueeze(0).to(self.device)
            next_state = exp['next_state'].unsqueeze(0).to(self.device)
            reward = torch.tensor([[exp['reward']]], device=self.device, dtype=torch.float)
            certainty = torch.tensor([[exp['certainty']]], device=self.device)
            head_loss = torch.tensor([[1.0 / (1 + float(exp.get('head_loss', 10.0)))]], device=self.device)

            code_tokens = exp['code_tokens']
            attention_mask = exp['attention_mask']
            ast_metric = exp['ast_metric']

            master_pred = exp["master_pred"]
            # Critic update
            with torch.no_grad():
                dist_next, next_state = self.target_actor(code_tokens, attention_mask,
                                                 {"cpu": next_state[0, 1], "memory": next_state[0, 0]}, ast_metric, master_pred)
                target_q = self.target_critic(next_state, dist_next.probs)
                target_value = reward + self.gamma * target_q * head_loss

            dist, state = self.actor(code_tokens, attention_mask,
                                 {"cpu": state[0, 1], "memory": state[0, 0]}, ast_metric, master_pred)
            current_q = self.critic(state, dist.probs)

            critic_loss = F.mse_loss(current_q, target_value)
            self.critic_optimizer.zero_grad()
            critic_loss.backward()
            torch.nn.utils.clip_grad_norm_(self.critic.parameters(), 1.0)
            self.critic_optimizer.step()
            total_critic_loss += critic_loss.item()

            # Actor update

            self.actor_optimizer.zero_grad()
            state_detached = state.detach()

            dist_actor, _ = self.actor(
                code_tokens, attention_mask,
                {"cpu": state_detached[0, 1], "memory": state_detached[0, 0]},
                ast_metric, master_pred
            )
            actor_loss = -self.critic(state_detached, dist_actor.probs)
            actor_loss = (actor_loss * head_loss * certainty).mean()
            actor_loss.backward()
            torch.nn.utils.clip_grad_norm_(self.actor.parameters(), 1.0)
            self.actor_optimizer.step()
            total_actor_loss += actor_loss.item()

        # Soft update targets
        with torch.no_grad():
            for target_param, param in zip(self.target_critic.parameters(), self.critic.parameters()):
                target_param.data.copy_(param.data)
            for target_param, param in zip(self.target_actor.parameters(), self.actor.parameters()):
                target_param.data.copy_(param.data)

        # Save actor and critic model
        torch.save(self.actor.state_dict(), "./"+SAVE_DIR + "/actor_model.pt")
        torch.save(self.critic.state_dict(), "./"+SAVE_DIR + "/critic_model.pt")

        return
