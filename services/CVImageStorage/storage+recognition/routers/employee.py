import sys
from pathlib import Path
sys.path.append(str(Path(__file__).parent.parent))

from fastapi import APIRouter, Depends, HTTPException, Query, UploadFile, File
from sqlalchemy.orm import Session
from typing import Optional
import cv2
import numpy as np
import json

from database import get_db
from models import Employee, Manager, Photo
from storage import save_photo_to_db
from cv_service import cv_db
from cv_service import authenticate_face
router = APIRouter(prefix="/employees", tags=["Employees"])


@router.post("")
async def create_employee(
        name: str = Query(...),
        email: str = Query(...),
        manager_id: Optional[str] = Query(None),
        db: Session = Depends(get_db)
):
    existing = db.query(Employee).filter(Employee.email == email).first()
    if existing:
        raise HTTPException(400, "Employee with this email already exists")

    if manager_id and not db.query(Manager).filter(Manager.id == manager_id).first():
        raise HTTPException(404, "Manager not found")

    employee = Employee(name=name, email=email, manager_id=manager_id)
    db.add(employee)
    db.commit()
    db.refresh(employee)
    return {"id": employee.id, "name": employee.name, "email": employee.email}


@router.get("")
async def get_employees(db: Session = Depends(get_db)):
    employees = db.query(Employee).all()
    return [
        {"id": e.id, "name": e.name, "email": e.email, "manager_id": e.manager_id,
         "has_embedding": e.embedding is not None}
        for e in employees
    ]


@router.get("/{employee_id}")
async def get_employee(employee_id: str, db: Session = Depends(get_db)):
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")
    return {
        "id": employee.id, "name": employee.name, "email": employee.email,
        "manager_id": employee.manager_id,
        "photos_count": len(employee.photos),
        "screenshots_count": len(employee.screenshots)
    }


@router.post("/{employee_id}/photo")
async def upload_employee_photo(
        employee_id: str,
        file: UploadFile = File(...),
        db: Session = Depends(get_db)
):
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    content = await file.read()
    nparr = np.frombuffer(content, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    if img is None:
        raise HTTPException(400, "Invalid image")

    from detector import detect_faces
    from db import get_face_embedding

    faces = detect_faces(img)
    if not faces:
        raise HTTPException(400, "No face detected")

    x1, y1, x2, y2 = faces[0]
    face_crop = img[y1:y2, x1:x2]
    if face_crop.size == 0:
        raise HTTPException(400, "Face crop is empty")

    embedding = get_face_embedding(face_crop)
    photo = save_photo_to_db(
        db=db,
        employee_id=employee_id,
        file_content=content,
        filename=file.filename,
        content_type=file.content_type or "image/png"
    )

    employee.embedding = json.dumps(embedding.tolist())
    db.add(employee)
    db.commit()

    # Обновляем память
    cv_db[employee_id] = {"embedding": embedding, "name": employee.name, "email": employee.email}

    return {"employee_id": employee_id, "name": employee.name, "photo_id": photo.id, "status": "uploaded to DB"}


@router.post("/auth", tags=["Auth"])
async def authenticate(file: UploadFile = File(...), db: Session = Depends(get_db)):
    content = await file.read()
    nparr = np.frombuffer(content, np.uint8)
    img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)

    if img is None:
        raise HTTPException(400, "Invalid image")

    emp_id, confidence = authenticate_face(img)
    if not emp_id:
        raise HTTPException(401, "Authentication failed")

    employee = db.query(Employee).filter(Employee.id == emp_id).first()
    return {
        "employee_id": emp_id,
        "name": employee.name,
        "email": employee.email,
        "authenticated": True,
        "confidence": confidence
    }