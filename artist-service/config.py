# config.py
from pydantic_settings import BaseSettings
from pathlib import Path


class Settings(BaseSettings):
    db_url: str
    jwt_secret: str
    jwt_algorithm: str = "HS256"
    port: int = 8001
    content_base_path: str = "storage"  # 🔹 NUEVO: path base para el storage

    class Config:
        env_file = ".env"

    @property
    def storage_path(self) -> Path:
        """Retorna el Path completo del directorio de storage"""
        return Path(self.content_base_path)


settings = Settings()  # type: ignore
