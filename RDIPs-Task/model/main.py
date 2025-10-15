
import argparse
import os
from unix_socket import socket_listener, handle_connection

# Main loop
def main():
    parser = argparse.ArgumentParser(description="Run actor-critic socket listener.")
    parser.add_argument("--is_training", type=lambda x: x.lower() in ['true', '1', 'yes'], default=False,
                        help="Run in training mode if True, else inference mode.")
    parser.add_argument("--socket_path", type=str, default="../guess.sock",
                        help="Path to the UNIX socket file.")
    
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