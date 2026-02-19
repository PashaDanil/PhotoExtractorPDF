class StorageError(Exception):
    """Базовая ошибка хранилища."""
    def __init__(self, message: str, cause: Exception | None = None):
        super().__init__(message)
        self.__cause__ = cause


class StorageConnectionError(StorageError):
    """Не удалось подключиться к хранилищу."""
    pass


class StorageOperationError(StorageError):
    """Ошибка при работе с хранилищем."""
    pass


class StorageObjectNotFoundError(StorageError):
    """Объект не найден в хранилище."""
    pass


class StorageQuotaExceededError(StorageError):
    """Превышен лимит хранилища."""
    pass


class StorageTimeoutError(StorageError):
    """Таймаут операции."""
    pass

class StorageAccessError(StorageError):
    """Ошибка доступа."""

class StorageValidationError(StorageError):
    """Некорректные входные данные для операции."""

class StorageUnexpectedError(StorageError):
    """Непредвиденная ошибка."""