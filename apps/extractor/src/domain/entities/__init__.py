# adapters/ocr/__init__.py

from .image import BoundingBox, Image
from .document import Page, Document

__all__ = ["BoundingBox", "Image", "Page", "Document"]