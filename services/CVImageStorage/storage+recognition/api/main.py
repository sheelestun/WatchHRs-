from contextlib import asynccontextmanager

from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from config import MAX_FILE_SIZE, HOST, PORT
from face_store import ensure_dirs as ensure_face_dirs, load_all_embeddings
from file_storage import ensure_dirs as ensure_storage_dirs
from routers import auth, photo, screenshot


@asynccontextmanager
async def lifespan(_: FastAPI):
    ensure_storage_dirs()
    ensure_face_dirs()
    load_all_embeddings()
    yield


app = FastAPI(title="CV Image Storage", version="2.0.0", lifespan=lifespan)

app.include_router(photo.router)
app.include_router(screenshot.router)
app.include_router(auth.router)


@app.middleware("http")
async def limit_body_size(request: Request, call_next):
    content_length = request.headers.get("content-length")
    if content_length and int(content_length) > MAX_FILE_SIZE:
        return JSONResponse({"detail": "payload too large"}, status_code=413)
    return await call_next(request)


@app.get("/health")
async def health():
    return {"status": "ok"}


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("main:app", host=HOST, port=PORT, reload=False)
