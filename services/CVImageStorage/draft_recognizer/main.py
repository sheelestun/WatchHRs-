import cv2
import time
import sys

from config import TIMEOUT_SECONDS, CAMERA_INDEX
from db import load_db
from recognizer import identify_faces


def main():
    known_db = load_db()
    if not known_db:
        print("База пуста!")
        sys.exit(1)

    cap = cv2.VideoCapture(CAMERA_INDEX)
    if not cap.isOpened():
        print("Не удалось открыть камеру")
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

            name = identify_faces(frame, known_db)

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

            cv2.imshow('FaceID', frame)

            if cv2.waitKey(1) & 0xFF == ord('q'):
                break

    finally:
        cap.release()
        cv2.destroyAllWindows()

    print("-" * 50)
    if result_name:
        print(f"УСПЕХ: {result_name}")
    else:
        print("НЕ РАСПОЗНАН! Попробуйте снова.")
        sys.exit(1)


if __name__ == "__main__":
    main()