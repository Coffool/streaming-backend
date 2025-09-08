from sqlalchemy.ext.asyncio import AsyncSession
from strategies.base_strategy import SearchStrategy
from services.serializers import serialize_song, serialize_album, serialize_artist


class SearchService:
    def __init__(self, strategy: SearchStrategy):
        self.strategy = strategy

    async def search(
        self,
        session: AsyncSession,
        query: str,
        limit: int = 5,
        offset_songs: int = 0,
        offset_albums: int = 0,
        offset_artists: int = 0,
    ) -> dict:
        songs, albums, artists = await self.strategy.search(
            session, query, limit, offset_songs, offset_albums, offset_artists
        )

        return {
            "songs": {
                "page": (offset_songs // limit) + 1,
                "results": [serialize_song(s) for s in songs],
            },
            "albums": {
                "page": (offset_albums // limit) + 1,
                "results": [serialize_album(a) for a in albums],
            },
            "artists": {
                "page": (offset_artists // limit) + 1,
                "results": [serialize_artist(ar) for ar in artists],
            },
        }
