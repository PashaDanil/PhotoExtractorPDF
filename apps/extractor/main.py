
import asyncio
from pathlib import Path
import uuid
import numpy as np
from typing import List

from extractor.adapters.outbound.pdf.pymupdf_reader import PyMuPDFReader
from extractor.adapters.outbound.cv.yolo_detector import YoloDetector
from extractor.adapters.outbound.image.fitz_image_cropper import FitzImageCropper
from extractor.adapters.outbound.zip.batched_zip_archiver import BatchedZipArchiver
from extractor.domain.entities.page import Page

async def process_single_document(
        pdf_path: str,
        output_zip: str,
        pdf_reader: PyMuPDFReader,
        detector: YoloDetector,
        cropper: FitzImageCropper,
        archiver: BatchedZipArchiver,
):
    """
    Обрабатывает один PDF документ:
    1. Открывает PDF
    2. Детектирует изображения на страницах (батчами)
    3. Вырезает найденные изображения
    4. Сохраняет в ZIP архив
    """
    print(f"Processing: {pdf_path}")

    # Открываем документ
    async with pdf_reader.open_document(pdf_path) as doc_handle:
        # Открываем ZIP для сохранения результатов
        async with archiver.open_zip(output_zip) as zip_handle:

            total_images = 0

            # Обрабатываем документ батчами по 8 страниц
            async for batch_rasters in doc_handle.iter_batches(batch_size=8, dpi=150):

                # Детектируем изображения на всём батче сразу
                batch_bboxes = await detector.detect_batch(batch_rasters)

                # Обрабатываем каждую страницу из батча
                for page_raster, bboxes in zip(batch_rasters, batch_bboxes):

                    if not bboxes:
                        print(f"  Page {page_raster.page_num + 1}: No images found")
                        continue

                    # Создаём объект Page с найденными изображениями
                    page = Page(
                        id=uuid.uuid4(),
                        document_id=uuid.uuid4(),
                        page_number=1,
                        width=800.0,
                        height=600.0,
                        image=dummy_image,
                        images=None  # или пустой список []
                    )
                    page = Page(
                        page_raster,
                        bboxes)

                    # Получаем bytes страницы для вырезания
                    page_bytes = await doc_handle.get_page(page.page_number)

                    # Вырезаем и сохраняем изображения
                    count = await cropper.process_page(page_bytes, zip_handle, page)

                    total_images += count
                    print(f"  Page {page.page_number + 1}: Found and cropped {count} images")

            print(f"Total images extracted: {total_images}")

    print(f"Saved to: {output_zip}\n")


async def process_multiple_documents():
    """
    Пример обработки нескольких документов одновременно.
    Все сервисы - singleton, используются всеми документами.
    """

    # === ИНИЦИАЛИЗАЦИЯ SINGLETON СЕРВИСОВ ===
    print("Initializing services...")

    # PDF Reader - один на все документы
    pdf_reader = PyMuPDFReader(max_workers=4)

    # Image Detector - один на все документы
    detector = YoloDetector(
        model_path="yolov11m-doclaynet.pt",
        imgsz=960,
        classes=[6],  # только изображения
        device="cuda"
    )

    # Прогреваем модель YOLO
    print("Warming up YOLO model...")
    detector.warmup()

    # Image Cropper - один на все документы
    cropper = FitzImageCropper(max_workers=2)

    # ZIP Archiver - один на все документы
    archiver = BatchedZipArchiver(
        batch_size=2 * 1024 * 1024,  # 2MB батчи
        base_path="./output"
    )

    # Создаём выходную директорию
    Path("./output").mkdir(exist_ok=True)

    print("Services initialized!\n")

    # === ОБРАБОТКА ДОКУМЕНТОВ ===

    # Список документов для обработки
    documents = [
        ("50137291M.pdf", "document1_images.zip"),
        ("EN_DualSense_Wireless_Controller_CFI-ZCT1W_Instruction_Man_Web.pdf", "document2_images.zip"),
    ]

    # Вариант 1: Последовательная обработка
    print("=== Sequential Processing ===")
    for pdf_path, output_zip in documents:
        await process_single_document(
            pdf_path,
            output_zip,
            pdf_reader,
            detector,
            cropper,
            archiver,
        )

    # Вариант 2: Параллельная обработка (если достаточно памяти)
    # print("=== Parallel Processing ===")
    # tasks = [
    #     process_single_document(
    #         pdf_path,
    #         output_zip,
    #         pdf_reader,
    #         detector,
    #         cropper,
    #         archiver,
    #     )
    #     for pdf_path, output_zip in documents
    # ]
    # await asyncio.gather(*tasks)

    # === SHUTDOWN ===
    print("Shutting down services...")
    await pdf_reader.shutdown()
    detector.close()
    await cropper.shutdown()
    print("Done!")


async def process_with_iterator_example():
    """
    Пример использования iter_images() для более гибкого контроля.
    """

    pdf_reader = PyMuPDFReader(max_workers=4)
    detector = YoloDetector()
    cropper = FitzImageCropper(max_workers=4)

    detector.warmup()

    try:
        async with pdf_reader.open_document("document.pdf") as doc_handle:

            # Обрабатываем только первые 10 страниц
            async for batch_rasters in doc_handle.iter_batches(
                    batch_size=5,
                    dpi=150,
                    pages=range(10)  # только первые 10 страниц
            ):
                batch_bboxes = await detector.detect_batch(batch_rasters)

                for page_raster, bboxes in zip(batch_rasters, batch_bboxes):
                    if not bboxes:
                        continue

                    page = Page.from_raster(page_raster, bboxes)
                    page_bytes = await doc_handle.get_page(page.page_number)

                    # Используем итератор для контроля каждого изображения
                    async for image, image_bytes in cropper.iter_images(page_bytes, page):
                        # Можем делать что-то с каждым изображением
                        print(f"Cropped image with confidence: {image.confidence}")

                        # Сохраняем в отдельные файлы вместо ZIP
                        filename = f"page_{page.page_number}_img_{image}.pdf"
                        with open(filename, "wb") as f:
                            f.write(image_bytes)

    finally:
        await pdf_reader.shutdown()
        detector.close()
        await cropper.shutdown()


async def fastapi_like_example():
    """
    Пример интеграции в FastAPI-подобное приложение.
    """
    from dataclasses import dataclass

    @dataclass
    class AppState:
        """Глобальное состояние приложения (singleton сервисы)"""
        pdf_reader: PyMuPDFReader
        detector: YoloDetector
        cropper: FitzImageCropper
        archiver: BatchedZipArchiver

    # === STARTUP ===
    async def startup():
        print("Starting application...")

        pdf_reader = PyMuPDFReader(max_workers=4)
        detector = YoloDetector(device="cuda:0")
        detector.warmup()
        cropper = FitzImageCropper(max_workers=8)
        archiver = BatchedZipArchiver(base_path="./output")

        return AppState(
            pdf_reader=pdf_reader,
            detector=detector,
            cropper=cropper,
            archiver=archiver,
        )

    # === REQUEST HANDLER ===
    async def handle_request(app_state: AppState, pdf_path: str, output_zip: str):
        """
        Обработчик запроса.
        Все сервисы переиспользуются между запросами (thread-safe).
        """
        await process_single_document(
            pdf_path,
            output_zip,
            app_state.pdf_reader,
            app_state.detector,
            app_state.cropper,
            app_state.archiver,
        )

    # === SHUTDOWN ===
    async def shutdown(app_state: AppState):
        print("Shutting down application...")
        await app_state.pdf_reader.shutdown()
        app_state.detector.close()
        await app_state.cropper.shutdown()

    # Симуляция работы приложения
    app_state = await startup()

    try:
        # Симулируем несколько входящих запросов
        await handle_request(app_state, "doc1.pdf", "doc1.zip")
        await handle_request(app_state, "doc2.pdf", "doc2.zip")
        await handle_request(app_state, "doc3.pdf", "doc3.zip")
    finally:
        await shutdown(app_state)


if __name__ == "__main__":
    import torch
    import sys

    print(f"Python version: {sys.version}")
    print(f"PyTorch version: {torch.__version__}")
    print(f"CUDA available: {torch.cuda.is_available()}")
    print(f"CUDA device count: {torch.cuda.device_count()}")
    if torch.cuda.is_available():
        print(f"CUDA device name: {torch.cuda.get_device_name(0)}")
    # Выберите нужный пример:

    # Основной пример - обработка нескольких документов
    asyncio.run(process_multiple_documents())

    # Или пример с итератором
    # asyncio.run(process_with_iterator_example())

    # Или пример FastAPI-like приложения
    # asyncio.run(fastapi_like_example())