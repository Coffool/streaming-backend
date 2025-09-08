from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from database.models import Artist


class ArtistRepository:
    def __init__(self, session: AsyncSession):
        self.session = session

    async def search_by_name(self, query: str, limit: int, offset: int):
        """
        Busca artistas cuyo artist_name contenga 'query' (ilike), paginados.
        ORDEN ESPECÍFICO para que coincida con el serializer actual.
        """
        stmt = (
            select(
                Artist.id,  # índice [0] - para serializer
                Artist.artist_name,  # índice [1] - para serializer como "name"
                Artist.profile_pic,  # índice [2] - para serializer como "profile_pic"
                Artist.bio,  # índice [3] - adicional si lo necesitas
                Artist.user_id,  # índice [4] - adicional si lo necesitas
            )
            .where(Artist.artist_name.ilike(f"%{query}%"))
            .limit(limit)
            .offset(offset)
        )
        result = await self.session.execute(stmt)
        return result.all()
