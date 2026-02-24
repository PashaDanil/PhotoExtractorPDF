from extractor.application.ports.consumer import Consumer
from pika.adapters.blocking_connection import BlockingConnection


class RabbitConsumer(Consumer):

    def __init__(self, connection: BlockingConnection):
        # connection = pika.BlockingConnection(pika.ConnectionParameters('localhost'))
        self._connection = connection
        self._channel = None
        self._topics = []


    def connect(self):
        """
        Метод для подключения к брокеру сообщений
        """
        self._channel = self._connection.channel()


    def subscribe(self, topic: str):
        """
        Метод для подписки на топик
        Args:
            topic: Строка - название топика
        Returns:
            None
        """
        self._topics.append(topic)
        self._channel.queue_declare(queue=topic)


    def consume(self):
        """
        Начать считывать сообщения с брокера
        """
        for topic in self._topics:
            self._channel.basic_consume(queue=topic,
                                  on_message_callback=self.callback,
                                  auto_ack=True)
        self._channel.start_consuming()


    def disconnect(self):
        """
        Остановка считывания сообщений с брокера
        """
        ...


    def callback(ch, method, properties, body):
        print(f" [x] Received '{body.decode()}'")
