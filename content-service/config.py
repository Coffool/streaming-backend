from pydantic_settings import BaseSettings
from pathlib import Path
from pydantic import Field


class Settings(BaseSettings):
    db_url: str = Field(alias="db_url_py")
    jwt_secret: str
    jwt_algorithm: str = "HS256"
    port: int = 8001
    content_base_path: str = "storage"

    # 🔹 Nueva variable para RabbitMQ
    rabbitmq_url: str = Field(default="amqp://guest:guest@rabbitmq:5672/")

    class Config:
        env_file = ".env"

    @property
    def storage_path(self) -> Path:
        return Path(self.content_base_path)


settings = Settings()  # type: ignore
