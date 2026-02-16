import fitz
import asyncio
import numpy as np
from concurrent.futures import ThreadPoolExecutor
from contextlib import asynccontextmanager
from typing import AsyncIterator, Iterable, List
from itertools import islice
from extractor.domain.entities.page import PageRaster
from extractor.application.ports.pdf_reader import AsyncDocumentHandle, PDFReader


class PyMuPDFHandle(AsyncDocumentHandle):
    def __init__(self, pdf_doc: fitz.Document, executor: ThreadPoolExecutor):
        self._doc = pdf_doc
        self._executor = executor

    @property
    def page_count(self) -> int:
        return len(self._doc)

    def _render_page_sync(self, page_num: int, dpi: int = 72) -> PageRaster:
        """
        Синхронная функция для потока.
        Рендерит страницу сразу в Numpy массив (быстро).
        """
        page = self._doc[page_num]

        # Получаем pixmap (RGB)
        pix = page.get_pixmap(dpi=dpi, alpha=False)

        # Конвертируем буфер напрямую в numpy (Zero-copy)
        # 3 - channels (RGB)
        img_array = np.frombuffer(pix.samples, dtype=np.uint8).reshape(pix.h, pix.w, 3)

        # Важно: YOLO cv2 обычно хочет BGR, а fitz дает RGB.
        # Либо конвертим тут, либо указываем YOLO формат RGB.
        # Проще оставить RGB, ultralytics умеет с ним работать.
        # Это решим позже

        return PageRaster(
            page_num=page_num,
            width=pix.w,
            height=pix.h,
            image=img_array,
            render_dpi = dpi
        )

    async def iter_batches(
            self,
            batch_size: int = 8,
            dpi: int = 72,
            pages: Iterable[int] | None = None
    ) -> AsyncIterator[List[PageRaster]]:
        """
        Главный метод: выдает сразу ПАЧКУ страниц.
        """
        # возможно динамически высчитывать batch_size опираясь на размер и DPI
        # если будут переполнения памяти или близкие к этому ситуации
        loop = asyncio.get_running_loop()

        if pages is None:
            iterator = iter(range(self.page_count))
        else:
            iterator = iter(pages)

        while True:
            batch = list(islice(iterator, batch_size))
            if not batch:
                break

            tasks = [
                loop.run_in_executor(
                    self._executor,
                    self._render_page_sync,
                    idx,
                    dpi
                )
                for idx in batch
            ]


            batch_results = list(await asyncio.gather(*tasks))
            yield batch_results

    async def get_page(
            self,
            page_number: int,
            dpi: int = 72,
    ) -> bytes:
        """
        Async метод для получения страницы.
        Возвращает bytes страницы.
        """
        new_doc = fitz.open()
        new_doc.insert_pdf(self._doc, from_page=page_number, to_page=0)  # вставляем только страницу 0

        # Получаем bytes
        pdf_bytes = new_doc.tobytes()
        new_doc.close()
        return pdf_bytes



class PyMuPDFReader(PDFReader):
    def __init__(self, max_workers: int = 4):
        self._executor = ThreadPoolExecutor(max_workers=max_workers)

    @asynccontextmanager
    async def open_document(self, file_path: str) -> AsyncIterator[PyMuPDFHandle]:
        """
        Открывает документ.
        Принимаем путь к файлу, не загружает сразу весь в память.
        """
        pdf_doc = await asyncio.to_thread(fitz.open, file_path)

        try:
            yield PyMuPDFHandle(pdf_doc, self._executor)
        finally:
            await asyncio.to_thread(pdf_doc.close)

    async def shutdown(self):
        self._executor.shutdown(wait=True)
