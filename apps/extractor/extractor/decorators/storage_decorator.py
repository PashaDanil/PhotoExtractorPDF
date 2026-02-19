import functools
import logging
import inspect
import asyncio
from miniopy_async.error import S3Error
from aiohttp import ClientConnectionError, ServerTimeoutError
from extractor.exceptions.storage_exp import (StorageError,
                                              StorageConnectionError,
                                              StorageOperationError ,
                                              StorageObjectNotFoundError,
                                              StorageQuotaExceededError,
                                              StorageTimeoutError,
                                              StorageAccessError,
                                              StorageValidationError,
                                              StorageUnexpectedError)

logger = logging.getLogger(__name__)

def handle_s3_errors(func):
    sig = inspect.signature(func)

    @functools.wraps(func)
    async def wrapper(self, *args, **kwargs):
        # берём object_name из аргументов
        operation = func.__name__
        # корректно достаём object_name из сигнатуры
        bound = sig.bind(self, *args, **kwargs)
        bound.apply_defaults()
        object_name = bound.arguments.get("object_name", "unknown")

        try:
            return await func(self, *args, **kwargs)

        except S3Error as e:
            if e.code == "NoSuchBucket":
                logger.error("Бакет %s не существует", self._bucket)
                raise StorageObjectNotFoundError(
                    f"Бакет '{self._bucket}' не найден", cause=e
                )
            if e.code == "AccessDenied":
                logger.error("Нет доступа к бакету %s", self._bucket)
                raise StorageAccessError(
                    f"Нет прав на доступ к '{self._bucket}'", cause=e
                )
            if e.code == "EntityTooLarge":
                raise StorageQuotaExceededError(
                    f"Файл слишком большой", cause=e
                )

            logger.exception("S3 ошибка при операции '%s' для %s", operation, object_name)
            raise StorageOperationError(
                f"S3 ошибка при '{operation}' для '{object_name}': {e.code}", cause=e
            )

        except (StorageError, StorageValidationError):
            raise # проброс кастомых ошибок

        except (ClientConnectionError, ConnectionError, OSError) as e:
            logger.exception("Нет связи с хранилищем при '%s'", operation)
            raise StorageConnectionError(
                "Не удалось подключиться к хранилищу", cause=e
            )

        except (asyncio.TimeoutError, ServerTimeoutError) as e:
            logger.exception("Таймаут при '%s' для %s", operation, object_name)
            raise StorageTimeoutError(
                f"Таймаут при '{operation}' для '{object_name}'", cause=e
            )

        except Exception as e:
            logger.exception("Непредвиденная ошибка при '%s' для %s", operation, object_name)
            raise StorageUnexpectedError(
                f"Непредвиденная ошибка при '{operation}' для '{object_name}'", cause=e
            )

    return wrapper