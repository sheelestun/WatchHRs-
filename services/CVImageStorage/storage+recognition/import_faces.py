import os
import cv2
import uuid
import json
from datetime import datetime
from pathlib import Path
from sqlalchemy import create_engine, Column, String, DateTime, Text, ForeignKey
from sqlalchemy.orm import declarative_base, sessionmaker

from config import KNOWN_FACES_DIR
from detector import detect_faces
from db import get_face_embedding

DATABASE_URL = "postgresql://postgres:MrLogan555@localhost:5432/cv_storage"
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(bind=engine)
Base = declarative_base()


class Employee(Base):
    __tablename__ = "employees"
    id = Column(String, primary_key=True)
    name = Column(String, nullable=False)
    email = Column(String, unique=True, nullable=False)
    embedding = Column(Text, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)


class Photo(Base):
    __tablename__ = "photos"
    id = Column(String, primary_key=True)
    user_id = Column(String, ForeignKey("employees.id"))
    filename = Column(String)
    path = Column(String)
    uploaded_at = Column(DateTime, default=datetime.utcnow)


Base.metadata.create_all(bind=engine)


def main():
    print("=" * 60)
    print("📥 Импорт фото из face_db/ в PostgreSQL")
    print("=" * 60)

    if not os.path.exists(KNOWN_FACES_DIR):
        print(f"❌ Папка {KNOWN_FACES_DIR} не найдена!")
        return

    db = SessionLocal()
    imported = 0
    skipped = 0

    try:
        for filename in os.listdir(KNOWN_FACES_DIR):
            if not filename.lower().endswith((".jpg", ".jpeg", ".png")):
                continue

            employee_id_name = os.path.splitext(filename)[0]
            email = f"{employee_id_name}@company.com"
            existing = db.query(Employee).filter(Employee.email == email).first()
            if existing:
                print(f"⏭️  Уже существует: {employee_id_name}")
                skipped += 1
                continue

            path = os.path.join(KNOWN_FACES_DIR, filename)
            image = cv2.imread(path)
            if image is None:
                print(f"❌ Не прочитано: {filename}")
                continue

            faces = detect_faces(image)
            if not faces:
                print(f"  ❌ Нет лица: {filename}")
                continue

            x1, y1, x2, y2 = faces[0]
            face_crop = image[y1:y2, x1:x2]

            if face_crop.shape[0] < 60 or face_crop.shape[1] < 60:
                print(f"  ❌ Лицо слишком маленькое: {filename}")
                continue

            print(f"  🔄 Обработка...")
            embedding = get_face_embedding(face_crop)

            emp_id = str(uuid.uuid4())
            employee = Employee(
                id=emp_id,
                name=employee_id_name.replace("_", " ").title(),
                email=email,
                embedding=json.dumps(embedding.tolist())
            )

            photos_dir = Path("storage/photos")
            photos_dir.mkdir(parents=True, exist_ok=True)
            photo_filename = f"{emp_id}.png"
            photo_path = photos_dir / photo_filename
            cv2.imwrite(str(photo_path), image)

            photo = Photo(
                id=str(uuid.uuid4()),
                user_id=emp_id,
                filename=filename,
                path=str(photo_path)
            )

            db.add(employee)
            db.flush()

            db.add(photo)
            imported += 1
            print(f"✅ Добавлен: {employee_id_name}")

        db.commit()
        print("=" * 60)
        print(f"📊 Итого: ✅ {imported} импортировано, ⏭️ {skipped} пропущено")
        print("=" * 60)

    except Exception as e:
        db.rollback()
        print(f"❌ Ошибка: {e}")
        raise
    finally:
        db.close()


if __name__ == "__main__":
    main()