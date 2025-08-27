import socket
import os
import json
import time
from typing import Dict, List
import pandas as pd

import torch
import torch.nn as nn
from transformers import RobertaTokenizer, RobertaModel
from torch.optim import AdamW
import numpy as np

SAVE_DIR = "./saved_model"
DEVICE = "cpu"
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
            nn.Linear(128, 1),
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
        return output


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


