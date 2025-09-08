# song_repository.py
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.future import select
from sqlalchemy.orm import selectinload
from database.models import Song, Album, Artist, User


class SongRepository:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def get_by_title_ilike(self, query: str, limit: int, offset: int):
        # Solo columnas necesarias: id, title, duration, audio_url, album_id
        stmt = (
            select(Song.id, Song.title, Song.duration, Song.audio_url, Song.album_id)
            .where(Song.title.ilike(f"%{query}%"))
            .limit(limit)
            .offset(offset)
        )
        result = await self.session.execute(stmt)
        return result.all()  # devuelve lista de tuplas
