import os

import torch
import torch.nn as nn
from transformers import RobertaTokenizer, RobertaModel
from torch.optim import AdamW
from utils.common import DEVICE

import numpy as np
SAVE_DIR = "./saved_model"
LR = 1e-4

class CodeWithMetricsModel(nn.Module):
    def __init__(self, model_name, num_extra_feats):
        super().__init__()
        self.encoder = RobertaModel.from_pretrained(model_name)

        self.feat_proj = nn.Sequential(
            nn.Linear(num_extra_feats, 64),  # Increased size
            nn.ReLU(),
            nn.Linear(64, 64)                # Increased size
        )

        self.norm = nn.LayerNorm(self.encoder.config.hidden_size + 64) # Adjust norm layer

        self.fc = nn.Sequential(
            nn.Linear(self.encoder.config.hidden_size + 64, 512), # Increased size
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

    def forward(self, input_ids_list, attention_mask_list, extra_feats):
        # Flatten chunks across batch
        flat_input_ids = torch.cat(input_ids_list, dim=0)        # [total_chunks, seq_len]
        flat_attention_mask = torch.cat(attention_mask_list, dim=0)

        outputs = self.encoder(flat_input_ids, attention_mask=flat_attention_mask)
        chunk_reps = outputs.last_hidden_state.mean(dim=1)       # [total_chunks, hidden]

        # Now regroup per sample
        pooled_outputs = []
        idx = 0
        for ids in input_ids_list:
            n_chunks = ids.size(0)
            pooled = chunk_reps[idx: idx + n_chunks].mean(dim=0)
            pooled_outputs.append(pooled)
            idx += n_chunks

        pooled_output = torch.stack(pooled_outputs, dim=0)       # [batch, hidden]
        feat_proj = self.feat_proj(extra_feats)
        x = torch.cat([pooled_output, feat_proj], dim=1)
        x = self.norm(x)
        output = self.fc(x)
        ic = output[:, 0]
        sigma_ic = output[:, 1]
        cycle = output[:, 2]
        sigma_cycle = output[:, 3]
        sigma_ic  = torch.exp(sigma_ic) + 1e-6
        sigma_cycle = torch.exp(sigma_cycle) + 1e-6

        return ic, sigma_ic, cycle, sigma_cycle

# class RelativeErrorWithSigmaLoss(nn.Module):
    # def __init__(self, sigma_reg_weight=0.1, min_sigma_ratio=0.01):
    #     super().__init__()
    #     self.sigma_reg_weight = sigma_reg_weight
    #     self.min_sigma_ratio = min_sigma_ratio
    
    # def forward(self, log10_pred, log10_sigma, log10_target):
    #     # Convert to linear space
    #     linear_pred = 10 ** log10_pred
    #     print("linear_pred", linear_pred.mean())
    #     linear_target = 10 ** log10_target
    #     print("log10_target", log10_target.mean())
    #     # Convert sigma to linear space (as standard deviation)
    #     linear_sigma_upper = 10 ** (log10_pred + log10_sigma)
    #     linear_sigma = linear_sigma_upper - linear_pred
        
    #     # Ensure reasonable sigma bounds
    #     min_sigma = self.min_sigma_ratio * linear_pred
    #     linear_sigma = torch.clamp(linear_sigma, min=min_sigma)
    #     print("linear_sigma", linear_sigma.mean())
    #     # Relative error
    #     relative_error = torch.abs(linear_pred - linear_target) / (linear_target + 1e-6)
    #     print("relative_error", relative_error.mean())
    #     # Weight the error by uncertainty (inverse weighting)
    #     # High sigma → low weight → lower penalty
    #     # Low sigma → high weight → higher penalty
    #     sigma_ratio = linear_sigma / linear_pred  # Coefficient of variation
    #     print("sigma_ratio", sigma_ratio.mean())
    #     weighted_error = relative_error / (sigma_ratio + 1e-6)
        
    #     # Regularize sigma to prevent it from becoming too large
    #     sigma_penalty = self.sigma_reg_weight * sigma_ratio.mean()
    #     print("sigma_penalty", sigma_penalty.mean())
    #     print("tessst:", weighted_error.mean() + sigma_penalty)
    #     return weighted_error.mean() + sigma_penalty

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

def load_model(model_name):
    if os.path.exists(SAVE_DIR + "/model.pt"):
        model = CodeWithMetricsModel(model_name, 7)
        model.load_state_dict(torch.load(SAVE_DIR + "/model.pt", map_location=DEVICE))
        optimizer = AdamW(model.parameters(), lr=LR, weight_decay=0.05)
        # Add learning rate scheduler for better convergence
        scheduler = torch.optim.lr_scheduler.ReduceLROnPlateau(
            optimizer, mode='min', factor=0.5, patience=10
        )
        tokenizer = RobertaTokenizer.from_pretrained(SAVE_DIR, local_files_only=True)
        return model, tokenizer, scheduler
    else:
        tokenizer = RobertaTokenizer.from_pretrained(model_name)
        return None, tokenizer, None


