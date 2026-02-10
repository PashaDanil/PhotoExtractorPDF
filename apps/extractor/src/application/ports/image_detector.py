from abc import ABC, abstractmethod
from contextlib import asynccontextmanager
from typing import AsyncIterator, List, Iterable
from dataclasses import dataclass
import numpy as np
import numpy.typing as npt
from pdf_reader import PageRaster

class CVModel(ABC):
    @abstractmethod
    def __init__(self, model_path: str) -> None:
        ...

