# album_repository.py
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.future import select
from database.models import Album


class AlbumRepository:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def get_by_title_ilike(self, query: str, limit: int, offset: int):
        # Solo columnas necesarias: id, title, cover_url
        stmt = (
            select(Album.id, Album.title, Album.cover_url)
            .where(Album.title.ilike(f"%{query}%"))
            .limit(limit)
            .offset(offset)
        )
        result = await self.session.execute(stmt)
        return result.all()
