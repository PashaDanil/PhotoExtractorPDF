from dataclasses import dataclass
from typing import Self


@dataclass(frozen=True, slots=True)
class BoundingBox:
    """
    Ограничивающая рамка изображения.

    Использует координатную систему (x1, y1, x2, y2),
    где (x1, y1) — левый верхний угол, (x2, y2) — правый нижний.

    Координаты нормализованы к размеру страницы [0.0, 1.0].
    """
    x1: float
    y1: float
    x2: float
    y2: float

    def __post_init__(self) -> None:
        """Валидация координат."""
        if self.x1 > self.x2:
            raise ValueError(f"x1 ({self.x1}) не может быть больше x2 ({self.x2})")
        if self.y1 > self.y2:
            raise ValueError(f"y1 ({self.y1}) не может быть больше y2 ({self.y2})")

    @property
    def width(self) -> float:
        return self.x2 - self.x1

    @property
    def height(self) -> float:
        return self.y2 - self.y1

    @property
    def center(self) -> tuple[float, float]:
        return (
            (self.x1 + self.x2) / 2,
            (self.y1 + self.y2) / 2
        )

    def as_xywh(self) -> tuple[float, float, float, float]:
        """Вернуть как (x, y, width, height)."""
        return (self.x1, self.y1, self.width, self.height)

    def as_xyxy(self) -> tuple[float, float, float, float]:
        """Вернуть как (x1, y1, x2, y2)."""
        return (self.x1, self.y1, self.x2, self.y2)

    @classmethod
    def from_xywh(cls, x: float, y: float, width: float, height: float) -> Self:
        """Создать из формата (x, y, width, height)."""
        return cls(
            x1=x,
            y1=y,
            x2=x + width,
            y2=y + height
        )

    @classmethod
    def from_center(cls, cx: float, cy: float, width: float, height: float) -> Self:
        """Создать из центра и размеров."""
        half_w = width / 2
        half_h = height / 2
        return cls(
            x1=cx - half_w,
            y1=cy - half_h,
            x2=cx + half_w,
            y2=cy + half_h
        )

    def scale(self, factor: float) -> Self:
        """Масштабировать от центра."""
        cx, cy = self.center
        new_w = self.width * factor
        new_h = self.height * factor
        return self.from_center(cx, cy, new_w, new_h)

    def pad(self, padding: float) -> Self:
        """Добавить отступ со всех сторон."""
        return self.__class__(
            x1=max(0.0, self.x1 - padding),
            y1=max(0.0, self.y1 - padding),
            x2=min(1.0, self.x2 + padding),
            y2=min(1.0, self.y2 + padding)
        )

    def padx(self, padding_left: float = 0, padding_right: float = 0) -> Self:
        """Добавить отступ по бокам."""
        return self.__class__(
            x1=max(0.0, self.x1 - padding_left),
            y1=self.y1,
            x2=min(1.0, self.x2 + padding_right),
            y2=self.y2
        )

    def pady(self, padding_top: float = 0, padding_bottom: float = 0) -> Self:
        """Добавить отступ сверху и снизу."""
        return self.__class__(
            x1=self.x1,
            y1=max(0.0, self.y1 - padding_top),
            x2=self.x2,
            y2=min(1.0, self.y2 + padding_bottom)
        )

    def to_absolute(self, page_width: int, page_height: int) -> tuple[float, float, float, float]:
        """
        Конвертировать нормализованные [0,1] в абсолютные пиксели.
        Возвращает координаты в формате (x1, y1, x2, y2).
        """
        return (
            self.x1 * page_width,
            self.y1 * page_height,
            self.x2 * page_width,
            self.y2 * page_height
        )

    def intersects(self, other: Self) -> bool:
        """Проверить пересечение с другим bbox."""
        return not (
            self.x2 <= other.x1 or
            self.x1 >= other.x2 or
            self.y2 <= other.y1 or
            self.y1 >= other.y2
        )
