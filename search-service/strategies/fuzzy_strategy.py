from rapidfuzz import fuzz
from typing import List, Tuple
from strategies.base_strategy import SearchStrategy
from repositories.song_repository import SongRepository
from repositories.album_repository import AlbumRepository
from repositories.artist_repository import ArtistRepository
from sqlalchemy.ext.asyncio import AsyncSession


class FuzzySearchStrategy(SearchStrategy):
    def __init__(self, threshold: int = 70):
        self.threshold = threshold

    async def _filter_tuples(self, tuples_list, query: str, field_index: int) -> List:
        """
        Filtra una lista de TUPLAS usando fuzzy matching sobre un campo específico.

        Args:
            tuples_list: Lista de tuplas retornadas por los repositorios
            query: Término de búsqueda
            field_index: Índice en la tupla donde está el campo a comparar
        """
        filtered = []
        for tuple_item in tuples_list:
            # Verificar que la tupla tenga suficientes elementos y el campo no sea None
            if len(tuple_item) > field_index and tuple_item[field_index]:
                similarity = fuzz.partial_ratio(
                    query.lower(), str(tuple_item[field_index]).lower()
                )
                if similarity >= self.threshold:
                    filtered.append((tuple_item, similarity))

        # Ordenar por similitud (mayor a menor) y retornar solo las tuplas
        filtered.sort(key=lambda x: x[1], reverse=True)
        return [item for item, score in filtered]

    async def search(
        self,
        session: AsyncSession,
        query: str,
        limit: int,
        offset_songs: int,
        offset_albums: int,
        offset_artists: int,
    ) -> Tuple[List, List, List]:
        song_repo = SongRepository(session)
        album_repo = AlbumRepository(session)
        artist_repo = ArtistRepository(session)

        try:
            # Traer resultados de los repositorios (ahora TODOS retornan tuplas)
            songs = await song_repo.get_by_title_ilike(query, limit * 3, offset_songs)
            albums = await album_repo.get_by_title_ilike(
                query, limit * 3, offset_albums
            )
            artists = await artist_repo.search_by_name(query, limit * 3, offset_artists)

            # Filtrar usando fuzzy matching sobre las tuplas
            # Índices: songs[1]=title, albums[1]=title, artists[1]=artist_name
            filtered_songs = await self._filter_tuples(
                songs, query, 1
            )  # title está en índice 1
            filtered_albums = await self._filter_tuples(
                albums, query, 1
            )  # title está en índice 1
            filtered_artists = await self._filter_tuples(
                artists, query, 1
            )  # artist_name está en índice 1

            return (
                filtered_songs[:limit],
                filtered_albums[:limit],
                filtered_artists[:limit],
            )

        except Exception as e:
            print(f"Error en búsqueda fuzzy: {e}")
            return [], [], []
