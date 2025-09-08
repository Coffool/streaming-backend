from fastapi import FastAPI
from core.handlers.album_handler import router as album_router
from core.handlers.song_handler import router as song_router
from middleware.auth_middleware import AuthMiddleware
import asyncio
from contextlib import asynccontextmanager

from events.consumer import consume_events  # consumer del artist_created
import uvicorn

# -------------------------
# Lifespan handler para startup y shutdown
# -------------------------
@asynccontextmanager
async def lifespan(_):
    # Startup: iniciar consumer en background
    task = asyncio.create_task(consume_events())
    print("[*] Consumer de artist_created iniciado en background.")
    yield
    # Shutdown: cancelar tarea si sigue corriendo
    task.cancel()
    try:
        await task
    except asyncio.CancelledError:
        print("[*] Consumer detenido correctamente.")


# -------------------------
# Crear app con lifespan
# -------------------------
app = FastAPI(title="Music Service", version="0.1", lifespan=lifespan)

# Middleware
app.add_middleware(AuthMiddleware)

# Rutas (sin duplicar prefix/tags)
app.include_router(album_router)
app.include_router(song_router)


# Health check
@app.get("/health")
def health_check():
    return {"status": "ok"}

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8002, reload=True)