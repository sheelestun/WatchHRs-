from pathlib import Path
import sys
sys.path.append(str(Path(__file__).parent.parent))

from fastapi import FastAPI
from contextlib import asynccontextmanager

from database import engine, Base, SessionLocal
from cv_service import load_cv_db_from_postgres
from routers import managers, employees, screenshots, work_sessions

Base.metadata.create_all(bind=engine)

@asynccontextmanager
async def lifespan(app: FastAPI):
    db = SessionLocal()
    try:
        load_cv_db_from_postgres(db)
    finally:
        db.close()
    yield

app = FastAPI(
    title="Employee Monitoring System",
    version="4.0 (PostgreSQL BYTEA Storage)",
    lifespan=lifespan
)

app.include_router(managers.router)
app.include_router(employees.router)
app.include_router(screenshots.router)
app.include_router(work_sessions.router)

@app.get("/")
async def root():
    return {
        "service": "Employee Monitoring System",
        "version": "4.0",
        "storage": "PostgreSQL (BYTEA) - все файлы в БД",
        "docs": "/docs",
        "endpoints": {
            "Managers": "GET/POST /managers",
            "Employees": "GET/POST /employees, POST /employees/{id}/photo, POST /auth",
            "Screenshots": "POST /screenshots, GET /screenshots/employee/{id}, GET /screenshot/{id}/download",
            "Work Sessions": "POST /work-sessions/start, POST /work-sessions/{id}/end, GET /work-sessions/employee/{id}"
        }
    }

if __name__ == "__main__":
    import uvicorn
    uvicorn.run("api.main:app", host="127.0.0.1", port=8000, reload=True)