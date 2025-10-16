import os

import torch
import torch.nn as nn
import numpy as np

from transformers import RobertaTokenizer, RobertaModel
from torch.optim import AdamW
from utils.constant import DEVICE, SAVE_DIR
from collections import deque
import torch.nn.functional as F
from torch.distributions import Categorical

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
    """ACTOR: Predicts performance ast_metric (IC, cycles)"""
    def __init__(self, encoder_hidden_size, dropout, init_weigh=None):
        super().__init__()
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

    def forward(self, state, action):
        out = self.fc1(state)
        out = self.relu(out)
        # debug()
        out = self.fc2(torch.cat([out,action], dim=1))
        out = self.relu(out)
        out = self.fc3(out)
        return out
    
    def calculate_loss(self, y_pred, y_true):
        loss = F.mse_loss(y_pred, y_true)
        return loss

class ActorModel(nn.Module):
    def __init__(self, num_extra_feats, n_states = 7, n_actions=2, model_name="microsoft/codebert-base", heads=5):
        '''
        """state: [ic, sigma_ic, cycle, sigma_cycle, certainty, cpu%, memory%]"""
        '''
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

        # ensemble: predict performance ast_metric state from code features
        dropout_list = list(np.random.choice(np.arange(0.1, 0.6, 0.1), size=heads, replace=False))
        log_values = np.random.uniform(-6, -1, size=heads * 2)
        unique_vals = np.unique(np.round(log_values, 6))[:heads]  # ensure unique
        init_weigh = list(10 ** unique_vals)        
        self.heads = nn.ModuleList([
            self.load_code_with_metrics_model("./"+SAVE_DIR+"/"+"model.pt", hidden_size, dropout_list[i], init_weigh = init_weigh[i])
            for i in range(heads)
        ])

        self.shared_fc = nn.Sequential(
            nn.Linear(n_states, 128),
            nn.ReLU(),
            nn.Linear(128, 64),
            nn.ReLU(),
            nn.Linear(64, n_actions)
        )
        
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

        head_outputs = []
        for head in self.heads:
            ic, sigma_ic, cycle, sigma_cycle = head(x)
            head_outputs.append((ic, sigma_ic, cycle, sigma_cycle))

        self.head_outputs = head_outputs
        self.code_features = x  # Save for later use

        ics, sigmas_ic, cycles, sigmas_cycle = self.get_avg_head_outputs(head_outputs)
        certainty = self._calculate_certainty_per_sample(head_outputs)
        m_ic = torch.tensor(np.log10(master_predictions["ic"]))
        m_cycle = torch.tensor(np.log10(master_predictions["cycle"]))
        blended_ic = certainty * ics + (1 - certainty) * m_ic
        blended_cycle = certainty * cycles + (1 - certainty) * m_cycle
        
        # Build state with uncertainty information
        state = torch.cat([
            blended_ic.unsqueeze(0),
            blended_cycle.unsqueeze(0),
            certainty.unsqueeze(0),
            (ics - m_ic).abs().unsqueeze(0),  # Disagreement with base model
            (cycles - m_cycle).abs().unsqueeze(0),
            state["memory"].unsqueeze(0).unsqueeze(0),
            state["cpu"].unsqueeze(0).unsqueeze(0),
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

    def _calculate_certainty_per_sample(self, head_outputs=None):
        if head_outputs is None:
            head_outputs = self.head_outputs
        
        ics = torch.stack([h[0] for h in head_outputs], dim=0)  # [num_heads, batch]
        cycles = torch.stack([h[2] for h in head_outputs], dim=0)
        sigma_ics = torch.stack([h[1] for h in head_outputs], dim=0)
        sigma_cycles = torch.stack([h[3] for h in head_outputs], dim=0)

        mean_ic = ics.mean(dim=0)
        mean_cycle = cycles.mean(dim=0)
        
        epistemic_ic = ics.std(dim=0) / (mean_ic.abs() + 1e-6)
        epistemic_cycle = cycles.std(dim=0) / (mean_cycle.abs() + 1e-6)

        avg_sigma_ic = sigma_ics.mean(dim=0)
        avg_sigma_cycle = sigma_cycles.mean(dim=0)

        if not hasattr(self, 'uncertainty_weights'):
            self.uncertainty_weights = nn.Parameter(torch.ones(4))

        # Normalize uncertainties to similar scales
        epistemic = (epistemic_ic + epistemic_cycle) / 2
        aleatoric = (avg_sigma_ic + avg_sigma_cycle) / 2

        # Apply learned weights with softmax
        weights = F.softmax(self.uncertainty_weights, dim=0)
        total_uncertainty = (
            weights[0] * epistemic + 
            weights[1] * aleatoric
        )

        # Convert to certainty using exponential decay
        certainty = torch.exp(-weights[2] * total_uncertainty)

        # Apply temperature scaling
        temperature = torch.clamp(weights[3] * 10, 0.1, 10.0)
        certainty = torch.sigmoid((certainty - 0.5) * temperature)

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
        if len(self.memory) < batch_size:
            return
        indices = np.random.choice(len(self.memory), size=batch_size, replace=False)
        batch = [self.memory[idx] for idx in indices]
        return batch