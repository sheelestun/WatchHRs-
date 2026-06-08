import cv2
import time
import sys
import json
import numpy as np
from sqlalchemy import create_engine, text
from sqlalchemy.orm import sessionmaker
from config import TIMEOUT_SECONDS, CAMERA_INDEX
from detector import detect_faces
from db import get_face_embedding

DATABASE_URL = "postgresql://postgres:MrLogan555@localhost:5432/cv_storage"
engine = create_engine(DATABASE_URL)
SessionLocal = sessionmaker(bind=engine)


def load_db_from_postgres():
    """Загрузка эмбеддингов из таблицы employees"""
    db = {}
    session = SessionLocal()

    try:
        result = session.execute(
            text("SELECT id, name, embedding FROM employees WHERE embedding IS NOT NULL")
        )
        rows = result.fetchall()

        for emp_id, name, embedding_json in rows:
            try:
                embedding = np.array(json.loads(embedding_json))
                db[emp_id] = {"embedding": embedding, "name": name}
                print(f"✅ Загружен: {name}")
            except Exception as e:
                print(f"❌ Ошибка {name}: {e}")

        print(f"\n📦 База из PostgreSQL: {len(db)} человек(а)\n")

    except Exception as e:
        print(f"❌ Ошибка подключения: {e}")
        print("💡 Убедитесь, что PostgreSQL запущен и таблица employees существует")
    finally:
        session.close()

    return db


def main():
    known_db = load_db_from_postgres()

    if not known_db:
        print("❌ База пуста!")
        print("💡 Запустите: python import_faces.py")
        sys.exit(1)

    cap = cv2.VideoCapture(CAMERA_INDEX)
    if not cap.isOpened():
        print("❌ Не удалось открыть камеру")
        sys.exit(1)

    start_time = time.time()
    result_name = None
    consecutive = 0
    last_name = None

    try:
        while True:
            elapsed = time.time() - start_time
            remain = int(TIMEOUT_SECONDS - elapsed)

            if remain <= 0:
                break

            ret, frame = cap.read()
            if not ret:
                break
            faces = detect_faces(frame)
            name = None

            if faces:
                x1, y1, x2, y2 = faces[0]
                face_crop = frame[y1:y2, x1:x2]

                if face_crop.size > 0:
                    current_embedding = get_face_embedding(face_crop)

                    best_id = None
                    best_distance = float('inf')

                    for emp_id, data in known_db.items():
                        distance = np.linalg.norm(current_embedding - data["embedding"])
                        if distance < best_distance:
                            best_distance = distance
                            best_id = emp_id

                    if best_distance < 0.9:
                        name = known_db[best_id]["name"]

            if name == last_name and name is not None:
                consecutive += 1
                if consecutive >= 3:
                    result_name = name
                    break
            else:
                consecutive = 0
                last_name = name

            cv2.putText(frame, f"Time: {remain}s", (10, 30),
                        cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 255, 0), 2)
            cv2.putText(frame, f"Conf: {consecutive}/3", (10, 70),
                        cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 255, 255), 2)

            if name:
                cv2.putText(frame, f"Employee: {name}", (10, 110),
                            cv2.FONT_HERSHEY_SIMPLEX, 0.8, (0, 255, 0), 2)

            cv2.imshow('FaceID (PostgreSQL)', frame)

            if cv2.waitKey(1) & 0xFF == ord('q'):
                break

    finally:
        cap.release()
        cv2.destroyAllWindows()

    print("-" * 50)
    if result_name:
        print(f"✅ УСПЕХ: {result_name}")
    else:
        print("❌ НЕ РАСПОЗНАН!")
        sys.exit(1)


if __name__ == "__main__":
    main()