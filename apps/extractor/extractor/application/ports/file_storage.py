from abc import ABC, abstractmethod
from typing import BinaryIO


class FileStorage(ABC):
    """Интерфейс для взаимодействия с хранилищем файлов."""

    @abstractmethod
    async def upload_file(
        self,
        file_data: BinaryIO,
        object_name: str,
        content_type: str,
        metadata: dict | None = None
    ) -> str:
        ...

    @abstractmethod
    async def download_file(self, object_name: str) -> bytes:
        """Скачать файл из хранилища."""
        ...

    @abstractmethod
    async def delete_file(self, object_name: str) -> None:
        """Удалить файл из хранилища."""
        ...

    @abstractmethod
    async def exists(self, object_name: str) -> bool:
        """Проверить существование файла."""
        ...

    @abstractmethod
    async def close(self) -> None:
        """Освободить ресурсы."""
        ...
