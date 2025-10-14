import socket
import json
import os
import torch

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

def handle_connection(server_socket):
    """Accept and handle a single connection"""
    conn, _ = server_socket.accept()
    print("Connection established.")
    
    try:
        # Receive initial request
        initial_data = recv_msg(conn)
        if not initial_data:
            return
        
        ic = initial_data.get("pred_ic")
        input_ids = torch.tensor(initial_data.get("input_ids"))
        attention_masks = torch.tensor(initial_data.get("attention_masks"))
        
        cycle = initial_data.get("pred_cycle")
        
        # Process and generate prediction
        score, confidence = your_model_predict(ic, "", cycle)
        
        # Send response back
        response = {"score": score, "confidence": confidence}
        send_msg(conn, response)
        
        # If confidence < 0.6, expect actual_ic feedback
        if confidence < 0.6:
            feedback = recv_msg(conn)
            if feedback and "ic" in feedback:
                actual_ic = feedback["ic"]
                # Update your model or log the actual_ic
                print(f"Received actual IC: {actual_ic}")

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

def your_model_predict(ic, code, cycle):
    """Replace with your actual model prediction logic"""
    # Placeholder - implement your model here
    score = 0.8
    confidence = 0.7
    
    
    return score, confidence

# # Main loop
# def main():

#     server = socket_listener("../guess.sock")
    
#     try:
#         while True:
#             handle_connection(server)
#     except KeyboardInterrupt:
#         print("\nShutting down...")
#     finally:
#         server.close()
#         try:
#             os.unlink("../guess.sock")
#         except:
#             pass

# if __name__ == "__main__":
#     main()