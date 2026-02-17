from dataclasses import dataclass
from uuid import UUID
from typing import List
from .page import Page

@dataclass
class Document:
    """Класс Документ"""
    id: UUID
    file_name: str
    pages: List[Page] | None

    page_count: int
    title: str

