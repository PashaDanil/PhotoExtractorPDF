import pika

# Establish a connection to RabbitMQ
connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
channel = connection.channel()

# Declare the queue 'hello' (it needs to exist)
channel.queue_declare(queue='hello')

# Define a callback function to process incoming messages
def callback(ch, method, properties, body):
    print(f" [x] Received '{body.decode()}'")

# Tell RabbitMQ to consume messages from the 'hello' queue
channel.basic_consume(queue='hello',
                      on_message_callback=callback,
                      auto_ack=True)

print(' [*] Waiting for messages. To exit press CTRL+C')

# Start consuming messages (this will block until a message is received)
channel.start_consuming()