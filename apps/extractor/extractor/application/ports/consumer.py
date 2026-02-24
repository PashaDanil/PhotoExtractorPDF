from abc import ABC, abstractmethod

class Consumer(ABC):

    @abstractmethod
    def connect(self):
        """
        Метод для подключения к брокеру сообщений
        """
        ...

    @abstractmethod
    def subscribe(self, topic: str):
        """
        Метод для подписки на топик
        Args:
            topic: Строка - название топика
        Returns:
            None
        """
        ...

    @abstractmethod
    def consume(self):
        """
        Начать считывать сообщения с брокера
        """
        ...

    @abstractmethod
    def close(self):
        """
        Остановка считывания сообщений с брокера
        """
        ...

