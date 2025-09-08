from fastapi import FastAPI
from handlers.search_handler import router as search_router
from middleware.auth_middleware import AuthMiddleware
import uvicorn

app = FastAPI(title="Search Service", version="0.1")

# Middleware global
app.add_middleware(AuthMiddleware)

# Router de búsqueda
app.include_router(search_router, prefix="/search", tags=["Search"])


# Health check
@app.get("/health")
def health_check():
    return {"status": "ok"}


if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8004, reload=True)
