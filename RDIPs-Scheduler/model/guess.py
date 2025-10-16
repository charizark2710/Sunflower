import socket
import os
import json
import torch
from base_model.model import load_model
from utils.common import extract_numeric_features, tokenize_codes, DEVICE

SOCKET_PATH = "../guess.sock"

# Load model and tokenizer
model, tokenizer, scheduler = load_model("microsoft/codebert-base")
if (model is None) or (tokenizer is None):
    exit(1)    

model.to(DEVICE)
model.eval()

# Ensure old socket is removed
if os.path.exists(SOCKET_PATH):
    os.remove(SOCKET_PATH)
# Create UNIX socket server
server_sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
server_sock.bind(SOCKET_PATH)
server_sock.listen(1)
print(f"Listening on {SOCKET_PATH}")

def handle_request(data: dict):
    """
    data is expected to be a dict with the required inputs.
    Example:
    { "code": "function example() { return 42; }", "metrics": {...} }
    """
    # Example: just using code text for tokenization
    code_str = data.get("code", "")
    input_ids, attention_mask = tokenize_codes(tokenizer, [code_str])

    input_ids = [x.to(DEVICE) for x in input_ids]
    attention_mask = [x.to(DEVICE) for x in attention_mask]
    extra_feats = extract_numeric_features([data])

    # Model forward
    with torch.no_grad():
        outputs = model(input_ids, attention_mask, extra_feats) # type: ignore

    ic = pow(10, outputs[0].detach().cpu().numpy().flatten())
    cycle = pow(10, outputs[2].detach().cpu().numpy().flatten())

    return {
        "ic": ic.tolist()[0],
        "cycle": cycle.tolist()[0],
        "input_ids": torch.cat(input_ids).view(-1).tolist(),
        "attention_mask": torch.cat(attention_mask).view(-1).tolist(),
    }

while True:
    conn, _ = server_sock.accept()
    with conn:
        print("Client connected")
        try:
            raw_data = conn.recv(4096)
            if not raw_data:
                continue

            # Decode & parse
            data_str = raw_data.decode("utf-8").strip()
            try:
                data_json = json.loads(data_str)
            except json.JSONDecodeError:
                conn.sendall(b'{"error": "Invalid JSON"}\n')
                continue

            # Process through model
            result = handle_request(data_json)
            # Send back JSON
            conn.sendall((json.dumps(result) + "\n").encode("utf-8"))

        except Exception as e:
            err_msg = json.dumps({"error": str(e)}) + "\n"
            conn.sendall(err_msg.encode("utf-8"))
