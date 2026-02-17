import fitz
import asyncio
from collections.abc import AsyncIterator
from concurrent.futures import ThreadPoolExecutor

from extractor.application.ports.image_cropper import ImageCropper
from extractor.domain.entities.page import Page
from extractor.domain.entities.image import Image
from extractor.application.ports.zip_archiver import ZipArchiveHandle


class FitzImageCropper(ImageCropper):
    """
    Обработчик изображений на базе PyMuPDF (fitz).

    Thread-safe, stateless сервис для вырезания изображений из PDF.
    Предназначен для использования как singleton на уровне приложения.

    Один экземпляр обслуживает все запросы через общий пул потоков.
    """

    def __init__(self, max_workers: int | None = None):
        """
        Args:
            max_workers: Количество потоков для параллельной обработки.
                        None = автоматический выбор (обычно CPU count)
        """
        self._executor = ThreadPoolExecutor(max_workers=max_workers)

    @staticmethod
    def _crop_image_sync(
            page_bytes: bytes,
            image: Image,
            page_width: float,
            page_height: float,
            padding: float = 0.02
    ) -> bytes:
        # Открываем страницу из bytes
        pdf_doc = fitz.open(stream=page_bytes, filetype="pdf")

        try:
            page = pdf_doc[0]

            page_rect = page.rect
            pdf_page_width = page_rect.width  # в points (72 DPI)
            pdf_page_height = page_rect.height  # в points (72 DPI)

            if padding:
                coords = image.bounding_box.pad_relative(padding).to_absolute(
                    pdf_page_width,
                    pdf_page_height
                )

            clip_rect = fitz.Rect(coords)

            # Создаём новый документ
            new_doc = fitz.open()

            try:
                new_page = new_doc.new_page(
                    width=clip_rect.width,
                    height=clip_rect.height
                )

                new_page.show_pdf_page(
                    new_page.rect,
                    pdf_doc,
                    0,
                    clip=clip_rect
                )

                pdf_bytes = new_doc.tobytes()
                return pdf_bytes

            finally:
                new_doc.close()

        finally:
            pdf_doc.close()

    @staticmethod
    def _get_page_dimensions(page_bytes: bytes) -> tuple[float, float]:
        """
        Получает размеры страницы из PDF bytes.

        Args:
            page_bytes: PDF страница в виде bytes

        Returns:
            Кортеж (ширина, высота)
        """
        pdf_doc = fitz.open(stream=page_bytes, filetype="pdf")
        try:
            fitz_page = pdf_doc.load_page(0)
            rect = fitz_page.rect
            return rect.width, rect.height
        finally:
            pdf_doc.close()

    async def process_page(
            self,
            page_bytes: bytes,
            zip_handler: ZipArchiveHandle,
            page: Page,
    ) -> int:
        """
        Обрабатывает страницу: вырезает все изображения и сохраняет в zip.

        Thread-safe метод, может вызываться одновременно из разных запросов.

        Args:
            page_bytes: PDF страница в виде bytes
            zip_handler: Хэндлер для сохранения изображений
            page: Метаданные страницы с найденными изображениями

        Returns:
            Количество обработанных изображений
        """
        if not page.images:
            return 0

        # Получаем размеры страницы в отдельном потоке
        page_width, page_height = await asyncio.get_event_loop().run_in_executor(
            self._executor,
            self._get_page_dimensions,
            page_bytes,
        )

        # Обрабатываем изображения
        processed_count = 0
        for img_index, image in enumerate(page.images, start=1):
            # Вырезаем изображение в отдельном потоке
            image_bytes = await asyncio.get_event_loop().run_in_executor(
                self._executor,
                self._crop_image_sync,
                page_bytes,
                image,
                page_width,
                page_height,
            )

            # Формируем имя файла: p{номер_страницы}_i{номер_изображения}.pdf
            filename = f"p{page.page_number + 1}_i{img_index}.pdf"

            # Сохраняем в zip
            await zip_handler.add_file(filename, image_bytes)
            processed_count += 1

        return processed_count

    async def iter_images(
            self,
            page_bytes: bytes,
            page: Page,
    ) -> AsyncIterator[tuple[Image, bytes]]:
        """
        Async-итератор по изображениям со страницы.

        Thread-safe метод, может вызываться одновременно из разных запросов.

        Args:
            page_bytes: PDF страница в виде bytes
            page: Метаданные страницы с найденными изображениями

        Yields:
            Кортежи (метаданные_изображения, байты_pdf_изображения)
        """
        if not page.images:
            return

        # Получаем размеры страницы в отдельном потоке
        page_width, page_height = await asyncio.get_event_loop().run_in_executor(
            self._executor,
            self._get_page_dimensions,
            page_bytes,
        )

        # Итерируемся по изображениям
        for image in page.images:
            # Вырезаем изображение в отдельном потоке
            image_bytes = await asyncio.get_event_loop().run_in_executor(
                self._executor,
                self._crop_image_sync,
                page_bytes,
                image,
                page_width,
                page_height,
            )

            yield (image, image_bytes)

    async def shutdown(self) -> None:
        """
        Закрывает ThreadPoolExecutor при остановке сервера.
        Вызывается один раз при shutdown приложения.
        """
        self._executor.shutdown(wait=True)