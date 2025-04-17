from pika import ConnectionParameters, BlockingConnection
import json

connection_params = ConnectionParameters(
    host="rabbitmq",
    # host="localhost",
    port=5672,
)


def connectRabbitMQ(path: str, name: str):
    with BlockingConnection(connection_params) as conn:
        with conn.channel() as ch:
            ch.queue_declare(queue="file")

            message = {
                "path": path,
                "name": name
            }

            ch.basic_publish(
                exchange="",
                routing_key="file",
                body=json.dumps(message).encode("utf-8"),
            )
            print("Message sent")

