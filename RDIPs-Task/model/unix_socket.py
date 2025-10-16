import socket
import json
import os

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

def recv_msg(conn: socket.socket):
    """Receive and parse JSON message"""
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

    if not data:
        return None

    try:
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
