from abc import ABC, abstractmethod
from typing import List
from apps.extractor.extractor.domain.entities.page import PageRaster
from apps.extractor.extractor.domain.entities.bounding_box import BoundingBox

class ImageDetector(ABC):
    """Детектор изображений на странице документа"""

    @abstractmethod
    def detect(self, page: PageRaster) -> List[BoundingBox] | None:
        """
        Находит изображения на одной странице.

        Args:
            page: страница с растеризованным изображением

        Returns:
            список bounding boxes с координатами и confidence
        """
        ...

    @abstractmethod
    def detect_batch(self, pages: List[PageRaster]) -> List[List[BoundingBox]] | None:
        """
        Находит изображения на батче страниц.

        Args:
            pages: список страниц

        Returns:
            список со списком bounding boxes с координатами
            и confidence для каждого изображения
        """
        ...

    def warmup(self) -> None:
        """
        Прогрев модели (первый inference).
        """
        pass