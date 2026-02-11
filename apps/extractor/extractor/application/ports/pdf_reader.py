from abc import ABC, abstractmethod
from contextlib import asynccontextmanager
from typing import AsyncIterator, List, Iterable
from src.entities.page import PageRaster


class AsyncDocumentHandle(ABC):
    """Хэндлер открытого документа"""
    @property
    @abstractmethod
    def page_count(self) -> int: ...

    @abstractmethod
    def _render_page_sync(
            self,
            page_num: int,
            dpi: int = 72
    ) -> PageRaster:
        """
        Синхронная функция для потока.
        Рендерит страницу в Numpy массив.
        """
        ...

    @abstractmethod
    async def iter_batches(
            self,
            batch_size: int = 8,
            dpi: int = 72,
            pages: Iterable[int] | None = None,
    ) -> AsyncIterator[List[PageRaster]]:
        """
        Async итератор по DPF страницам.
        Возвращает батч страниц.
        """
        ...


class PDFReader(ABC):

    @abstractmethod
    @asynccontextmanager
    async def open_document(self, file_path: str) -> AsyncIterator[AsyncDocumentHandle]:
        """
        Async context manager для работы с PDF
        Возвращает асинхронный hadler с итератором по документу
        """
        ...

    @abstractmethod
    async def shutdown(self) -> None:
        """Завершение работы пулла потоков"""
        ...