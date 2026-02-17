from abc import ABC, abstractmethod
from typing import AsyncIterator, AsyncContextManager


class ZipArchiveHandle(ABC):

    @abstractmethod
    async def add_file(
        self,
        filename: str,
        file_data: bytes,
    ) -> None:
        ...

    @abstractmethod
    async def close(self) -> None:
        ...


class ZipArchiver(ABC):

    @abstractmethod
    def open_zip(
        self,
        zip_name: str,
    ) -> AsyncContextManager[ZipArchiveHandle]:
        """Возвращает async context manager для работы с zip-архивом"""
        ...