# main.py
from fastapi import FastAPI
from handlers.playlist_handlers import router as playlist_router
from middleware.auth_middleware import AuthMiddleware
import uvicorn

app = FastAPI(
    title="Playlist Service",
    version="0.1",
    description="Microservicio para gestión de playlists de música",
    docs_url="/docs",
    redoc_url="/redoc",
)

# Middleware global
app.add_middleware(AuthMiddleware)

# Router de playlists
app.include_router(playlist_router, prefix="/playlists", tags=["Playlists"])


# Health check
@app.get("/health")
def health_check():
    return {"status": "ok", "service": "playlist-service"}


# Root endpoint
@app.get("/")
def root():
    return {
        "message": "Playlist Service API",
        "version": "0.1",
        "docs": "/docs",
        "health": "/health",
    }


if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8006, reload=True)
