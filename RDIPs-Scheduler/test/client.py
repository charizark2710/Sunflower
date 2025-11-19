import pika
import csv
import time
import json
from typing import Optional

class RabbitMQPublisher:
    def __init__(self, host: str = 'localhost', port: int = 5672, 
                 username: str = 'admin', password: str = 'admin',
                 exchange: str = 'Scheduler-Gateway',
                 routing_key: str = 'Scheduler-Gateway.JS_CODE'):
        """
        Initialize RabbitMQ connection for stream queue publishing
        """
        self.host = host
        self.port = port
        self.username = username
        self.password = password
        self.exchange = exchange
        self.routing_key = routing_key
        self.connection = None
        self.channel = None
    
    def connect(self):
        """Establish connection to RabbitMQ"""
        credentials = pika.PlainCredentials(self.username, self.password)
        parameters = pika.ConnectionParameters(
            host=self.host,
            port=self.port,
            credentials=credentials,
            heartbeat=600,
            blocked_connection_timeout=300
        )
        
        self.connection = pika.BlockingConnection(parameters)
        self.channel = self.connection.channel()
        
        # Declare exchange (if it doesn't exist)
        self.channel.exchange_declare(
            exchange=self.exchange,
            exchange_type='fanout',
            durable=True
        )
        
        print(f"Connected to RabbitMQ at {self.host}:{self.port}")
    
    def publish_message(self, message: str) -> bool:
        """
        Publish a single message to the queue
        
        Args:
            message: The message content to publish
            
        Returns:
            bool: True if successful, False otherwise
        """
        try:
            self.channel.basic_publish(
                exchange=self.exchange,
                routing_key=self.routing_key,
                body=message,
                properties=pika.BasicProperties(
                    delivery_mode=2,  # Make message persistent
                    content_type='text/plain'
                )
            )
            return True
        except Exception as e:
            print(f"Error publishing message: {e}")
            return False
    
    def close(self):
        """Close RabbitMQ connection"""
        if self.connection and not self.connection.is_closed:
            self.connection.close()
            print("Connection closed")


def read_and_publish_csv(csv_file_path: str, 
                         publisher: RabbitMQPublisher,
                         batch_size: int = 10,
                         delay_between_batches: float = 1.0,
                         delay_between_messages: float = 0.1):
    """
    Read CSV file and publish 'code' column values to RabbitMQ stream queue
    
    Args:
        csv_file_path: Path to the CSV file
        publisher: RabbitMQPublisher instance
        batch_size: Number of messages to send before pausing
        delay_between_batches: Seconds to wait between batches
        delay_between_messages: Seconds to wait between individual messages
    """
    total_sent = 0
    total_failed = 0
    
    try:
        with open(csv_file_path, 'r', encoding='utf-8') as csvfile:
            reader = csv.DictReader(csvfile)
            
            # Check if 'code' column exists
            if 'code' not in reader.fieldnames:
                print(f"Error: 'code' column not found in CSV. Available columns: {reader.fieldnames}")
                return
            
            print(f"Starting to process CSV file: {csv_file_path}")
            print(f"Batch size: {batch_size}, Delay between batches: {delay_between_batches}s")
            print("-" * 60)
            
            batch_count = 0
            
            for row in reader:
                code = row.get('code', '').strip()
                
                # Skip empty codes
                if not code:
                    continue
                
                # Publish message
                if publisher.publish_message(code):
                    total_sent += 1
                    print(f"Sent [{total_sent}]: {code[:50]}{'...' if len(code) > 50 else ''}")
                else:
                    total_failed += 1
                    print(f"Failed to send: {code[:50]}{'...' if len(code) > 50 else ''}")
                
                batch_count += 1
                
                # Delay between messages
                time.sleep(delay_between_messages)
                
                # Batch delay
                if batch_count >= batch_size:
                    print(f"\n--- Batch complete ({batch_size} messages). Waiting {delay_between_batches}s... ---\n")
                    time.sleep(delay_between_batches)
                    batch_count = 0
            
            print("-" * 60)
            print(f"\nProcessing complete!")
            print(f"Total messages sent: {total_sent}")
            print(f"Total messages failed: {total_failed}")
            
    except FileNotFoundError:
        print(f"Error: File '{csv_file_path}' not found")
    except Exception as e:
        print(f"Error reading CSV file: {e}")


def main():
    # Configuration
    CSV_FILE = 'input.csv'  # Change this to your CSV file path
    
    # RabbitMQ Configuration
    RABBITMQ_HOST = 'localhost'
    RABBITMQ_PORT = 5672
    RABBITMQ_USER = 'admin'
    RABBITMQ_PASSWORD = 'admin'
    EXCHANGE = 'Scheduler-Gateway'
    ROUTING_KEY = 'Scheduler-Gateway.JS_CODE'
    
    # Rate limiting configuration
    BATCH_SIZE = 2  # Send 5 messages per batch
    DELAY_BETWEEN_BATCHES = 10.0  # Wait 2 seconds between batches
    DELAY_BETWEEN_MESSAGES = 0.1  # Wait 0.1 seconds between individual messages
    
    # Create publisher and connect
    publisher = RabbitMQPublisher(
        host=RABBITMQ_HOST,
        port=RABBITMQ_PORT,
        username=RABBITMQ_USER,
        password=RABBITMQ_PASSWORD,
        exchange=EXCHANGE,
        routing_key=ROUTING_KEY
    )
    
    try:
        publisher.connect()
        
        # Read CSV and publish messages
        read_and_publish_csv(
            csv_file_path=CSV_FILE,
            publisher=publisher,
            batch_size=BATCH_SIZE,
            delay_between_batches=DELAY_BETWEEN_BATCHES,
            delay_between_messages=DELAY_BETWEEN_MESSAGES
        )
        
    except Exception as e:
        print(f"Error: {e}")
    finally:
        publisher.close()


if __name__ == "__main__":
    main()