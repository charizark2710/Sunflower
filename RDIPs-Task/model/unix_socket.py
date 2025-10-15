import socket
import json
import os
import torch
import psutil
import argparse

from ActorCritic import ActorCritic

def socket_listener(socket_path="../guess.sock"):
    # Remove existing socket file if it exists
    try:
        os.unlink(socket_path)
    except FileNotFoundError:
        pass
    
    server = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    server.bind(socket_path)
    server.listen(1)
    
    print(f"Listening on {socket_path} — waiting for connection...")
    
    return server

def handle_connection(server_socket, is_training=False):
    """Accept and handle a single connection"""
    conn, _ = server_socket.accept()
    print("Connection established.")
    
    actorCritic = ActorCritic(
        num_extra_feats=7,
        conn=conn
    )
    try:
        # Receive initial request
        initial_data = recv_msg(conn)
        if not initial_data:
            return
        
        ic = initial_data.get("ic")
        input_ids = torch.tensor(initial_data.get("input_ids"))
        attention_mask = torch.tensor(initial_data.get("attention_mask"))
        cycle = initial_data.get("cycle")
        ast_metric = initial_data.get("ast_metric")

        memory_total = psutil.virtual_memory().total
        cpu_usage = psutil.cpu_percent(interval=0) / 100
        mem_usage = psutil.virtual_memory().used / memory_total

        response = None
        if is_training:
            action, state = actorCritic.train([{
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
        response = {"isSkip": action != 1, "confidence": state["certainty"]}
        send_msg(conn, response)

    except Exception as e:
        print(f"Error handling connection: {e}")
    finally:
        conn.close()

def recv_msg(conn: socket.socket):
    """Receive and parse JSON message"""
    try:
        data = b''
        not_done = True
        while not_done:
            chunk = conn.recv(4096)
            if not chunk:
                break
            if chunk.endswith(b'#END#'):
                not_done = False
                chunk = chunk.replace(b'#END#', b'')
            data += chunk
        if not data or data == None or data == b'' or len(data) == 0:
            return None

        # Parse JSON
        return json.loads(data.decode('utf-8'))
    except json.JSONDecodeError as e:
        print(f"Failed to decode JSON: {e}")
        return None
    except Exception as e:
        print(f"Error receiving message: {e}")
        return None

def send_msg(conn: socket.socket, data: dict):
    """Send JSON message"""
    try:
        message = json.dumps(data).encode('utf-8')
        conn.sendall(message)
    except Exception as e:
        print(f"Error sending message: {e}")
        raise
