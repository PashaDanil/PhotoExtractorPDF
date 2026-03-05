import asyncio
from pathlib import Path
from uuid import uuid4
from extractor.adapters.outbound.pdf.pymupdf_reader import PyMuPDFReader
from extractor.adapters.outbound.cv.yolo_detector import YoloDetector
from extractor.adapters.outbound.image.fitz_image_cropper import FitzImageCropper
from extractor.adapters.outbound.zip.batched_zip_archiver import BatchedZipArchiver
from extractor.domain.entities.page import Page
from extractor.domain.entities.document import Document
from extractor.domain.entities.image import Image
from logger import setup_logger

logger = setup_logger(name="extractor", log_file="logs/app.log")


async def process_single_document(
        pdf_path: str,
        output_zip: str,
        pdf_reader: PyMuPDFReader,
        detector: YoloDetector,
        cropper: FitzImageCropper,
        archiver: BatchedZipArchiver,
) -> Document:
    """
    Обрабатывает один PDF документ:
    1. Открывает PDF
    2. Детектирует изображения на страницах (батчами)
    3. Вырезает найденные изображения
    4. Сохраняет в ZIP архив
    5. Возвращает Document с метаданными
    """
    print(f"Processing: {pdf_path}")

    # Генерируем ID для документа
    document_id = uuid4()
    file_name = Path(pdf_path).name

    # Список обработанных страниц
    processed_pages = []

    # Открываем документ
    async with pdf_reader.open_document(pdf_path) as doc_handle:
        page_count = doc_handle.page_count

        # Открываем ZIP для сохранения результатов
        async with archiver.open_zip(output_zip) as zip_handle:

            total_images = 0

            # Обрабатываем документ батчами по 8 страниц
            async for batch_rasters in doc_handle.iter_batches(batch_size=8, dpi=74):

                # Детектируем изображения на всём батче сразу
                batch_bboxes = await detector.detect_batch(batch_rasters)

                # Обрабатываем каждую страницу из батча
                for page_raster, bboxes in zip(batch_rasters, batch_bboxes):

                    # Генерируем ID для страницы
                    page_id = uuid4()

                    # Создаём объекты Image из bounding boxes
                    images = None
                    if bboxes:
                        images = [
                            Image(
                                id=uuid4(),
                                page_number=page_raster.page_num,
                                bounding_box=bbox,
                            )
                            for bbox in bboxes
                        ]

                    # Создаём объект Page
                    page = Page(
                        id=page_id,
                        document_id=document_id,
                        page_number=page_raster.page_num,
                        width=page_raster.width,
                        height=page_raster.height,
                        images=images,
                        render_dpi=page_raster.render_dpi,
                    )

                    processed_pages.append(page)

                    if not images:
                        print(f"  Page {page.page_number + 1}: No images found")
                        continue

                    # Получаем bytes страницы для вырезания
                    page_bytes = await doc_handle.get_page(page.page_number)

                    # Вырезаем и сохраняем изображения
                    count = await cropper.process_page(page_bytes, zip_handle, page)

                    total_images += count
                    print(f"  Page {page.page_number + 1}: Found and cropped {count} images")

            print(f"Total images extracted: {total_images}")

    # Создаём объект Document
    document = Document(
        id=document_id,
        file_name=file_name,
        pages=processed_pages,
        page_count=page_count,
        title=file_name,
    )

    print(f"Saved to: {output_zip}\n")

    return document


async def process_multiple_documents(documents: list[tuple[str, str]], device: str = 'cpu', model: str = "models/yolov11m-doclaynet.pt"):
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
        model_path=model,
        imgsz=960,
        classes=[6],  # только изображения
        device=device
    )

    # Прогреваем модель YOLO
    print("Warming up YOLO model...")
    detector.warmup()
    import time
    start_time = time.time()
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

    # === SHUTDOWN ===
    print("Shutting down services...")
    await pdf_reader.shutdown()
    detector.close()
    await cropper.shutdown()
    end_time = time.time()
    execution_time = end_time - start_time

    print(f"Время выполнения: {execution_time:.4f} секунд")
    print("Done!")


if __name__ == "__main__":

    logger.info("Приложение запущено")
    # указываешь документ + имя зип архива
    documents = [
        ("test_pdf/50137291M.pdf", "document1_images.zip"),
        ("test_pdf/Курсовая (Комаров Б9123-09.03.04).pdf", "document3_images.zip"),
    ]

    asyncio.run(process_multiple_documents
        (
        documents,
        device='cpu', # device = cuda for CUDA GPU
        model="models/yolov11m-doclaynet.pt"
    )
    )