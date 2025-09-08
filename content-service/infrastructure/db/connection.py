from typing import AsyncGenerator
from sqlalchemy.ext.asyncio import create_async_engine, async_sessionmaker, AsyncSession
from sqlalchemy.orm import DeclarativeBase
from config import settings


engine = create_async_engine(
    settings.db_url,
    echo=True,
    connect_args={
        "statement_cache_size": 0,  # 🔹 Desactiva cache de asyncpg
        "prepared_statement_cache_size": 0,  # 🔹 Más seguro con PgBouncer
    },
    pool_pre_ping=True,  # 🔹 Revisa conexiones muertas (recomendado en PgBouncer)
)


# Factory de sesiones asincrónicas
AsyncSessionLocal = async_sessionmaker(
    bind=engine,
    expire_on_commit=False,
)


class Base(DeclarativeBase):
    pass


# Dependency para FastAPI
async def get_db() -> AsyncGenerator[AsyncSession, None]:
    async with AsyncSessionLocal() as session:
        yield session
