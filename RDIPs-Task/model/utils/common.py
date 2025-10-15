import torch
from typing import Dict, List
import os
SAVE_DIR = "./saved_model"
from torch.optim import AdamW
from base_model.model import ActorModel, CriticModel

DEVICE = 'cuda' if torch.cuda.is_available() else 'cpu'
MAX_TOKEN_LEN = 512
CPU_RATE = 3.5e9
LR = 1e-4

def load_model(num_extra_feats):
    if not os.path.exists(SAVE_DIR + "/actor_model.pt") or not os.path.exists(SAVE_DIR + "/critic_model.pt"):
        actor = ActorModel(num_extra_feats = num_extra_feats)
        critic = CriticModel()
        actor_optimizer = AdamW(actor.parameters(), lr=LR)
        critic_optimizer = AdamW(critic.parameters(), lr=LR)
        return actor, critic
        
    actor_model = ActorModel(num_extra_feats = num_extra_feats)
    actor_model.load_state_dict(torch.load(SAVE_DIR + "/actor_model.pt", map_location=DEVICE))
    critic_model = CriticModel()
    critic_model.load_state_dict(torch.load(SAVE_DIR + "/critic_model.pt", map_location=DEVICE))

    return actor_model, critic_model