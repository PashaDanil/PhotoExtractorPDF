from dataclasses import dataclass
from uuid import UUID
from .image import Image
from typing import List
import numpy as np
import numpy.typing as npt


@dataclass
class Page:
    """Класс Страница документа"""
    id: UUID
    document_id: UUID
    page_number: int
    width: float
    height: float
    images: List[Image] | None
    render_dpi: int = 72

    @property
    def width_px(self) -> float:
        """Ширина при текущем DPI"""
        return self.width *  self.render_dpi / 72

    @property
    def height_px(self) -> float:
        """Высота при текущем DPI"""
        return self.height * self.render_dpi / 72


@dataclass(frozen=True)
class PageRaster:
    """DTO изображения страницы"""
    page_num: int
    width: int
    height: int
    image: npt.NDArray[np.uint8]
    render_dpi: int = 72


