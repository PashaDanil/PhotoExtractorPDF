from typing import BinaryIO, AsyncIterator
from contextlib import asynccontextmanager
from extractor.application.ports.zip_archiver import ZipArchiveHandle, ZipArchiver
import zipstream
import asyncio
import os


class BatchedZipArchiveHandle(ZipArchiveHandle):
    """
    Обработчик zip-архива с батчевой записью.
    Запись в файл идёт в пуле потоков через asyncio.to_thread.
    """

    def __init__(self, file_handle: BinaryIO, batch_size: int = 2 * 1024 * 1024):
        self._file = file_handle
        self._batch_size = batch_size
        self._buffer = bytearray()
        self._zs = zipstream.ZipFile(mode='w', compression=zipstream.ZIP_DEFLATED)
        self._files_added: list[tuple[str, bytes]] = []
        self._closed = False

    async def add_file(self, filename: str, file_data: bytes) -> None:
        if self._closed:
            raise RuntimeError("Archive is closed")
        self._files_added.append((filename, file_data))

    async def _flush_buffer(self) -> None:
        """Сброс буфера на диск в отдельном потоке."""
        if not self._buffer:
            return

        data = bytes(self._buffer)
        self._buffer.clear()

        # запись в файл — синхронная операция, уводим её в пул потоков
        await asyncio.to_thread(self._file.write, data)

    async def _process_zip_stream(self) -> None:
        """Генерация zip-стрима и батчевая запись в файл."""

        # Добавляем все файлы в zipstream
        for filename, file_data in self._files_added:
            self._zs.write_iter(filename, iter([file_data]))

        # Генерируем архивацию и пишем батчами
        for chunk in self._zs:
            self._buffer.extend(chunk)

            if len(self._buffer) >= self._batch_size:
                await self._flush_buffer()
                # даём другим корутинам шанс выполниться
                await asyncio.sleep(0)

        # Финальный сброс
        await self._flush_buffer()

    async def close(self) -> None:
        """Финализирует архив и закрывает файл."""
        if self._closed:
            return

        await self._process_zip_stream()

        # Закрытие файла — тоже в отдельном потоке
        await asyncio.to_thread(self._file.close)

        self._closed = True


class BatchedZipArchiver(ZipArchiver):
    """
    Архиватор, создающий BatchedZipArchiveHandle.
    """

    def __init__(
        self,
        batch_size: int = 2 * 1024 * 1024,
        base_path: str = "",
    ):
        self._batch_size = batch_size
        self._base_path = base_path

    @asynccontextmanager
    async def open_zip(self, zip_name: str) -> AsyncIterator[ZipArchiveHandle]:
        """Открывает ZIP-архив для записи (асинхронный контекстный менеджер)."""

        if self._base_path:
            full_path = os.path.join(self._base_path, zip_name)
        else:
            full_path = zip_name

        # Обычный синхронный open
        file_handle = open(full_path, 'wb')

        handle = BatchedZipArchiveHandle(
            file_handle=file_handle,
            batch_size=self._batch_size,
        )

        try:
            yield handle
        finally:
            await handle.close()