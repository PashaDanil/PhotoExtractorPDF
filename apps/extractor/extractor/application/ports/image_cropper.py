# extractor/application/ports/image_cropper.py
from abc import ABC, abstractmethod
from collections.abc import AsyncIterator
from extractor.domain.entities.page import Page
from extractor.domain.entities.image import Image
from extractor.application.ports.zip_archiver import ZipArchiveHandle


class ImageCropper(ABC):
    """Интерфейс для обработки и вырезания изображений из PDF страниц"""

    @abstractmethod
    async def process_page(
            self,
            page_bytes: bytes,
            zip_handler: ZipArchiveHandle,
            page: Page,
    ) -> int:
        """
        Обрабатывает одну страницу PDF:
        1. Принимает bytes страницы
        2. Вырезает изображения согласно page.images
        3. Сохраняет их в zip через zip_handler

        Args:
            page_bytes: PDF страница в виде bytes
            zip_handler: Хэндлер для сохранения изображений
            page: Метаданные страницы с найденными изображениями

        Returns:
            Количество обработанных изображений
        """
        ...

    @abstractmethod
    async def iter_images(
            self,
            page_bytes: bytes,
            page: Page,
    ) -> AsyncIterator[tuple[Image, bytes]]:
        """
        Async-итератор по изображениям со страницы.

        Args:
            page_bytes: PDF страница в виде bytes
            page: Метаданные страницы с найденными изображениями

        Yields:
            Кортежи (метаданные_изображения, байты_pdf_изображения)
        """
        ...

    @abstractmethod
    async def shutdown(self) -> None:
        """
        Закрывает ресурсы при остановке сервера.
        Вызывается один раз при shutdown приложения.
        """
        ...