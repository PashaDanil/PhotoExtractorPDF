from extractor.application.ports.consumer import Consumer
from extractor.exceptions.broker_exp import (
    BrokerConnectionError,
    BrokerFatalError,
    BrokerValidationError,
    BrokerUnexpectedError,
)
import aio_pika
import aio_pika.exceptions
import logging

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

    def __init__(self, url: str) -> None:
        """
        Конструктор consumer(а)
        Args:
            :param url: amqp://user:password@host:port/vhost
        Returns:
            None
        """
        self._url = url
        self._connection: aio_pika.abc.AbstractConnection | None = None
        self._channel: aio_pika.abc.AbstractChannel | None = None
        self._topics: list[str] = []
        logger.info("Инициализирован RabbitMQ consumer")

    async def connect(self) -> None:
        """Открыть соединение и канал."""
        try:
            self._connection = await aio_pika.connect_robust(self._url)
            self._channel = await self._connection.channel()
            # Prefetch: обрабатываем по одному сообщению за раз
            await self._channel.set_qos(prefetch_count=1)
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
            await self._channel.declare_queue(topic, durable=True)
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

    async def consume(self):
        """
        Запустить асинхронное считывание сообщений.
        Не блокирует event loop — использует async generator от aio_pika.
        """
        if not self._topics:
            raise BrokerValidationError(
                "Нет топиков для подписки. Вызови subscribe() перед consume()"
            )

        try:
            for topic in self._topics:
                queue = await self._channel.get_queue(topic)
                # on_message_callback вызывается как корутина
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

    async def _callback(self, message: aio_pika.abc.AbstractIncomingMessage):
        """
        Вызывается event loop'ом для каждого входящего сообщения.
        """
        async with message.process(requeue=True):          # nack+requeue при исключении
            try:
                payload = message.body.decode("utf-8")
                logger.debug("Получено сообщение: %s", payload)
                logger.info(
                    "Получено сообщение delivery_tag=%s", message.delivery_tag
                )

                # TODO: передать payload в use-case
                await self._handle(payload)

            except UnicodeDecodeError as e:
                logger.exception(
                    "Ошибка декодирования delivery_tag=%s", message.delivery_tag
                )
                raise BrokerValidationError(
                    f"Ошибка декодирования delivery_tag={message.delivery_tag}"
                ) from e

    async def _handle(self, payload: str):
        """Заинжектить use-case сюда."""
        raise NotImplementedError