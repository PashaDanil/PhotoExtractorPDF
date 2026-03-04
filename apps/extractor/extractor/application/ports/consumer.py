from abc import ABC, abstractmethod

class Consumer(ABC):

    @abstractmethod
    async def connect(self):
        """
        Метод для подключения к брокеру сообщений
        """
        ...

    @abstractmethod
    async def subscribe(self, topic: str):
        """
        Метод для подписки на топик
        Args:
            topic: Строка - название топика
        Returns:
            None
        """
        ...

    @abstractmethod
    async def consume(self):
        """
        Начать считывать сообщения с брокера
        """
        ...

    @abstractmethod
    async def disconnect(self):
        """
        Остановка считывания сообщений с брокера
        """
        ...

