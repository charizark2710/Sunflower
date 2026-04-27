import os

import torch
import torch.nn as nn
import numpy as np
import random

from transformers import RobertaTokenizer, RobertaModel
from torch.optim import AdamW
from utils.constant import DEVICE, SAVE_DIR
from collections import deque
import torch.nn.functional as F

LR = 1e-5

def init_w(m):
    m.weight.data.normal_(1.0, 0.02)
    m.bias.data.fill_(0)

class RelativeErrorWithSigmaLoss(nn.Module):
    def __init__(self, cycle_rate=3.5e9, ic_w=0.2, cycle_w=0.1, cpu_time_w=0.4, eps=1e-2):
        super().__init__()
        self.cycle_rate = np.log10(cycle_rate)
        self.ic_w = ic_w
        self.cycle_w = cycle_w
        self.cpu_time_w = cpu_time_w
        self.eps = eps

    def forward(self, ic_pred, sigma_ic,
               cycle_pred, sigma_cycle, ic_target, cycle_target):

        # Avoid divide by zero
        ic_pred = ic_pred.clamp(min=self.eps)
        ic_target = ic_target.clamp(min=self.eps)
        sigma_ic = sigma_ic.clamp(min=self.eps)
        sigma_cycle = sigma_cycle.clamp(min=self.eps)

        # log-likelihood style losses

        ic_loss  = 0.5 * (((ic_pred - ic_target) / sigma_ic)**2 
                          + 2 * torch.log(sigma_ic))
        cycle_loss = 0.5 * (((cycle_pred - cycle_target) / sigma_cycle)**2 
                            + 2 * torch.log(sigma_cycle))
        cpu_time_target = cycle_target / self.cycle_rate
        cpu_time_pred = (ic_pred * (cycle_pred / ic_pred)) / self.cycle_rate
        cpu_time_loss = 0.5 * (((cpu_time_pred - cpu_time_target) / sigma_cycle)**2
                            + 2 * torch.log(sigma_cycle))

        return (self.ic_w * ic_loss.mean().abs() +
                self.cycle_w * cycle_loss.mean().abs() +
                self.cpu_time_w * cpu_time_loss.mean().abs())
class DuelingHeadNet(nn.Module):
    def __init__(self, n_actions=2):
        super(DuelingHeadNet, self).__init__()
        mult = 64*7*7
        self.split_size = 512
        self.fc1 = nn.Linear(mult, self.split_size*2)
        self.value = nn.Linear(self.split_size, 1)
        self.advantage = nn.Linear(self.split_size, n_actions)
        self.fc1.apply(init_w)
        self.value.apply(init_w)
        self.advantage.apply(init_w)

    def forward(self, x):
        x1,x2 = torch.split(F.relu(self.fc1(x)), self.split_size, dim=1)
        value = self.value(x1)
        advantage = self.advantage(x2)
        # value is shape [batch_size, 1]
        # advantage is shape [batch_size, n_actions]
        q = value + torch.sub(advantage, torch.mean(advantage, dim=1, keepdim=True))
        return q

class CodeWithMetricsModel(nn.Module):
    """ACTOR: Predicts performance ast_metric (IC, cycles)"""
    def __init__(self, encoder_hidden_size):
        super().__init__()
        self.fc = nn.Sequential(
            nn.Linear(encoder_hidden_size + 64, 512), # Increased size
            nn.ReLU(),
            nn.Dropout(0.3),
            nn.Linear(512, 256),             # Added a new layer
            nn.ReLU(),
            nn.Dropout(0.3),
            nn.Linear(256, 128),
            nn.ReLU(),
            nn.Dropout(0.3),
            nn.Linear(128, 4),
        )
        
        self.optimizer = AdamW(self.parameters(), lr=LR, weight_decay=0.01)
        self.loss_fn = RelativeErrorWithSigmaLoss()

    def loss_function(self, ic_target, cycle_target):
        """Actor loss: prediction error"""
        ic_pred, sigma_ic, cycle_pred, sigma_cycle = self.result
        loss_val = self.loss_fn(torch.log10(ic_pred).squeeze(), sigma_ic.squeeze(), torch.log10(cycle_pred).squeeze(), sigma_cycle.squeeze(), ic_target.squeeze(), cycle_target.squeeze())
        loss_val.backward(retain_graph=True)
        torch.nn.utils.clip_grad_norm_(self.parameters(), max_norm=1.0)
        self.optimizer.step()
        self.optimizer.zero_grad()
        return loss_val.item()

    def forward(self, x):
        output = self.fc(x)
        ic = output[:, 0]
        sigma_ic = torch.exp(output[:, 1]) + 1e-6
        cycle = output[:, 2]
        sigma_cycle = torch.exp(output[:, 3]) + 1e-6
        
        self.result = (ic, sigma_ic, cycle, sigma_cycle)
        return ic, sigma_ic, cycle, sigma_cycle

class CriticModel(nn.Module):
    def __init__(self, nb_states=8, nb_actions=2):
        super().__init__()
        self.fc1 = nn.Linear(nb_states, 400)
        self.fc2 = nn.Linear(400 + nb_actions, 300)
        self.fc_certainty = nn.Linear(300 + 1, 150)
        self.fc3 = nn.Linear(150, 1)
        self.relu = nn.ReLU()
        
        # optional stabilization
        nn.init.uniform_(self.fc3.weight, -3e-3, 3e-3)
        nn.init.uniform_(self.fc3.bias, -3e-3, 3e-3)

    def forward(self, state, action, certainty):
        x = self.relu(self.fc1(state))
        x = self.relu(self.fc2(torch.cat([x, action], dim=1)))
        x = self.relu(self.fc_certainty(torch.cat([x, certainty.view(1, 1)], dim=1)))
        q_value = self.fc3(x)
        return q_value

    def calculate_loss(self, y_pred, y_true):
        loss = F.mse_loss(y_pred, y_true)
        return loss

class ActorModel(nn.Module):
    def __init__(self, num_extra_feats, n_states = 8, n_actions=2, model_name="microsoft/codebert-base", heads=5):
        '''
        """state: [ic, sigma_ic, cycle, sigma_cycle, m_ic, m_cycle, cpu%, memory%]"""
        '''
        super().__init__()
        
        # Shared encoder for extracting code features
        self.encoder = RobertaModel.from_pretrained(model_name)
        hidden_size = self.encoder.config.hidden_size
        self.norm = nn.LayerNorm(hidden_size + 64)
        
        self.n_actions = n_actions
        self.feat_proj = nn.Sequential(
            nn.Linear(num_extra_feats, 64),
            nn.ReLU(),
            nn.Linear(64, 64)
        )

        self.core = self.load_code_with_metrics_model("./"+SAVE_DIR+"/"+"model.pt", hidden_size)
        self.heads = nn.ModuleList([DuelingHeadNet(n_actions=n_actions) for k in range(heads)])

        self.shared_fc = nn.Sequential(
            nn.Linear(n_states, 128),
            nn.ReLU(),
            nn.Linear(128, 64),
            nn.ReLU(),
            nn.Linear(64, 64*7*7)
        )
        
        self.mix_logit = nn.Parameter(torch.tensor(0.0))
        self.scale_param = nn.Parameter(torch.tensor(0.0))
        self.temp_param = nn.Parameter(torch.tensor(-2.0))

    def extract_code_embeddings(self, input_ids_list, attention_mask_list, ast_metric):
        # Extract code embeddings
        flat_input_ids = torch.cat(input_ids_list, dim=0)
        flat_attention_mask = torch.cat(attention_mask_list, dim=0)

        outputs = self.encoder(flat_input_ids, attention_mask=flat_attention_mask)
        chunk_reps = outputs.last_hidden_state.mean(dim=1)

        # Pool back per sample
        pooled_outputs, idx = [], 0
        for ids in input_ids_list:
            n_chunks = ids.size(0)
            pooled = chunk_reps[idx:idx + n_chunks].mean(dim=0)
            pooled_outputs.append(pooled)
            idx += n_chunks

        code_emb = torch.stack(pooled_outputs, dim=0)
        feat_proj = self.feat_proj(ast_metric)
        x = torch.cat([code_emb, feat_proj], dim=1)
        x = self.norm(x)
        
        return x

    def forward(self, input_ids_list, attention_mask_list, state, ast_metric, master_predictions):
        """
        Returns: value estimate from critic
        """
        x = self.extract_code_embeddings(input_ids_list, attention_mask_list, ast_metric)
        x = self.core(x)
        ic, sigma_ic, cycle, sigma_cycle = x
        m_ic = torch.tensor(np.log10(master_predictions["ic"]))
        m_cycle = torch.tensor(np.log10(master_predictions["cycle"]))

        # print(f"torch.log10(ic).squeeze(0).unsqueeze(0).shape: {torch.log10(ic).squeeze(0).unsqueeze(0).shape}")
        # print(f"torch.log10(cycle).squeeze(0).unsqueeze(0).shape: {torch.log10(cycle).squeeze(0).unsqueeze(0).shape}")
        # print(f"sigma_ic.squeeze(0).unsqueeze(0).shape: {sigma_ic.squeeze(0).unsqueeze(0).shape}")
        # print(f"sigma_cycle.squeeze(0).unsqueeze(0).shape: {sigma_cycle.squeeze(0).unsqueeze(0).shape}")
        # print(f"m_ic.unsqueeze(0).shape: {m_ic.unsqueeze(0).shape}")
        # print(f"m_cycle.unsqueeze(0).shape: {m_cycle.unsqueeze(0).shape}")
        # print(f"state['memory'].unsqueeze(0).shape: {state['memory'].unsqueeze(0).shape}")
        # print(f"state['cpu'].unsqueeze(0).shape: {state['cpu'].unsqueeze(0).shape}")
        state = torch.cat([
            torch.log10(ic).reshape(1, 1).float(),
            torch.log10(cycle).reshape(1, 1).float(),
            sigma_ic.reshape(1, 1).float(),
            sigma_cycle.reshape(1, 1).float(),
            m_ic.reshape(1, 1).float(),
            m_cycle.reshape(1, 1).float(),
            state["memory"].reshape(1, 1).float(),
            state["cpu"].reshape(1, 1).float(),
        ], dim=1)

        x = self.shared_fc(state)

        self.head_outputs = [net(x) for net in self.heads]
        self.code_features = x  # Save for later use

        certainty = self._calculate_certainty_per_sample()
        
        avg_head_output = torch.mean(torch.stack(self.head_outputs), dim=0)
        return avg_head_output, state, certainty

    def load_code_with_metrics_model(self, checkpoint_path, encoder_hidden_size=768,
                                    ):
        """Load a CodeWithMetricsModel from checkpoint"""
        # Determine mode
        model = CodeWithMetricsModel(
            encoder_hidden_size=encoder_hidden_size,
        )
        if not os.path.exists(checkpoint_path):
            print(f"✗ Checkpoint not found: {checkpoint_path}")
            return model
        checkpoint = torch.load(checkpoint_path, map_location=DEVICE)

        if isinstance(checkpoint, dict) and 'model_state_dict' in checkpoint:
            state_dict = checkpoint['model_state_dict']
        else:
            state_dict = checkpoint
        
        # Load weights (strict=False allows loading head weights without encoder)
        missing, unexpected = model.load_state_dict(state_dict, strict=False)

        if missing:
            print(f"Missing keys: {missing}")
        if unexpected:
            print(f"Unexpected keys: {unexpected}")

        model = model.to(DEVICE)
        model.eval()

        print(f"✓ Model loaded successfully")
        return model

    def _calculate_certainty_per_sample(self):
        head_outputs = self.head_outputs
        advantages = [h[0] - h[0].mean(dim=0, keepdim=True) for h in head_outputs]
        var_q = torch.var(torch.stack(advantages), dim=0)
        uncertainty = var_q.mean(dim=-1)
        certainty = torch.exp(-uncertainty)
        return certainty

    def loss_fn(self, q_values, actions, rewards):
        """
        Contextual bandit loss for multi-head critic.
        
        q_values: list of [batch_size, num_actions] from each head
        actions: [batch_size]  (action indices taken)
        rewards: [batch_size]  (observed rewards)
        """
        q_pred = q_values.gather(1, actions.unsqueeze(1)).squeeze(1)
        # MSE between predicted Q and observed reward
        loss = F.mse_loss(q_pred, rewards)
        return loss

class ReplayMemory():
    def __init__(self, maxlen=1000):
        self.memory = deque([], maxlen=maxlen)
    def append(self, transition: dict):
        self.memory.append(transition)
    def __len__(self):
        return len(self.memory)
    
    def clear(self):
        self.memory.clear()
    
    def sample(self, batch_size):
        
        if len(self.memory) < batch_size:
            batch_size = len(self.memory)    
        return random.sample(self.memory, batch_size)