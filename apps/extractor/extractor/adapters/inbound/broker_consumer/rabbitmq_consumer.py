import logging
import aio_pika
import aio_pika.exceptions
from aio_pika.abc import AbstractRobustConnection, AbstractRobustChannel, AbstractIncomingMessage
from infrastructure.config.rabbitmq import RabbitMQConfig
from extractor.application.ports.consumer import Consumer
from extractor.exceptions.broker_exp import (
    BrokerConnectionError,
    BrokerFatalError,
    BrokerValidationError,
    BrokerUnexpectedError,
)

RETRIABLE_EXCEPTIONS = (
    aio_pika.exceptions.AMQPConnectionError,
    ConnectionError,
)

FATAL_EXCEPTIONS = (
    aio_pika.exceptions.AuthenticationError,
    aio_pika.exceptions.ProbableAuthenticationError,
    aio_pika.exceptions.IncompatibleProtocolError,
    aio_pika.exceptions.ChannelInvalidStateError,
)

logger = logging.getLogger(__name__)


class RabbitConsumer(Consumer):

    def __init__(self, config: RabbitMQConfig) -> None:
        """
        Инициализирует consumer для RabbitMQ.
        Args:
            config: Конфигурация подключения к RabbitMQ
        Returns:
            None
        """
        self._config = config
        self._connection: AbstractRobustConnection | None = None
        self._channel: AbstractRobustChannel | None = None
        self._topics: list[str] = []
        logger.info("Инициализирован RabbitMQ consumer")

    async def connect(self) -> None:
        """Открыть соединение и канал."""
        try:
            self._connection = await aio_pika.connect_robust(self._config.url)
            self._channel = await self._connection.channel()
            await self._channel.set_qos(prefetch_count=self._config.prefetch_count)
            logger.info("Подключение к RabbitMQ: канал открыт")

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
        """Закрыть канал и соединение."""
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

    async def subscribe(self, topic: str) -> None:
        """Объявить очередь и сохранить имя топика.
        Args:
            topic: Строка - название топика
        Returns:
            None
        """
        if not topic or not topic.strip():
            raise BrokerValidationError("Название топика не может быть пустой строкой")

        try:
            await self._channel.declare_queue(topic, durable=self._config.durable)
            self._topics.append(topic)
            logger.info("Подписка на topic: %s", topic)

        except FATAL_EXCEPTIONS as e:
            logger.exception("Фатальная ошибка при объявлении очереди %s", topic)
            raise BrokerFatalError(f"Ошибка при объявлении очереди {topic}") from e

        except RETRIABLE_EXCEPTIONS as e:
            logger.warning("Соединение потеряно при объявлении очереди %s", topic)
            raise BrokerConnectionError("Потеряно соединение с брокером") from e

        except Exception as e:
            logger.exception("Непредвиденная ошибка при объявлении очереди %s", topic)
            raise BrokerUnexpectedError("Непредвиденная ошибка") from e

    async def consume(self) -> None:
        """
        Запустить асинхронное считывание сообщений.
        Не блокирует event loop — использует async generator от aio_pika.
        Returns:
            None
        """
        if not self._topics:
            raise BrokerValidationError(
                "Нет топиков для подписки. Вызови subscribe() перед consume()"
            )

        try:
            for topic in self._topics:
                queue = await self._channel.declare_queue(topic, durable=self._config.durable, passive=True)
                await queue.consume(self._callback)
            logger.info("Начато потребление сообщений из: %s", self._topics)

        except FATAL_EXCEPTIONS as e:
            logger.exception("Фатальная ошибка при подписке на топики")
            raise BrokerFatalError("Фатальная ошибка при подписке на топики") from e

        except RETRIABLE_EXCEPTIONS as e:
            logger.warning("Соединение потеряно при подписке: %s", e)
            raise BrokerConnectionError("Не удалось подписаться на топики") from e

        except Exception as e:
            logger.exception("Непредвиденная ошибка при подписке на топики")
            raise BrokerUnexpectedError("Непредвиденная ошибка при подписке") from e

    async def _callback(self, message: AbstractIncomingMessage) -> None:
        """
        Обрабатывает входящее сообщение от RabbitMQ.

        Вызывается event loop'ом для каждого полученного сообщения.

        Args:
            message: Входящее сообщение из очереди RabbitMQ.
        """
        try:
            payload = message.body.decode("utf-8")
            logger.debug("Получено сообщение: %s", payload)
            logger.info("Получено сообщение delivery_tag=%s", message.delivery_tag)
            await self._handle(payload)
            await message.ack()

        except UnicodeDecodeError as e:
            logger.exception("Ошибка декодирования delivery_tag=%s", message.delivery_tag)
            await message.nack(requeue=False)
            raise BrokerValidationError(
                f"Ошибка декодирования delivery_tag={message.delivery_tag}"
            ) from e

        except BrokerConnectionError:
            await message.nack(requeue=True)
            raise

        except Exception:
            await message.nack(requeue=False)
            raise

    async def _handle(self, payload: str) -> None:
        """TODO: Заинжектить use-case сюда."""
        raise NotImplementedError