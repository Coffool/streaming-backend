import os
from pathlib import Path
from fastapi import UploadFile
from config import settings
from typing import Union, Any


class FileUploader:
    """Helper para manejar la subida de archivos usando configuración centralizada"""

    @staticmethod
    async def upload_profile_picture(
        file: UploadFile,
        artist_id: Union[int, Any],  # 🔹 Acepta tanto int como objetos SQLAlchemy
    ) -> str:
        """
        Sube la foto de perfil del artista a {CONTENT_BASE_PATH}/{artist_id}/utils/
        """
        # Usar path centralizado desde config
        storage_path = settings.storage_path

        # 🔹 Convertir artist_id a string (funciona con int y Column)
        artist_id_str = str(artist_id)

        # Carpeta utils del artista
        artist_utils_folder = storage_path / artist_id_str / "utils"
        artist_utils_folder.mkdir(parents=True, exist_ok=True)

        # Asegurar que filename nunca sea None
        original_filename = file.filename or "file.jpg"

        # Obtener la extensión original
        _, ext = os.path.splitext(original_filename)
        if not ext:
            ext = ".jpg"  # fallback si no tiene extensión

        # Nombre fijo del archivo
        filename = f"profile_picture{ext}"
        file_location = artist_utils_folder / filename

        # Guardar (sobre-escribe si ya existía)
        with open(file_location, "wb") as buffer:
            buffer.write(await file.read())

        # Retornar la URL relativa usando el path base de config
        return f"/{artist_id_str}/utils/{filename}"
