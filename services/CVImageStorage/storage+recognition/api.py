import cv2
import numpy as np
import uuid
import json
from datetime import datetime
from pathlib import Path
from typing import Optional
from fastapi import FastAPI, UploadFile, File, HTTPException, Depends, Query
from sqlalchemy import create_engine, Column, String, DateTime, ForeignKey, Index, Text, Integer
from sqlalchemy.orm import declarative_base, relationship, sessionmaker, Session


from config import MATCH_THRESHOLD
from detector import detect_faces
from db import get_face_embedding

BASE_DIR = Path(__file__).parent
STORAGE_PATH = BASE_DIR / "storage"
PHOTOS_DIR = STORAGE_PATH / "photos"
SCREENSHOTS_DIR = STORAGE_PATH / "screenshots"
PHOTOS_DIR.mkdir(parents=True, exist_ok=True)
SCREENSHOTS_DIR.mkdir(parents=True, exist_ok=True)

DATABASE_URL = "postgresql://postgres:MrLogan555@localhost:5432/cv_storage"
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(bind=engine)
Base = declarative_base()

class Manager(Base):
    __tablename__ = "managers"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    name = Column(String, nullable=False)
    email = Column(String, unique=True, nullable=False)
    created_at = Column(DateTime, default=datetime.utcnow)

    employees = relationship("Employee", back_populates="manager", cascade="all, delete-orphan")


class Employee(Base):
    __tablename__ = "employees"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    name = Column(String, nullable=False)
    email = Column(String, unique=True, nullable=False)
    manager_id = Column(String, ForeignKey("managers.id"), nullable=True)
    embedding = Column(Text, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    manager = relationship("Manager", back_populates="employees")
    photos = relationship("Photo", back_populates="employee", cascade="all, delete-orphan")
    screenshots = relationship("Screenshot", back_populates="employee", cascade="all, delete-orphan")
    work_sessions = relationship("WorkSession", back_populates="employee", cascade="all, delete-orphan")


class Photo(Base):
    __tablename__ = "photos"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    user_id = Column(String, ForeignKey("employees.id"), index=True)
    filename = Column(String)
    path = Column(String)
    uploaded_at = Column(DateTime, default=datetime.utcnow)

    employee = relationship("Employee", back_populates="photos")


class Screenshot(Base):
    __tablename__ = "screenshots"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    employee_id = Column(String, ForeignKey("employees.id"), index=True)
    filename = Column(String)
    path = Column(String)
    cnt_mouse_clicks = Column(Integer, default=0)
    cnt_keyboard_clicks = Column(Integer, default=0)
    timestamp = Column(DateTime, default=datetime.utcnow, index=True)

    employee = relationship("Employee", back_populates="screenshots")

    __table_args__ = (
        Index('idx_employee_timestamp', 'employee_id', 'timestamp'),
    )


class WorkSession(Base):
    __tablename__ = "work_sessions"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    employee_id = Column(String, ForeignKey("employees.id"), index=True)
    start_time = Column(DateTime, nullable=False)
    end_time = Column(DateTime, nullable=True)
    total_time = Column(Integer, nullable=True)  # в секундах
    created_at = Column(DateTime, default=datetime.utcnow)

    employee = relationship("Employee", back_populates="work_sessions")

Base.metadata.create_all(bind=engine)
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()

cv_db = {}
def load_cv_db_from_postgres():
    """Загружает эмбеддинги сотрудников из PostgreSQL в память"""
    global cv_db
    db = SessionLocal()

    try:
        employees = db.query(Employee).filter(
            Employee.embedding.isnot(None)
        ).all()

        for employee in employees:
            try:
                embedding = np.array(json.loads(employee.embedding))
                cv_db[employee.id] = {
                    "embedding": embedding,
                    "name": employee.name,
                    "email": employee.email
                }
            except Exception as e:
                print(f"❌ Ошибка загрузки {employee.email}: {e}")

        print(f"✅ Загружено {len(cv_db)} сотрудников из PostgreSQL")
    finally:
        db.close()

load_cv_db_from_postgres()
def save_photo(file_content: bytes, employee_id: str) -> str:
    filename = f"{employee_id}.png"
    path = PHOTOS_DIR / filename
    with open(path, "wb") as f:
        f.write(file_content)
    return str(path)


def save_screenshot(file_content: bytes, employee_id: str) -> str:
    timestamp = datetime.utcnow().strftime("%Y%m%d_%H%M%S")
    unique_id = str(uuid.uuid4())[:8]
    filename = f"{employee_id}-screenshot-{timestamp}-{unique_id}.png"

    employee_dir = SCREENSHOTS_DIR / employee_id
    employee_dir.mkdir(exist_ok=True)

    path = employee_dir / filename
    with open(path, "wb") as f:
        f.write(file_content)
    return str(path)

app = FastAPI(title="Employee Monitoring System", version="3.0")
@app.post("/managers")
async def create_manager(
        name: str = Query(...),
        email: str = Query(...),
        db: Session = Depends(get_db)
):
    """Создать менеджера"""
    existing = db.query(Manager).filter(Manager.email == email).first()
    if existing:
        raise HTTPException(400, "Manager with this email already exists")

    manager = Manager(
        id=str(uuid.uuid4()),
        name=name,
        email=email
    )
    db.add(manager)
    db.commit()
    db.refresh(manager)

    return {"id": manager.id, "name": manager.name, "email": manager.email}


@app.get("/managers")
async def get_managers(db: Session = Depends(get_db)):
    """Получить всех менеджеров"""
    managers = db.query(Manager).all()
    return [
        {
            "id": m.id,
            "name": m.name,
            "email": m.email,
            "employees_count": len(m.employees)
        }
        for m in managers
    ]

@app.post("/employees")
async def create_employee(
        name: str = Query(...),
        email: str = Query(...),
        manager_id: Optional[str] = Query(None),
        db: Session = Depends(get_db)
):
    """Создать сотрудника"""
    existing = db.query(Employee).filter(Employee.email == email).first()
    if existing:
        raise HTTPException(400, "Employee with this email already exists")

    if manager_id:
        manager = db.query(Manager).filter(Manager.id == manager_id).first()
        if not manager:
            raise HTTPException(404, "Manager not found")

    employee = Employee(
        id=str(uuid.uuid4()),
        name=name,
        email=email,
        manager_id=manager_id
    )
    db.add(employee)
    db.commit()
    db.refresh(employee)

    return {"id": employee.id, "name": employee.name, "email": employee.email}


@app.get("/employees")
async def get_employees(db: Session = Depends(get_db)):
    """Получить всех сотрудников"""
    employees = db.query(Employee).all()
    return [
        {
            "id": e.id,
            "name": e.name,
            "email": e.email,
            "manager_id": e.manager_id,
            "has_embedding": e.embedding is not None
        }
        for e in employees
    ]


@app.get("/employees/{employee_id}")
async def get_employee(employee_id: str, db: Session = Depends(get_db)):
    """Получить информацию о сотруднике"""
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    return {
        "id": employee.id,
        "name": employee.name,
        "email": employee.email,
        "manager_id": employee.manager_id,
        "photos_count": len(employee.photos),
        "screenshots_count": len(employee.screenshots),
        "work_sessions_count": len(employee.work_sessions)
    }


@app.post("/employees/{employee_id}/photo")
async def upload_employee_photo(
        employee_id: str,
        file: UploadFile = File(...),
        db: Session = Depends(get_db)
):
    """Загрузить фото сотрудника для распознавания"""
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    content = await file.read()
    nparr = np.frombuffer(content, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    if img is None:
        raise HTTPException(400, "Invalid image file")

    faces = detect_faces(img)
    if not faces:
        raise HTTPException(400, "No face detected")

    x1, y1, x2, y2 = faces[0]
    face_crop = img[y1:y2, x1:x2]

    if face_crop.size == 0:
        raise HTTPException(400, "Face crop is empty")

    embedding = get_face_embedding(face_crop)
    path = save_photo(content, employee_id)

    photo = Photo(
        user_id=employee_id,
        filename=file.filename,
        path=path
    )

    employee.embedding = json.dumps(embedding.tolist())
    db.add(employee)
    db.add(photo)
    db.commit()

    # Обновляем CV базу
    cv_db[employee_id] = {
        "embedding": embedding,
        "name": employee.name,
        "email": employee.email
    }

    return {
        "employee_id": employee_id,
        "name": employee.name,
        "status": "photo uploaded and embedding saved"
    }


@app.post("/auth")
async def authenticate_employee(
        file: UploadFile = File(...),
        db: Session = Depends(get_db)
):
    """Аутентификация сотрудника по лицу"""
    content = await file.read()
    nparr = np.frombuffer(content, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    if img is None:
        raise HTTPException(400, "Invalid image file")

    faces = detect_faces(img)
    if not faces:
        raise HTTPException(400, "No face detected")

    x1, y1, x2, y2 = faces[0]
    face_crop = img[y1:y2, x1:x2]

    if face_crop.size == 0:
        raise HTTPException(400, "Face crop is empty")

    current_embedding = get_face_embedding(face_crop)

    best_employee_id = None
    best_distance = float('inf')

    for emp_id, data in cv_db.items():
        distance = np.linalg.norm(current_embedding - data["embedding"])
        if distance < best_distance:
            best_distance = distance
            best_employee_id = emp_id

    if best_distance >= MATCH_THRESHOLD:
        raise HTTPException(401, "Authentication failed: face not recognized")

    employee = db.query(Employee).filter(Employee.id == best_employee_id).first()

    return {
        "employee_id": best_employee_id,
        "name": employee.name,
        "email": employee.email,
        "authenticated": True,
        "confidence": float(1.0 - best_distance / MATCH_THRESHOLD)
    }

@app.post("/screenshots")
async def upload_screenshot(
        employee_id: str = Query(...),
        file: UploadFile = File(...),
        cnt_mouse_clicks: int = Query(0),
        cnt_keyboard_clicks: int = Query(0),
        db: Session = Depends(get_db)
):
    """Загрузить скриншот сотрудника"""
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    if not file.filename.lower().endswith(('.png', '.jpg', '.jpeg')):
        raise HTTPException(400, "Invalid format")

    content = await file.read()
    if len(content) > 10 * 1024 * 1024:
        raise HTTPException(400, "File too large")

    path = save_screenshot(content, employee_id)

    screenshot = Screenshot(
        employee_id=employee_id,
        filename=file.filename,
        path=path,
        cnt_mouse_clicks=cnt_mouse_clicks,
        cnt_keyboard_clicks=cnt_keyboard_clicks
    )
    db.add(screenshot)
    db.commit()
    db.refresh(screenshot)

    return {
        "id": screenshot.id,
        "employee_id": employee_id,
        "filename": file.filename,
        "cnt_mouse_clicks": cnt_mouse_clicks,
        "cnt_keyboard_clicks": cnt_keyboard_clicks,
        "timestamp": screenshot.timestamp.isoformat()
    }


@app.get("/employees/{employee_id}/screenshots")
async def get_employee_screenshots(
        employee_id: str,
        date: Optional[str] = Query(None),
        limit: int = Query(100),
        db: Session = Depends(get_db)
):
    """Получить скриншоты сотрудника"""
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    query = db.query(Screenshot).filter(Screenshot.employee_id == employee_id)

    if date:
        try:
            target_date = datetime.strptime(date, "%Y-%m-%d").date()
            query = query.filter(func.date(Screenshot.timestamp) == target_date)
        except ValueError:
            raise HTTPException(400, "Invalid date format. Use YYYY-MM-DD")

    screenshots = query.order_by(Screenshot.timestamp.desc()).limit(limit).all()

    return {
        "employee_id": employee_id,
        "count": len(screenshots),
        "screenshots": [
            {
                "id": s.id,
                "filename": s.filename,
                "path": s.path,
                "cnt_mouse_clicks": s.cnt_mouse_clicks,
                "cnt_keyboard_clicks": s.cnt_keyboard_clicks,
                "timestamp": s.timestamp.isoformat()
            }
            for s in screenshots
        ]
    }


@app.post("/work-sessions/start")
async def start_work_session(
        employee_id: str = Query(...),
        db: Session = Depends(get_db)
):
    """Начать рабочую сессию"""
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    session = WorkSession(
        employee_id=employee_id,
        start_time=datetime.utcnow()
    )
    db.add(session)
    db.commit()
    db.refresh(session)

    return {
        "session_id": session.id,
        "employee_id": employee_id,
        "start_time": session.start_time.isoformat(),
        "status": "started"
    }


@app.post("/work-sessions/{session_id}/end")
async def end_work_session(
        session_id: str,
        db: Session = Depends(get_db)
):
    """Завершить рабочую сессию"""
    session = db.query(WorkSession).filter(WorkSession.id == session_id).first()
    if not session:
        raise HTTPException(404, "Session not found")

    end_time = datetime.utcnow()
    session.end_time = end_time
    session.total_time = int((end_time - session.start_time).total_seconds())

    db.commit()

    return {
        "session_id": session_id,
        "end_time": end_time.isoformat(),
        "total_time_seconds": session.total_time,
        "status": "completed"
    }


@app.get("/employees/{employee_id}/work-sessions")
async def get_employee_work_sessions(
        employee_id: str,
        date: Optional[str] = Query(None),
        db: Session = Depends(get_db)
):
    """Получить рабочие сессии сотрудника"""
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    query = db.query(WorkSession).filter(WorkSession.employee_id == employee_id)

    if date:
        try:
            target_date = datetime.strptime(date, "%Y-%m-%d").date()
            query = query.filter(func.date(WorkSession.start_time) == target_date)
        except ValueError:
            raise HTTPException(400, "Invalid date format")

    sessions = query.order_by(WorkSession.start_time.desc()).all()

    return {
        "employee_id": employee_id,
        "count": len(sessions),
        "sessions": [
            {
                "id": s.id,
                "start_time": s.start_time.isoformat(),
                "end_time": s.end_time.isoformat() if s.end_time else None,
                "total_time_seconds": s.total_time
            }
            for s in sessions
        ]
    }

@app.get("/")
async def root():
    return {
        "service": "Employee Monitoring System",
        "version": "3.0",
        "database_structure": {
            "managers": "Managers table",
            "employees": "Employees with manager reference",
            "photos": "Employee photos with face embeddings",
            "screenshots": "Screenshots with mouse/keyboard activity",
            "work_sessions": "Work time tracking"
        },
        "endpoints": {
            "Managers": [
                "POST /managers",
                "GET /managers"
            ],
            "Employees": [
                "POST /employees",
                "GET /employees",
                "GET /employees/{id}",
                "POST /employees/{id}/photo",
                "POST /auth"
            ],
            "Screenshots": [
                "POST /screenshots",
                "GET /employees/{id}/screenshots"
            ],
            "Work Sessions": [
                "POST /work-sessions/start",
                "POST /work-sessions/{id}/end",
                "GET /employees/{id}/work-sessions"
            ]
        }
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="127.0.0.1", port=8000)