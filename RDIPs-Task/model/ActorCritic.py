
from base_model.model import ReplayMemory, ActorModel, CriticModel
from utils.constant import DEVICE
from unix_socket import send_msg, recv_msg
from base_model.env import TaskEvalEnv
from utils.common import load_model, profile_with_psutil
from utils.constant import SAVE_DIR

import pandas as pd
import numpy as np
import torch
import random
import os

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
        self.epsilon = 0
        self.epsilon_history = []
        self.device = DEVICE
        actor, critic = load_model(num_extra_feats)
        self.actor = actor.to(self.device)
        self.critic = critic.to(self.device)

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
    def select_action(self, code_tokens, attention_mask, state, ast_metric, master_pred):
        actions, pred_state, certainty = self.actor(code_tokens, attention_mask, state, ast_metric, master_pred)
        action = torch.argmax(actions, dim=-1).item()
        return action, pred_state, certainty

    # -------------------------------------------------------------------------
    # --- TRAINING LOOP
    # -------------------------------------------------------------------------
    def train(self, episodes):
        epsilon_decay = 0.995
        for _, episode in enumerate(episodes):
            ast_metric = episode["ast_metric"]
            code_tokens = episode["code_tokens"]
            attention_mask = episode["attention_mask"]
            master_pred = episode["master_pred"]
            state = torch.tensor(self.env.reset(), dtype=torch.float, device=self.device)
            certainty = torch.tensor(0.0, device=self.device)

            if random.random() < self.epsilon:
                action = self.env.action_space.sample()
            else:
                with torch.no_grad():
                    action, _, certainty = self.select_action(
                        code_tokens, attention_mask,
                        {"cpu": state[1], "memory": state[0]}, ast_metric, master_pred
                    )

            mem_usage = 0.0
            cpu_usage = 0.0

            result = None
            # execute
            if action == 1:
                send_and_recv = profile_with_psutil(lambda conn, msg: (send_msg(conn, msg), recv_msg(conn))[1])

                result, mem_usage, cpu_usage = send_and_recv(self.conn, {
                    "isSkip": False, "confidence": float(certainty.item()), "done": True,
                })
                if result is not None and result.get("actual_ic") is not None and result.get("actual_cycle") is not None:
                    result["actual_ic"] = torch.log10(torch.as_tensor(result["actual_ic"], dtype=torch.float32, device=self.device))
                    result["actual_cycle"] = torch.log10(torch.as_tensor(result["actual_cycle"], dtype=torch.float32, device=self.device))
            else:
                send_msg(self.conn, {"isSkip": True, "done": True})

            # update env
            self.env.update_context(result, {"certainty": certainty.item()})
            next_state, reward, _ = self.env.step(action, mem_usage, cpu_usage)

            next_state = torch.FloatTensor(next_state).to(self.device)

            self.memory.append({
                'state': state,
                'action': action,
                'reward': reward,
                'next_state': next_state,
                'certainty': certainty.item(),
                "code_tokens": code_tokens,
                "attention_mask": attention_mask,
                "ast_metric": ast_metric,
                "actual_cycle": result["actual_cycle"] if result is not None and result.get("actual_cycle") else None,
                "actual_ic": result["actual_ic"] if result is not None and result.get("actual_ic") else None,
                "master_pred": master_pred
            })

            state = next_state

            self.epsilon_history.append(reward)
            self.epsilon = max(0.1, self.epsilon * epsilon_decay)
        return self.update_model()

    # -------------------------------------------------------------------------
    # --- MODEL UPDATE
    # -------------------------------------------------------------------------
    def update_model(self):
        print("Updating model...")
        batch = self.memory.sample(self.batch_size)
        if batch is None:
            return

        total_actor_loss = 0.0
        total_critic_loss = 0.0

        for exp in batch:
            state = exp['state'].unsqueeze(0).to(self.device)
            reward = torch.tensor([[exp['reward']]], device=self.device, dtype=torch.float)
            certainty = torch.tensor([[exp['certainty']]], device=self.device)

            code_tokens = exp['code_tokens']
            attention_mask = exp['attention_mask']
            ast_metric = exp['ast_metric']
            master_pred = exp["master_pred"]
            action = exp['action']
            target_value = reward

            dist, state_feats, pred_certainty = self.actor(
                code_tokens, attention_mask,
                {"cpu": state[0, 1], "memory": state[0, 0]},
                ast_metric, master_pred
            )
            current_q = self.critic(state_feats, dist, pred_certainty)

            critic_loss =  self.critic.calculate_loss(current_q, target_value)

            advantage = reward - current_q.detach()
            if not isinstance(action, torch.Tensor):
                action_tensor = torch.tensor([action], device=self.device)
            else:
                action_tensor = action
            log_probs_all = torch.nn.functional.log_softmax(dist, dim=-1)
            log_prob = log_probs_all.gather(-1, action_tensor.unsqueeze(-1)).squeeze(-1)
            actor_loss = -log_prob * advantage
            total_loss = actor_loss + critic_loss
            print(f"Action: {action}, Intermediate Actor Loss: {actor_loss.item()}, Critic Loss: {critic_loss.item()}, Total Loss: {total_loss.item()}, certainty: {certainty.item()}, reward: {reward.item()}")
            all_trainable_parameters = list(self.actor.parameters()) + list(self.critic.parameters())
            torch.nn.utils.clip_grad_norm_(all_trainable_parameters, max_norm=1.0)
            self.actor_optimizer.zero_grad()
            self.critic_optimizer.zero_grad()
            total_loss.backward()
            self.critic_optimizer.step()
            self.actor_optimizer.step()
            
            total_actor_loss += actor_loss.item()
            total_critic_loss += critic_loss.item()

            data = pd.DataFrame({
                "actor_loss": [total_actor_loss / self.batch_size],
                "critic_loss": [total_critic_loss / self.batch_size],
                "total_loss": [(total_actor_loss + total_critic_loss) / self.batch_size]
            })
            if not os.path.exists('./loss.csv'):
                data.to_csv('./loss.csv', index=False)
            else:
                data.to_csv('./loss.csv', mode='a', header=False, index=False)
            
        torch.save(self.actor.state_dict(), "./"+SAVE_DIR + "/actor_model.pt")
        torch.save(self.critic.state_dict(), "./"+SAVE_DIR + "/critic_model.pt")
        print(f"Actor Loss: {total_actor_loss / self.batch_size}, Critic Loss: {total_critic_loss / self.batch_size}")
        return