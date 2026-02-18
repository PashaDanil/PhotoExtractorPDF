from extractor.application.ports.image_detector import ImageDetector
from typing import List
import numpy as np
from ultralytics import YOLO
from extractor.domain.entities.page import PageRaster
from extractor.domain.entities.bounding_box import BoundingBox
from concurrent.futures import ThreadPoolExecutor
import asyncio
import torch

class YoloDetector(ImageDetector):
    def __init__(
            self,
            model_path: str="models/yolov11m-doclaynet.pt",
            imgsz: int = 960,
            classes: List[int] = [6],
            device: str = "cuda"
    ):

        self.imgsz = imgsz
        self.classes = classes
        self.device = device
        self._model = YOLO(model_path)
        self._executor = ThreadPoolExecutor(max_workers=1)

    async def detect(self, page: PageRaster) -> List[BoundingBox] | None:
        loop = asyncio.get_running_loop()

        prediction = await loop.run_in_executor(
            self._executor,
            lambda: self._model.predict(page.image, imgsz=self.imgsz, classes=self.classes, device=self.device)
        )

        if prediction is None:
            return None

        norm_cord_tensor = prediction.boxes.xyxyn
        conf_tensor = prediction.boxes.conf.unsqueeze(1)
        prediction_boxes_data = torch.cat([norm_cord_tensor, conf_tensor], dim=1).cpu()

        prediction_boxes = [
            BoundingBox(
                x1=box[0],
                y1=box[1],
                x2=box[2],
                y2=box[3],
                confidence=box[4]
            )
            for box in prediction_boxes_data
        ]

        return prediction_boxes

    async def detect_batch(self, pages: List[PageRaster]) -> List[List[BoundingBox]] | None:
        loop = asyncio.get_running_loop()

        page_images = [page.image for page in pages]

        predictions = await loop.run_in_executor(
            self._executor,
            lambda: self._model.predict(page_images, imgsz=self.imgsz, classes=self.classes, device=self.device)
        )


        if predictions is None:
            return None
        
        batch_prediction_boxes_data = [
            torch.cat(
                [
                    prediction.boxes.xyxyn,
                    prediction.boxes.conf.unsqueeze(1)
                ],
                dim=1
            ).cpu()
            for prediction in predictions
        ]

        width = 1
        height = 1
        batch_prediction_boxes = [
            [
                BoundingBox(
                    x1=box[0],
                    y1=box[1],
                    x2=box[2],
                    y2=box[3],
                    confidence=box[4]
                )
                for box in prediction_boxes_data
            ]
            for prediction_boxes_data in batch_prediction_boxes_data
        ]

        return batch_prediction_boxes

    def warmup(self) -> None:
        """
        Прогрев модели фиктивным изображением.
        Это инициализирует CUDA, выделяет память и компилирует графы.
        """
        dummy_image = np.zeros((self.imgsz, self.imgsz, 3), dtype=np.uint8)
        for _ in range(3):
            self._model.predict(source=dummy_image, imgsz=self.imgsz, classes=self.classes, device=self.device, verbose=False)

    def close(self):
        self._executor.shutdown(wait=True)
