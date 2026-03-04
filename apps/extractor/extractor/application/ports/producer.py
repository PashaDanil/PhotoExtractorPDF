from abc import ABC, abstractmethod
from typing import Optional, Any
from dataclasses import dataclass, field

@dataclass(frozen=True)
class Message:
    """
    Класс сообщения, передаваемого брокеру.

    body: Тело сообщения.
    key: Ключ маршрутизации.
    headers: Дополнительные заголовки сообщения.
    content_type: MIME-тип содержимого.
    """

    body: bytes
    key: Optional[str] = None
    headers: dict[str, Any] = field(default_factory=dict)
    content_type: str = "application/json" #

    @classmethod
    def from_str(cls, text: str, **kwargs) -> "Message":
        return cls(body=text.encode("utf-8"), **kwargs)

    @classmethod
    def from_dict(cls, data: dict, **kwargs) -> "Message":
        import json
        return cls(body=json.dumps(data).encode("utf-8"), **kwargs)


class Producer(ABC):

    @abstractmethod
    async def connect(self) -> None:
        """Установить соединение с брокером сообщений."""
        ...

    @abstractmethod
    async def disconnect(self) -> None:
        """Закрыть соединение с брокером сообщений."""
        ...

    @abstractmethod
    async def publish(
            self,
            destination: str,
            message: Message
    ) -> None:
        """
        Опубликовать сообщение.

        Args:
            message: Сообщение для отправки.
            destination: Имя назначения (имя очереди/топика/потока)
        """
        ...

    @abstractmethod
    async def is_connected(self) -> bool:
        """Проверить, активно ли соединение."""
        ...
