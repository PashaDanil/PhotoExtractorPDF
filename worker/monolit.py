from ultralytics import YOLO
import fitz
import pymupdf
import numpy as np
import polars as pl
from typing import List
import os
import time

class Image:
    """Класс будет хранить
    1. Координаты объекта на странице
    2. Номер страницы
    """

class Page:
    def __init__(self, page_dpf : pymupdf.Page, page_num : int, parent_doc : pymupdf.Document, parent_path : str):
        self.objects_df = None
        self.page_num = page_num
        self.page_pdf = page_dpf
        self.parent_doc = parent_doc
        self.parent_path = parent_path

    def set_objects_df(self, df : pl.DataFrame):
        self.objects_df = df

    def get_image_boxes(self, confidence_rate = 0.3) -> pl.DataFrame:
        return self.objects_df.filter(
            (pl.col('class') == 6)
            &
            (pl.col('confidence') > confidence_rate)
        )

    def convert_to_array(self, dpi=72) -> np.ndarray | None:
        try:
            pix = self.page_pdf.get_pixmap(dpi=dpi)

            img = np.frombuffer(pix.samples, dtype=np.uint8).reshape(
                pix.height, pix.width, pix.n
            )
            return img
        except Exception as e:
            print(f'Error while converting pdf page to numpy.ndarray: {e}')
            return None

    def clip_all_images(self) -> bool:
        if self.objects_df is None:
            return False

        output_folder = self.parent_path[:-5]
        if not os.path.exists(output_folder):
            os.mkdir(output_folder)
        for index, box in enumerate(self.get_image_boxes()['box']):
            output_path = output_folder + '/p' + str(self.page_num + 1) + '_n' + str(index + 1) + '.pdf'
            self.clip_page(output_path, list(box.values()))
        return True

    def clip_page(self, output_path, bbox):

        # Определяем область для вырезания
        clip_rect = fitz.Rect(bbox)
        # Создаём новый документ
        new_doc = fitz.open()

        # Создаём страницу с размером вырезаемой области
        new_page = new_doc.new_page(
            width=clip_rect.width,
            height=clip_rect.height
        )

        # Копируем содержимое из исходной области
        new_page.show_pdf_page(
            new_page.rect,  # куда вставить (вся новая страница)
            self.parent_doc,  # исходный документ
            self.page_num,  # номер страницы
            clip=clip_rect  # область для вырезания
        )

        new_doc.save(output_path)
        new_doc.close()

        print(f"Сохранено: {output_path}")
        print(f"Сохранено: {output_path}")



class Document:
    """Добавить ориентацию страницы"""
    def __init__(self, document_path : str) -> None:
        self.document_path = document_path
        self.document = None
        self.number_of_pages = 0
        self.current_page = 0
        try:
            self.document = fitz.open(document_path)
            self.number_of_pages = len(self.document)
        except Exception as e:
            print(f"Error: {e}")

    def get_number_of_pages(self) -> int:
        return self.number_of_pages

    def get_path(self) -> str:
        return self.document_path

    def get_all_pages(self) -> List[Page] | None:
        if self.document is not None:
            return [Page(page, page_num, self.document, self.document_path) for page_num, page in enumerate(self.document)]
        return None

    def get_all_pages_as_images(self) -> List[np.ndarray]:
        pages = self.get_all_pages()
        images = [page.convert_to_array() for page in pages]
        return images

    def get_page(self) -> Page | None:
        if self.current_page < self.number_of_pages:
            page = Page(self.document[self.current_page], self.current_page, self.document, self.document_path)
            self.current_page += 1
            return page
        return None

    def close(self):
        self.document.close()



class OCRModel:

    def __init__(self, model_name="yolov11m-doclaynet.pt"):
        self.model = YOLO(model_name)

    # def detect_objects(self, document: Document, pages: list = None):
    #     if list is not None:
    #
    #     # if pages:
    #     #     iter_obj = pages
    #     # else:
    #     #     iter = range(0, Document.get_number_of_pages())
    #     # for pag_num in iter:
    #     page = Document.get_page()
    #     while page:
    #         page = Document.get_page()

    def process_all_pages(self, document: Document):
        all_pages = document.get_all_pages()
        images = document.get_all_pages_as_images()
        images_detection_result = self.model.predict(images)

        for page, detection_result in zip(all_pages, images_detection_result):
            page.set_objects_df(detection_result.to_df())
            page.clip_all_images()
    # def process_page_detection(self, page: Page, dpi=72):
    #     imgs =
    #     result = self.model.predict(img)[0]


if __name__ == "__main__":
    ocr = OCRModel()
    doc = Document("50137291M.pdf")
    try:
        start = time.time()
        ocr.process_all_pages(doc)
        end = time.time()

        print(f"Время выполнения: {end - start:.4f} секунд")
    finally:
        doc.close()