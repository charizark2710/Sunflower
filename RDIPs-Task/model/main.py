
import argparse
import os
from unix_socket import socket_listener, recv_msg, send_msg
from utils.common import extract_numeric_features

import torch
import psutil
import argparse

from ActorCritic import ActorCritic

actorCritic: ActorCritic
def handle_connection(server_socket, is_training=False):
    global actorCritic
    """Accept and handle a single connection"""
    conn, _ = server_socket.accept()
    print("Connection established.")
    actorCritic.updateConn(conn)
    try:
        # Receive initial request
        initial_data = recv_msg(conn)
        if not initial_data:
            return
        
        ic = initial_data.get("pred_ic")
        input_ids = [torch.tensor(initial_data.get("input_ids")).unsqueeze(0)]
        attention_mask = [torch.tensor(initial_data.get("attention_mask")).unsqueeze(0)]
        cycle = initial_data.get("pred_cycle")
        ast_metric = extract_numeric_features([initial_data.get("ast_metric")])

        memory_total = psutil.virtual_memory().total
        cpu_usage = psutil.cpu_percent(interval=0) / 100
        mem_usage = psutil.virtual_memory().used / memory_total

        response = None
        if is_training:
            actorCritic.train([{
                "ast_metric": ast_metric,
                "code_tokens": input_ids,
                "attention_mask": attention_mask,
                "master_pred": {
                    "ic": ic,
                    "cycle": cycle
                }
            }])
        else:
            action, state = actorCritic.select_action(input_ids, attention_mask, state={
                "cpu": cpu_usage,
                "memory": mem_usage
            }, ast_metric=ast_metric, master_pred={
                "ic": ic,
                "cycle": cycle
            })
            # Send response back
            response = {"isSkip": action != 1, "confidence": state["certainty"], "done": True}
            send_msg(conn, response)

    except Exception as e:
        print(f"Error handling connection: {e}")
    finally:
        conn.close()


# Main loop
def main():
    global actorCritic
    parser = argparse.ArgumentParser(description="Run actor-critic socket listener.")
    parser.add_argument("--is_training", type=lambda x: x.lower() in ['true', '1', 'yes'], default=True,
                        help="Run in training mode if True, else inference mode.")
    parser.add_argument("--socket_path", type=str, default="../guess.sock",
                        help="Path to the UNIX socket file.")
    actorCritic = ActorCritic(
            num_extra_feats=7,
            conn=None,
    )
    args = parser.parse_args()
    print(f"Running in {'training' if args.is_training else 'inference'} mode.")

    server = socket_listener(args.socket_path)
    try:
        while True:
            handle_connection(server, is_training=args.is_training)
    except KeyboardInterrupt:
        print("\nShutting down...")
    finally:
        server.close()
        try:
            os.unlink(args.socket_path)
        except:
            pass
if __name__ == "__main__":
    main()