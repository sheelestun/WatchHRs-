from pathlib import Path
import sys
sys.path.append(str(Path(__file__).parent.parent))

import json
import numpy as np
from sqlalchemy.orm import Session
from models import Employee
from config import MATCH_THRESHOLD

from detector import detect_faces
from db import get_face_embedding


cv_db: dict[str, dict] = {}
def load_cv_db_from_postgres(db: Session):
    """Загружает эмбеддинги сотрудников из БД в память"""
    global cv_db
    employees = db.query(Employee).filter(Employee.embedding.isnot(None)).all()

    for emp in employees:
        try:
            embedding = np.array(json.loads(emp.embedding))
            cv_db[emp.id] = {
                "embedding": embedding,
                "name": emp.name,
                "email": emp.email
            }
        except Exception as e:
            print(f"Ошибка загрузки {emp.email}: {e}")

    print(f"Загружено {len(cv_db)} сотрудников в память")
    return cv_db


def authenticate_face(image: np.ndarray) -> tuple[str | None, float | None]:
    faces = detect_faces(image)
    if not faces:
        return None, None

    x1, y1, x2, y2 = faces[0]
    face_crop = image[y1:y2, x1:x2]

    if face_crop.size == 0:
        return None, None

    current_embedding = get_face_embedding(face_crop)

    best_id = None
    best_distance = float('inf')

    for emp_id, data in cv_db.items():
        distance = np.linalg.norm(current_embedding - data["embedding"])
        if distance < best_distance:
            best_distance = distance
            best_id = emp_id

    if best_distance >= MATCH_THRESHOLD:
        return None, None

    confidence = 1.0 - best_distance / MATCH_THRESHOLD
    return best_id, confidence