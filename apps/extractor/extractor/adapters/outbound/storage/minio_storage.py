from miniopy_async import Minio
from miniopy_async.error import S3Error
import logging
from extractor.application.ports.file_storage import FileStorage
from extractor.decorators.storage_decorator import handle_s3_errors
from typing import BinaryIO
from extractor.exceptions.storage_exp import StorageValidationError

logger = logging.getLogger(__name__)


class MinioStorage(FileStorage):

    def __init__(
            self,
            client: Minio,
            bucket_name: str,
    ):
        self._client = client
        self._bucket = bucket_name

    @handle_s3_errors
    async def upload_file(
            self,
            file_data: BinaryIO,
            object_name: str,
            content_type: str,
            metadata: dict | None = None
    ) -> str:

        """
        Метод для загрузки файла в Minio
        Args:
            file_data: Поток байтов файла
            object_name: Полное имя объекта в бакете (например, "archives/doc123.zip")
            content_type: MIME-тип файла
            metadata: Пользовательские метаданные

        Returns:
            object_name загруженного файла.
        """
        if not object_name or not object_name.strip():
            raise StorageValidationError("object_name не может быть пустым")

        # определяем размер
        file_size = self._get_file_size(file_data)
        if file_size == 0:
           raise StorageValidationError("Файл пустой, загрузка отменена")

        await self._client.put_object(
            self._bucket,
            object_name,
            data=file_data,
            length=file_size,
            content_type=content_type,
            metadata=metadata
        )

        logger.info("Файл '%s' успешно загружен в '%s'", object_name, self._bucket)
        return object_name


    @handle_s3_errors
    async def download_file(self, object_name: str) -> bytes:
        """Метод для загрузки файла из хранилища.
        Args:
            object_name: Полное имя объекта в бакете (например, "archives/doc123.zip")
        Returns:
            file_data: bytes загруженного файла.
        """
        if not object_name or not object_name.strip():
            raise StorageValidationError("object_name не может быть пустым")

        response = await self._client.get_object(
            self._bucket,
            object_name
        )
        try:

            data = await response.read()
            return data

        finally:
            response.close()
            await response.release()


    @handle_s3_errors
    async def exists(self, object_name: str) -> bool:
        """Метод для проверки на существование файла в хранилище.
        Args:
            object_name: Полное имя объекта в бакете (например, "archives/doc123.zip")

        Returns:
            exists: Булевое значение - существует ли файл в хранилище."""
        if not object_name or not object_name.strip():
            raise StorageValidationError("object_name не может быть пустым")

        try:
            await self._client.stat_object(
                self._bucket,
                object_name
            )
            logger.info(
                "Файл %s есть в хранилище", object_name
            )
            return True
        except S3Error as e:
            if e.code in ("NoSuchKey", "NoSuchObject"):
                logger.info("Файл %s отсутствует в хранилище", object_name)
                return False
            raise


    @handle_s3_errors
    async def delete_file(self, object_name: str) -> None:
        """Метод для удаления файла из хранилища.
        Args:
            object_name: Полное имя объекта в бакете (например, "archives/doc123.zip")
        Returns: None."""
        if not object_name or not object_name.strip():
            raise StorageValidationError("object_name не может быть пустым")

        await self._client.remove_object(
            self._bucket,
            object_name
        )
        logger.info("Объект %s удален из %s", object_name, self._bucket)


    @handle_s3_errors
    async def close(self) -> None:
        """Сброс connection pool."""
        if hasattr(self._client, '_http') and self._client._http:
            await self._client._http.close()


    @staticmethod
    def _get_file_size(file_data: BinaryIO) -> int:
        """Метод для получения размера файла.
        Args:
            file_data: Бинарный поток загруженного файла
        Returns:
            size: Целочисленное значение - размер переданного файла."""
        file_data.seek(0, 2)
        size = file_data.tell()
        file_data.seek(0)
        return size
