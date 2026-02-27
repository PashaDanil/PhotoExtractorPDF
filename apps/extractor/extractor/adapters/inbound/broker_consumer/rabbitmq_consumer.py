from extractor.application.ports.consumer import Consumer
from pika.adapters.blocking_connection import BlockingConnection
from extractor.exceptions.broker_exp import (BrokerConnectionError,
                                             BrokerPublishError,
                                             BrokerConsumeError,
                                             BrokerFatalError,
                                             BrokerValidationError,
                                             BrokerUnexpectedError
                                             )
import pika
import logging

RETRIABLE_EXCEPTIONS = (
    pika.exceptions.ConnectionClosed,
    pika.exceptions.ConnectionClosedByBroker,
    pika.exceptions.StreamLostError,
    pika.exceptions.AMQPHeartbeatTimeout,
    pika.exceptions.ConnectionBlockedTimeout,
)

FATAL_EXCEPTIONS = (
    pika.exceptions.ChannelClosedByBroker,
    pika.exceptions.ChannelWrongStateError,
    pika.exceptions.AuthenticationError,
    pika.exceptions.ProbableAuthenticationError,
    pika.exceptions.IncompatibleProtocolError,
)

logger = logging.getLogger(__name__)


class RabbitConsumer(Consumer):

    def __init__(self, connection: BlockingConnection):
        self._connection = connection
        self._channel = None
        self._topics = []
        logger.info("Инициализирован RabbitMQ consumer")


    def connect(self):
        """
        Метод для подключения к брокеру сообщений
        """
        try:
            self._channel = self._connection.channel()
            logger.info("Подключение к RabbitMQ брокеру, канал открыт")
        except FATAL_EXCEPTIONS as e:
            logger.exception(
                "Фатальная ошибка при подключении к брокеру",
            )
            raise BrokerFatalError(
                f"Фатальная ошибка при подключении к брокеру"
            ) from e

        except RETRIABLE_EXCEPTIONS as e:
            logger.warning(
                "Не удалось подключиться к брокеру",
            )
            raise BrokerConnectionError(
                "Не удалось подключиться к брокеру"
            ) from e

        except Exception as e:
            logger.exception(
                "Непредвиденная ошибка при подключении к брокеру",
            )
            raise BrokerUnexpectedError(
                "Непредвиденная ошибка при обработке сообщения"
            ) from e

    def subscribe(self, topic: str):
        if not topic or not topic.strip():
            logger.error("Название топика не может быть пустым")
            raise BrokerValidationError("Название топика не может быть пустой строкой")

        try:
            self._channel.queue_declare(queue=topic)
            self._topics.append(topic)
            logger.info("Подписка на topic: %s", topic)

        except FATAL_EXCEPTIONS as e:
            logger.exception("Фатальная ошибка при объявлении очереди %s", topic)
            raise BrokerFatalError(f"Ошибка при объявлении очереди {topic}") from e

        except RETRIABLE_EXCEPTIONS as e:
            logger.warning("Соединение потеряно при объявлении очереди %s: %s", topic, e)
            raise BrokerConnectionError("Потеряно соединение с брокером") from e

        except Exception as e:
            logger.exception("Непредвиденная ошибка при объявлении очереди %s", topic)
            raise BrokerUnexpectedError("Непредвиденная ошибка") from e

    def consume(self):
        """
        Начать считывать сообщения с брокера
        """
        if not self._topics:
            raise BrokerValidationError("Нет топиков для подписки. Вызови subscribe() перед consume()")
        try:
            # У всех топиков один callback
            for topic in self._topics:
                self._channel.basic_consume(queue=topic,
                                            on_message_callback=self.callback,
                                            auto_ack=False)
            self._channel.start_consuming()

        except FATAL_EXCEPTIONS as e:
            logger.exception(
                "Фатальная ошибка при подписке на топик."
            )
            raise BrokerFatalError(
                f"Фатальная ошибка при при подписке на топик"
            ) from e

        except RETRIABLE_EXCEPTIONS as e:
            logger.warning(
                "Не удалось подписаться на топик: %s.",
                e
            )
            raise BrokerConnectionError(
                "Не удалось подписаться на топик"
            ) from e

        except Exception as e:
            logger.exception(
                "Непредвиденная ошибка попытке подписаться на топик.",
            )
            raise BrokerUnexpectedError(
                "Непредвиденная ошибка при попытке подписаться на топик"
            ) from e


    def disconnect(self):
        """
        Остановка считывания сообщений с брокера
        """
        try:
            if self._channel and self._channel.is_open:
                self._channel.close()
                logger.info("Канал закрыт")
        except Exception as e:
            logger.warning("Ошибка при закрытии канала: %s", e)

        try:
            if self._connection and self._connection.is_open:
                self._connection.close()
                logger.info("Соединение закрыто")
        except Exception as e:
            logger.warning("Ошибка при закрытии соединения: %s", e)

    def callback(self, ch, method, properties, body):
        try:
            payload = body.decode("utf-8")
            logger.debug("Получено сообщение: %s", payload)
            logger.info("Получено сообщение delivery_tag=%s",
                        method.delivery_tag)

            ch.basic_ack(delivery_tag=method.delivery_tag)

        except UnicodeDecodeError as e:
            logger.exception(
                "Ошибка в данных сообщения delivery_tag=%s",
                method.delivery_tag
            )

            raise BrokerValidationError(
                f"Ошибка при декодировании данных сообщения "
                f"delivery_tag={method.delivery_tag}"
            ) from e

        except FATAL_EXCEPTIONS as e:
            logger.exception(
                "Фатальная ошибка при ACK delivery_tag=%s.",
                method.delivery_tag
            )
            raise BrokerFatalError(
                f"Фатальная ошибка при ACK delivery_tag={method.delivery_tag}"
            ) from e

        except RETRIABLE_EXCEPTIONS as e:
            logger.warning(
                "Не удалось подтвердить сообщение delivery_tag=%s: %s.",
                method.delivery_tag,
                e
            )
            raise BrokerConnectionError(
                "Не удалось подключиться к брокеру"
            ) from e

        except Exception as e:
            logger.exception(
                "Непредвиденная ошибка при ACK delivery_tag=%s.",
                method.delivery_tag,
            )
            raise BrokerUnexpectedError(
                "Непредвиденная ошибка при обработке сообщения"
            ) from e


