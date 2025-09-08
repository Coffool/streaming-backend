from fastapi import APIRouter, Depends, Query
from sqlalchemy.ext.asyncio import AsyncSession
from services.search_service import SearchService
from strategies.fuzzy_strategy import FuzzySearchStrategy
from database.connection import get_db  # Tu función que devuelve AsyncSession

router = APIRouter()


@router.get("/search")
async def search(
    q: str = Query(..., description="Texto a buscar"),
    song_page: int = Query(1, ge=1, description="Página de canciones"),
    album_page: int = Query(1, ge=1, description="Página de álbumes"),
    artist_page: int = Query(1, ge=1, description="Página de artistas"),
    limit: int = Query(5, ge=1, le=50, description="Número de resultados por página"),
    db: AsyncSession = Depends(get_db),
):
    """
    Endpoint de búsqueda que devuelve canciones, álbumes y artistas
    con paginado independiente.
    """
    strategy = FuzzySearchStrategy(threshold=70)
    service = SearchService(strategy)

    # Calculamos offsets según la página
    offset_songs = (song_page - 1) * limit
    offset_albums = (album_page - 1) * limit
    offset_artists = (artist_page - 1) * limit

    # Llamada async al servicio
    result = await service.search(
        session=db,
        query=q,
        limit=limit,
        offset_songs=offset_songs,
        offset_albums=offset_albums,
        offset_artists=offset_artists,
    )

    return result
