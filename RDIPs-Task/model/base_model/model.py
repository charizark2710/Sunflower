import os

import torch
import torch.nn as nn
import numpy as np

from transformers import RobertaTokenizer, RobertaModel
from torch.optim import AdamW
from utils.common import CPU_RATE, DEVICE
from collections import deque
import torch.nn.functional as F
from torch.distributions import Categorical

SAVE_DIR = "./saved_model"
LR = 1e-4

def fanin_init(size, fanin=None):
    fanin = fanin or size[0]
    v = 1. / np.sqrt(fanin)
    return torch.Tensor(size).uniform_(-v, v)

class RelativeErrorWithSigmaLoss(nn.Module):
    def __init__(self, cycle_rate=3.5e9, ic_w=0.2, cycle_w=0.1, cpu_time_w=0.4, eps=1e-2):
        super().__init__()
        self.cycle_rate = cycle_rate
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
        # cpu_time_loss = 0.5 * (((cpu_time_pred - cpu_time_target) / sigma_cpu)**2 
        #                        + 2 * torch.log(sigma_cpu))
        ic_loss  = 0.5 * (((ic_pred - ic_target) / sigma_ic)**2 
                          + 2 * torch.log(sigma_ic))
        cycle_loss = 0.5 * (((cycle_pred - cycle_target) / sigma_cycle)**2 
                            + 2 * torch.log(sigma_cycle))
        cpu_time_loss = cycle_loss / self.cycle_rate

        return (self.ic_w * ic_loss.mean().abs() +
                self.cycle_w * cycle_loss.mean().abs() +
                self.cpu_time_w * cpu_time_loss.mean().abs())

class CodeWithMetricsModel(nn.Module):
    """ACTOR: Predicts performance ast_metrics (IC, cycles)"""
    def __init__(self, encoder_hidden_size, dropout, init_weigh=None):
        super().__init__(encoder_hidden_size, dropout, init_weigh)

        self.fc = nn.Sequential(
            nn.Linear(encoder_hidden_size + 64, 512),
            nn.ReLU(),
            nn.Dropout(dropout),
            nn.Linear(512, 256),
            nn.ReLU(),
            nn.Dropout(dropout),
            nn.Linear(256, 128),
            nn.ReLU(),
            nn.Dropout(dropout),
            nn.Linear(128, 4),  # [ic, log_sigma_ic, cycle, log_sigma_cycle]
        )
        
        self.optimizer = AdamW(self.parameters(), lr=LR, weight_decay=0.01)
        self.loss_fn = RelativeErrorWithSigmaLoss()
        self.init_weighs(init_weigh)

    def init_weighs(self, init_w):
        # Get all Linear layers from the Sequential
        linear_layers = [module for module in self.fc.modules() if isinstance(module, nn.Linear)]
        
        # Initialize hidden layers with fanin
        for layer in linear_layers[:-1]:
            layer.weight.data = fanin_init(layer.weight.data.size())
        
        # Initialize output layer with uniform distribution
        linear_layers[-1].weight.data.uniform_(-init_w, init_w)

    def loss_function(self, ic_target, cycle_target):
        """Actor loss: prediction error"""
        ic_pred, sigma_ic, cycle_pred, sigma_cycle = self.result
        loss_val = self.loss_fn(ic_pred, sigma_ic, cycle_pred, sigma_cycle, ic_target, cycle_target)
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
    def __init__(self, nb_states = 7, nb_actions = 2):
        """CRITIC: Evaluates quality of actor's predictions"""
        """action: [1, 0]"""
        """state: [ic, sigma_ic, cycle, sigma_cycle, certainty, cpu%, memory%]"""
        super().__init__()
        self.fc1 = nn.Linear(nb_states, 400)
        self.fc2 = nn.Linear(400+nb_actions, 300)
        self.fc3 = nn.Linear(300, 1)
        self.relu = nn.ReLU()

    def forward(self, xs):
        state, action = xs
        out = self.fc1(state)
        out = self.relu(out)
        # debug()
        out = self.fc2(torch.cat([out,action],1))
        out = self.relu(out)
        out = self.fc3(out)
        return out
    
    def calculate_loss(self, y_pred, y_true):
        loss = F.mse_loss(y_pred, y_true)
        return loss

class ActorModel(nn.Module):
    def __init__(self, num_extra_feats, n_states = 7, n_actions=2, model_name="microsoft/codebert-base", heads=5):
        super().__init__()
        
        # Shared encoder for extracting code features
        self.encoder = RobertaModel.from_pretrained(model_name)
        hidden_size = self.encoder.config.hidden_size
        self.norm = nn.LayerNorm(hidden_size + 64)
        
        self.feat_proj = nn.Sequential(
            nn.Linear(num_extra_feats, 64),
            nn.ReLU(),
            nn.Linear(64, 64)
        )

        # ensemble: predict performance ast_metrics state from code features
        dropout_list = list(np.random.choice(np.arange(0.1, 0.6, 0.1), size=heads, replace=False))
        init_weigh = list(np.random.choice(np.arange(1e-1, 1e-6, 1e-1), size=heads, replace=False))
        
        self.heads = nn.ModuleList([
            self.load_code_with_metrics_model("./saved_model", hidden_size, dropout_list[i], init_weigh = init_weigh[i])
            for i in range(heads)
        ])

        self.shared_fc = nn.Sequential(
            nn.Linear(n_states, 128),
            nn.ReLU(),
            nn.Linear(128, 64),
            nn.ReLU(),
            nn.Linear(64, n_actions)
        )
        
    def extract_code_embeddings(self, input_ids_list, attention_mask_list, ast_metrics):
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
        feat_proj = self.feat_proj(ast_metrics)
        x = torch.cat([code_emb, feat_proj], dim=1)
        x = self.norm(x)
        
        return x

    def forward(self, input_ids_list, attention_mask_list, state, ast_metrics, master_predictions):
        """
        Returns: value estimate from critic
        """

        x = self.extract_code_embeddings(input_ids_list, attention_mask_list, ast_metrics)

        head_outputs = []
        for head in self.heads:
            ic, sigma_ic, cycle, sigma_cycle = head(x)
            head_outputs.append((ic, sigma_ic, cycle, sigma_cycle))

        self.head_outputs = head_outputs
        self.code_features = x  # Save for later use

        ics, sigmas_ic, cycles, sigmas_cycle = self.get_avg_head_outputs(head_outputs)
        certainty = self._calculate_certainty_per_sample(master_predictions, head_outputs)
        
        state = torch.cat([
            ics.unsqueeze(1), 
            sigmas_ic.unsqueeze(1),
            cycles.unsqueeze(1), 
            sigmas_cycle.unsqueeze(1),
            certainty.unsqueeze(1),  # [batch, 1]
            state["memory"].unsqueeze(1),
            state["cpu"].unsqueeze(1),
        ], dim=1)

        action = self.shared_fc(state)
        distribution = Categorical(F.softmax(action, dim=-1))
        return distribution, state
    def load_code_with_metrics_model(self, checkpoint_path, encoder_hidden_size=768, dropout=0.3,
                                    init_weigh=1e-3):
        """Load a CodeWithMetricsModel from checkpoint"""
        # Determine mode
        model = CodeWithMetricsModel(
            encoder_hidden_size=encoder_hidden_size,
            dropout=dropout,
            init_weigh=init_weigh
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

    def get_avg_head_outputs(self, head_outputs = None):
        if head_outputs is None:
            head_outputs = self.head_outputs
        """Get average actor predictions"""
        ics = torch.stack([h[0] for h in head_outputs], dim=0)
        sigmas_ic = torch.stack([h[1] for h in head_outputs], dim=0)
        cycles = torch.stack([h[2] for h in head_outputs], dim=0)
        sigmas_cycle = torch.stack([h[3] for h in head_outputs], dim=0)

        return ics.mean(dim=0), sigmas_ic.mean(dim=0), cycles.mean(dim=0), sigmas_cycle.mean(dim=0)

    def _calculate_certainty_per_sample(self, master_predictions, head_outputs=None):
        if head_outputs is None:
            head_outputs = self.head_outputs
        
        # Stack slave predictions
        ics = torch.stack([h[0] for h in head_outputs], dim=0)  # [num_heads, batch]
        cycles = torch.stack([h[2] for h in head_outputs], dim=0)
        sigma_ics = torch.stack([h[1] for h in head_outputs], dim=0)
        sigma_cycles = torch.stack([h[3] for h in head_outputs], dim=0)
        
        # --- 1. EPISTEMIC UNCERTAINTY (slave heads disagreement) ---
        var_ic = ics.var(dim=0, unbiased=False)  # [batch]
        var_cycle = cycles.var(dim=0, unbiased=False)
        
        # Normalize by magnitude to make comparable across samples
        epistemic_ic = var_ic / (ics.mean(dim=0).abs() + 1e-6)
        epistemic_cycle = var_cycle / (cycles.mean(dim=0).abs() + 1e-6)
        epistemic_uncertainty = (epistemic_ic + epistemic_cycle) / 2
        
        # --- 2. ALEATORIC UNCERTAINTY (average confidence from heads) ---
        # High sigma = low confidence
        avg_sigma_ic = sigma_ics.mean(dim=0)
        avg_sigma_cycle = sigma_cycles.mean(dim=0)
        
        # Normalize by prediction magnitude
        aleatoric_ic = avg_sigma_ic / (ics.mean(dim=0).abs() + 1e-6)
        aleatoric_cycle = avg_sigma_cycle / (cycles.mean(dim=0).abs() + 1e-6)
        aleatoric_uncertainty = (aleatoric_ic + aleatoric_cycle) / 2
        
        # --- 3. MASTER-SLAVE GAP (hardware compatibility) ---
        m_ic, m_sigma_ic, m_cycle, m_sigma_cycle = master_predictions
        
        # Average slave prediction
        slave_ic_avg = ics.mean(dim=0)
        slave_cycle_avg = cycles.mean(dim=0)
        
        # Relative disagreement (normalized by magnitude)
        ic_gap = torch.abs(m_ic - slave_ic_avg) / (torch.abs(m_ic) + torch.abs(slave_ic_avg) + 1e-6)
        cycle_gap = torch.abs(m_cycle - slave_cycle_avg) / (torch.abs(m_cycle) + torch.abs(slave_cycle_avg) + 1e-6)
        
        # Option: Weight gap by confidence intervals
        # If both master and slave are uncertain (high sigma), gap matters less
        ic_gap_weighted = ic_gap * torch.exp(-(m_sigma_ic + avg_sigma_ic) / 2)
        cycle_gap_weighted = cycle_gap * torch.exp(-(m_sigma_cycle + avg_sigma_cycle) / 2)

        hardware_gap = (ic_gap_weighted + cycle_gap_weighted) / 2
        
        # --- COMBINE ALL UNCERTAINTIES ---
        # Weights can be tuned based on importance
        w_epistemic = 0.3    # Slave ensemble disagreement
        w_aleatoric = 0.2    # Prediction confidence (sigma)
        w_hardware = 0.5     # Master-slave gap (most important for execution safety)
        
        total_uncertainty = (
            w_epistemic * epistemic_uncertainty +
            w_aleatoric * aleatoric_uncertainty +
            w_hardware * hardware_gap
        )
        
        # Normalize to [0, 1] range per batch
        normalized_uncertainty = total_uncertainty / (total_uncertainty.max() + 1e-8)
        
        # Convert to certainty: high uncertainty = low certainty
        certainty = 1.0 - torch.clamp(normalized_uncertainty, 0, 1)
        
        return certainty
    def train_head(self, actual_ic, actual_cycle):
        """Train all actor heads with ground truth"""
        total_loss = 0
        for head in self.heads:
            loss = head.loss_function(actual_ic, actual_cycle)
            total_loss += loss
        return total_loss / len(self.heads)

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
        indices = np.random.choice(len(self.memory), size=batch_size, replace=False)
        batch = [self.memory[idx] for idx in indices]
        return batch