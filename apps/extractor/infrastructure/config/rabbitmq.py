from pydantic_settings import BaseSettings, SettingsConfigDict
from pydantic import Field

class RabbitMQConfig(BaseSettings):
    host: str = Field("localhost", alias="RABBITMQ_HOST")
    port: int = Field(5672, alias="RABBITMQ_PORT")
    username: str = Field("guest", alias="RABBITMQ_USERNAME")
    password: str = Field("guest", alias="RABBITMQ_PASSWORD")
    virtual_host: str = Field("/", alias="RABBITMQ_VIRTUAL_HOST")

    consume_queue: str = Field(..., alias="RABBITMQ_CONSUME_QUEUE")
    publish_queue: str = Field(..., alias="RABBITMQ_PUBLISH_QUEUE")

    prefetch_count: int = Field(1, alias="RABBITMQ_PREFETCH_COUNT")

    exchange: str = Field("", alias="RABBITMQ_EXCHANGE")
    exchange_type: str = Field("direct", alias="RABBITMQ_EXCHANGE_TYPE")
    durable: bool = Field(True, alias="RABBITMQ_DURABLE")
    delivery_mode: int = Field(2, alias="RABBITMQ_DELIVERY_MODE")

    @property
    def url(self) -> str:
        return f"amqp://{self.username}:{self.password}@{self.host}:{self.port}/{self.virtual_host}"

    model_config = SettingsConfigDict(
        env_file=".env",
        populate_by_name=True,
    )