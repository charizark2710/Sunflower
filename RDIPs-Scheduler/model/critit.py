import socket
import os
import json
import time
from typing import Dict, List
import pandas as pd

import torch
import torch.nn as nn
from torch.optim import AdamW
from torch.amp.autocast_mode import autocast
import numpy as np

from base_model.model import CodeWithMetricsModel, load_model
from utils.common import extract_numeric_features, tokenize_codes, DEVICE

autocast('cuda', enabled=True, dtype=torch.bfloat16)

SOCKET_PATH = "../training.sock"
DELIM = "\nEND\n"
MODEL_NAME = "microsoft/codebert-base"
BATCH_SIZE = 20
LR = 1e-4
SAVE_DIR = "./saved_model"

loaded_model , loaded_tokenizer, loaded_scheduler = load_model(MODEL_NAME)

model = loaded_model
optimizer = None
scheduler = loaded_scheduler
tokenizer = loaded_tokenizer
loss_fn = nn.L1Loss()

# ----------------- training setup -----------------
buffer: List[Dict] = []
training_history = []
eval_history = []

# Set random seeds for reproducibility
torch.manual_seed(42)
np.random.seed(42)


def normalize_time_ic(ic):
    return np.log10(ic + 1)

def unnormalize_time_ic(norm_val):
    return np.pow(10, norm_val) + 1

def unnormalize_time_ic_torch(normalized):
    return torch.pow(10, normalized) + 1

def exportToCsv(obj):
    df = pd.DataFrame([obj])
    df = df.drop('code', axis=1)
    df.to_csv("output.csv", mode='a', header=not os.path.exists("output.csv"), index=False)

def save_model(model, tokenizer, save_dir=SAVE_DIR):
    os.makedirs(save_dir, exist_ok=True)
    if model is not None:
        torch.save(model.state_dict(), save_dir + "/model.pt")
        tokenizer.save_pretrained(save_dir)
    else:
        print("Model is None, cannot save.")

def split_code_by_semicolon(ids, max_length, sep_token_id=None):
    """
    Split token ids into chunks of max_length, ensuring splits end at the nearest semicolon.
    
    Args:
        ids: List[int] - tokenized sequence
        max_length: int - max chunk length
        sep_token_id: int - token ID of ';' (must be provided from tokenizer)
    """
    chunks = []
    i = 0
    n = len(ids)

    while i < n:
        end = min(i + max_length, n)

        # Try to cut at the nearest `;` before end
        cut = end
        if sep_token_id is not None:
            for j in range(end - 1, i, -1):  # search backward
                if ids[j] == sep_token_id:
                    cut = j + 1  # include `;`
                    break

        # If no `;` found, fallback to max_length cut
        if cut == i:
            cut = end

        chunks.append(ids[i:cut])
        i = cut  # move to next chunk

    return chunks


def prepare_data(idxs):
    codes = [batch[i]["code"] for i in idxs]
    ic_raw = [batch[i]["ic"] for i in idxs]

    target_train = []
    for t in ic_raw:
        transformed = normalize_time_ic(t)
        target_train.append(transformed)
    
    target_train = torch.tensor(target_train, dtype=torch.float, device=DEVICE).unsqueeze(1)
    
    input_ids, attention_mask = tokenize_codes(tokenizer, codes)
    input_ids = [x.to(DEVICE) for x in input_ids]
    attention_mask = [x.to(DEVICE) for x in attention_mask]
    extra_feats = extract_numeric_features([batch[i] for i in idxs])
    return input_ids, attention_mask, extra_feats, target_train, ic_raw

def train_on_batch(batch: List[Dict]):
    global model, optimizer, scheduler, training_history, eval_history, tokenizer
    batch_seed = hash(str(sorted([item.get('code', '')[:50] for item in batch]))) % 2**32
    rng = np.random.RandomState(batch_seed)
    
    # Shuffle batch indices and split into 80% train / 20% eval
    indices = list(range(len(batch)))
    rng.shuffle(indices)
    split_idx = max(1, int(0.8 * len(batch)))  # Ensure at least 1 sample for eval
    train_idx = indices[:split_idx]
    eval_idx = indices[split_idx:] if split_idx < len(batch) else indices[-1:]
    save_model(model, tokenizer)


    # Initialize model if first run
    if model is None:
        num_extra_feats = extract_numeric_features(batch).shape[1]
        model = CodeWithMetricsModel(MODEL_NAME, num_extra_feats).to(DEVICE)
        optimizer = AdamW(model.parameters(), lr=LR, weight_decay=0.01)
        # Add learning rate scheduler for better convergence
        scheduler = torch.optim.lr_scheduler.ReduceLROnPlateau(
            optimizer, mode='min', factor=0.5, patience=10
        )
    else:
        if optimizer is None:
            optimizer = AdamW(model.parameters(), lr=LR, weight_decay=0.01)
        scheduler = loaded_scheduler

    # Train split
    input_ids_train, attention_mask_train, extra_feats_train, targets_train, ic_raw_train = prepare_data(train_idx)
    model.train()
    optimizer.zero_grad()
    preds_train = model(input_ids_train, attention_mask_train, extra_feats_train)
    
    # Custom loss that handles different time ranges differently
    def custom_loss(preds, targets, ic_raw):
        # The preds and targets are already log-transformed
        return loss_fn(preds.squeeze(), targets.squeeze())

    loss_train = custom_loss(preds_train, targets_train, ic_raw_train)

    loss_train.backward()
    # # Gradient clipping to prevent explosion
    torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)
    optimizer.step()

    if len(eval_idx) > 0:
        input_ids_eval, attention_mask_eval, extra_feats_eval, targets_eval, ic_raw_eval = prepare_data(eval_idx)
        model.eval()
        with torch.no_grad():
            preds_eval = model(input_ids_eval, attention_mask_eval, extra_feats_eval)
            loss_eval = custom_loss(preds_eval, targets_eval, ic_raw_eval)
    else:
        loss_eval = loss_train

    if scheduler:
        scheduler.step(loss_eval.item())

    # Store history for trend analysis
    training_history.append(loss_train.item())
    eval_history.append(loss_eval.item())

    # Keep only recent history to save memory
    if len(training_history) > 100:
        training_history = training_history[-100:]
        eval_history = eval_history[-100:]

    # Convert predictions back to original scale for reporting
    preds_train_cpu = preds_train.detach().cpu().numpy().flatten()
    pred_ic_train = []
    
    for i, (pred, _) in enumerate(zip(preds_train_cpu, ic_raw_train)):
        pred_ic_train.append(pred)

    return loss_train.item(), loss_eval.item(), pred_ic_train

# ----------------- socket server -----------------
if os.path.exists(SOCKET_PATH):
    os.remove(SOCKET_PATH)

server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
server.bind(SOCKET_PATH)
server.listen(1)
print("Listening on", SOCKET_PATH, " — waiting for connection...")

conn, _ = server.accept()
print("Client connected")
recv_buffer = ""
batch_count = 0

try:
    while True:
        chunk = conn.recv(4096)
        if not chunk:
            print("Client finished")
            # Drain remaining buffer
            while len(buffer) >= BATCH_SIZE:
                batch = buffer[:BATCH_SIZE]
                t0 = time.time()
                loss_train, loss_eval, pred_ic = train_on_batch(batch)
                dt = time.time() - t0
                batch_count += 1
                
                print(f"Batch loss: train={loss_train:.6f} eval={loss_eval:.6f} time={dt:.3f}s")
                # Calculate aggregate error metrics
                actual_ics = [batch[i]['ic'] for i in range(min(3, len(batch)))]
                predicted_ics = pred_ic[:3]
                
                for i in range(min(3, len(batch))):
                    actual_ic = actual_ics[i]
                    predicted_ic = predicted_ics[i]
                    print(f"n_actual_ic: {normalize_time_ic(actual_ic):.3f} n_predicted_ic: {(predicted_ic):.3f}")
                    print(f"actual_ic: {actual_ic:.3f} predicted_ic: {unnormalize_time_ic(predicted_ic):.3f}")
                
                buffer = buffer[BATCH_SIZE:]
            break

        recv_buffer += chunk.decode("utf-8")

        # Process complete JSON messages delimited by DELIM
        while DELIM in recv_buffer:
            raw, recv_buffer = recv_buffer.split(DELIM, 1)
            raw = raw.strip()
            if not raw:
                continue
            try:
                obj = json.loads(raw)
                ic = float(obj.get("ic", 0.0))
                obj["ic"] = ic
                exportToCsv(obj)
                buffer.append(obj)
            except Exception as e:
                print("Bad message:", e)
                print(raw[:200] + "..." if len(raw) > 200 else raw)
                continue

            if len(buffer) >= BATCH_SIZE:
                batch = buffer[:BATCH_SIZE]
                t0 = time.time()
                loss_train, loss_eval, pred_ic = train_on_batch(batch)
                dt = time.time() - t0
                print(f"Batch loss: train={loss_train:.6f} eval={loss_eval:.6f} time={dt:.3f}s")
                batch_count += 1
                
                # Calculate aggregate error metrics
                actual_ics = [batch[i]['ic'] for i in range(min(3, len(batch)))]
                predicted_ics = pred_ic[:3]
                
                for i in range(min(3, len(batch))):
                    actual_ic = actual_ics[i]
                    predicted_ic = predicted_ics[i]
                    print(f"n_actual_ic: {normalize_time_ic(actual_ic):.3f} n_predicted_ic: {predicted_ic:.3f}")     # print(f"actual_ic: { actual_ic:.3f} predicted_ic: {predicted_ic:.3f}")                

                    print(f"actual_ic: {actual_ic:.3f} predicted_ic: {unnormalize_time_ic(predicted_ic):.3f}")        
                buffer = buffer[BATCH_SIZE:]
finally:
    try: conn.close()
    except: pass
    try: server.close()
    except: pass
    try: os.remove(SOCKET_PATH)
    except: pass
    print("Server shutdown")
    if training_history:
        pd.DataFrame({"train_loss": training_history, "eval_loss": eval_history}).to_csv("history.csv", index=False, mode='a', header=not os.path.exists("output.csv"))
        print(f"\nTraining Summary:")
        print(f"  Total batches: {batch_count}")
        print(f"  Final train loss: {training_history[-1]:.6f}")
        print(f"  Final eval loss: {eval_history[-1]:.6f}")
        print(f"  Best train loss: {min(training_history):.6f}")
        print(f"  Best eval loss: {min(eval_history):.6f}")
