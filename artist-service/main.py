from fastapi import FastAPI
from handlers.artist_handler import router as artist_router
from middleware.auth_middleware import AuthMiddleware
import uvicorn

app = FastAPI(title="Artist Service", version="0.1")

# Registrar middleware de autenticación
app.add_middleware(AuthMiddleware)

# Rutas
app.include_router(artist_router, prefix="/artists", tags=["artists"])


@app.get("/health")
def health_check():
    return {"status": "ok"}


# Permitir ejecución directa con python3 main.py
if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8001, reload=True)
