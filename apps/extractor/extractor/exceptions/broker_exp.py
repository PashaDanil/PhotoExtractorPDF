class BrokerError(Exception):
    """Базовая ошибка брокера."""
    def __init__(self, message: str, cause: Exception | None = None):
        super().__init__(message)
        self.__cause__ = cause

class BrokerConnectionError(BrokerError):
    """Не удалось подключиться к брокеру."""
    pass

class BrokerPublishError(BrokerError):
    """Не удалось опубликовать сообщение"""
    pass

class BrokerConsumeError(BrokerError):
    """Ошибка при получении/обработке сообщения"""
    pass

class BrokerValidationError(BrokerError):
    """Некорректные входные данные для операции."""

class BrokerFatalError(BrokerError):
    """Логическая/конфигурационная ошибка"""

class BrokerUnexpectedError(BrokerError):
    """Непредвиденная ошибка."""