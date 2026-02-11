from dataclasses import dataclass
from typing import Self
from uuid import UUID, uuid4
from .bounding_box import BoundingBox

@dataclass(frozen=True, slots=True)
class Image:
    """
    Изображение, обнаруженное на странице PDF.

    Представляет результат работы детектора до извлечения.
    """
    id: UUID
    page_number: int
    bounding_box: BoundingBox
    confidence: float

    def __post_init__(self):
        """Валидация"""
        if self.page_number < 0:
            raise ValueError(f"page_number должен быть >= 0, получено: {self.page_number}")
        if (self.confidence < 0) or (self.confidence > 1):
            raise ValueError(f"confidence должен быть от 0 до 1")

    @classmethod
    def create(
        cls,
        page_number: int,
        bounding_box: BoundingBox,
        confidence: float,
    ) -> Self:
        """Создать с автогенерацией ID."""
        return cls(
            id=uuid4(),
            page_number=page_number,
            bounding_box=bounding_box,
            confidence=confidence,
        )

    def with_custom_bbox(self, padding: float) -> Self:
        """Создать копию с расширенным bbox."""
        return self.__class__(
            id=self.id,
            page_number=self.page_number,
            bounding_box=self.bounding_box.pad(padding),
            confidence=self.confidence,
            image_type=self.image_type
        )