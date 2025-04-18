import time
from pika import ConnectionParameters, BlockingConnection
import json
from dotenv import find_dotenv, load_dotenv
import os

load_dotenv(find_dotenv()) 
RABBITMQ_HOST = os.environ.get("RABBITMQ_HOST")

connection_params = ConnectionParameters(
    host=RABBITMQ_HOST,
    port=5672,
)


def sendToRabbitMQ(path: str, name: str):
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


def getFromRabbitMQ():
    result = {"url": None}

    def process_message(ch, method, properties, body):
        message = json.loads(body.decode("utf-8"))
        print(f"Получено сообщение: {message['url']}")
        result["url"] = message["url"]
        ch.basic_ack(delivery_tag=method.delivery_tag)
        ch.stop_consuming()

    with BlockingConnection(connection_params) as conn:
        with conn.channel() as ch:
            ch.queue_declare(queue="get_link")

            ch.basic_consume(
                queue="get_link",
                on_message_callback=process_message,
            )
            print("Жду сообщений")
            ch.start_consuming()

    return result["url"]

