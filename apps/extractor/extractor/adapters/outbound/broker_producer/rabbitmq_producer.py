import logging
import aio_pika
import aio_pika.exceptions
from aio_pika.abc import AbstractRobustConnection, AbstractRobustChannel, AbstractRobustExchange
from extractor.application.ports.producer import Producer, Message
from infrastructure.config.rabbitmq import RabbitMQConfig
from extractor.exceptions.broker_exp import (
    BrokerConnectionError,
    BrokerFatalError,
    BrokerUnexpectedError,
)

RETRIABLE_EXCEPTIONS = (
    aio_pika.exceptions.AMQPConnectionError,
    ConnectionError
)

FATAL_EXCEPTIONS = (
    aio_pika.exceptions.AuthenticationError,
    aio_pika.exceptions.ProbableAuthenticationError,
    aio_pika.exceptions.IncompatibleProtocolError,
    aio_pika.exceptions.ChannelInvalidStateError,
)

logger = logging.getLogger(__name__)


class RabbitMQProducer(Producer):
    """
    Producer для брокера сообщений RabbitMQ
    """

    def __init__(self,
                 config: RabbitMQConfig
                 ):
        self._config: RabbitMQConfig = config
        self._connection: AbstractRobustConnection | None = None
        self._channel: AbstractRobustChannel | None = None
        self._exchange: AbstractRobustExchange | None = None
        logger.info("Инициализирован RabbitMQ producer")

    async def connect(self) -> None:
        """
        Устанавливает соединение с брокером сообщений.
        """
        try:
            self._connection = await aio_pika.connect_robust(
                self._config.url
            )
            self._channel = await self._connection.channel()
            self._exchange = await self._channel.declare_exchange(
                name=self._config.exchange,
                type=self._config.exchange_type,
                durable=self._config.durable
            )
            logger.info("Установлено соединение с брокером RabbitMQ")
        except FATAL_EXCEPTIONS as e:
            logger.exception("Фатальная ошибка при подключении к брокеру")
            raise BrokerFatalError("Фатальная ошибка при подключении к брокеру") from e

        except RETRIABLE_EXCEPTIONS as e:
            logger.warning("Не удалось подключиться к брокеру")
            raise BrokerConnectionError("Не удалось подключиться к брокеру") from e

        except Exception as e:
            logger.exception("Непредвиденная ошибка при подключении к брокеру")
            raise BrokerUnexpectedError("Непредвиденная ошибка при подключении") from e


    async def disconnect(self) -> None:
        """
        Закрывает соединение с брокером сообщений.
        """
        try:
            if self._channel and not self._channel.is_closed:
                await self._channel.close()
                logger.info("Канал закрыт")
        except Exception as e:
            logger.warning("Ошибка при закрытии канала: %s", e)

        try:
            if self._connection and not self._connection.is_closed:
                await self._connection.close()
                logger.info("Соединение закрыто")
        except Exception as e:
            logger.warning("Ошибка при закрытии соединения: %s", e)


    async def publish(
                self,
                destination: str,
                message: Message
        ) -> None:
        """
        Публикует сообщение в RabbitMQ exchange.

        Args:
            message: Сообщение для отправки.
            destination: Ключ маршрутизации - имя очереди.
        Returns:
            None
        """
        if not self._exchange:
            logger.warning("Ошибка при попытке отправить сообщение в брокер RabbitMQ. Соединение не инициализировано.")
            raise BrokerConnectionError("Нет подключения к брокеру. Необходимо вызывать connect() перед publish()")
        try:
            amqp_message = aio_pika.Message(
                body=message.body,
                content_type=message.content_type,
                headers=message.headers,
                delivery_mode=self._config.delivery_mode,
            )

            await self._exchange.publish(
                amqp_message,
                routing_key=destination,
            )
            logger.info("Сообщение отправлено в брокер RabbitMQ")
        except FATAL_EXCEPTIONS as e:
            logger.exception("Фатальная ошибка при публикации сообщения")
            raise BrokerFatalError("Фатальная ошибка при публикации сообщения") from e

        except RETRIABLE_EXCEPTIONS as e:
            logger.warning("Не удалось подключиться к брокеру")
            raise BrokerConnectionError("Не удалось подключиться к брокеру") from e

        except Exception as e:
            logger.exception("Непредвиденная ошибка при публикации сообщения")
            raise BrokerUnexpectedError("Непредвиденная ошибка при публикации сообщения") from e

    async def is_connected(self) -> bool:
        """
        Проверить, активно ли соединение.

        Returns:
            bool: Возвращает True, если соединение с брокером активно (не закрыто), иначе False.
        """
        if not self._connection:
            return False
        try:
            return not self._connection.is_closed
        except Exception as e:
            logger.exception("Непредвиденная ошибка при проверке подключения к брокеру")
            return False

