import os
import cv2
import numpy as np
import torch
from facenet_pytorch import InceptionResnetV1
from config import KNOWN_FACES_DIR
from detector import detect_faces

print("Загрузка FaceNet...")
recognizer = InceptionResnetV1(pretrained='vggface2').eval()
print("✓ FaceNet готов")


def get_face_embedding(face_image):
    face_rgb = cv2.cvtColor(face_image, cv2.COLOR_BGR2RGB)
    face_resized = cv2.resize(face_rgb, (160, 160))
    face_float = face_resized.astype(np.float32) / 255.0
    face_normalized = (face_float - 0.5) / 0.5
    face_tensor = torch.from_numpy(face_normalized.transpose(2, 0, 1)).unsqueeze(0)
    with torch.no_grad():
        embedding = recognizer(face_tensor)

    return embedding.squeeze().numpy()


def load_db():
    db = {}
    if not os.path.exists(KNOWN_FACES_DIR):
        print(f"Папка {KNOWN_FACES_DIR} не найдена")
        return db

    for filename in os.listdir(KNOWN_FACES_DIR):
        if filename.endswith((".jpg", ".jpeg", ".png", ".JPG", ".PNG")):
            path = os.path.join(KNOWN_FACES_DIR, filename)
            name = os.path.splitext(filename)[0]

            image = cv2.imread(path)
            if image is None:
                print(f"Не прочитано: {filename}")
                continue

            faces = detect_faces(image)
            if not faces:
                print(f"Нет лица: {filename}")
                continue

            x1, y1, x2, y2 = faces[0]
            face_crop = image[y1:y2, x1:x2]
            if face_crop.shape[0] < 60 or face_crop.shape[1] < 60:
                print(f"Лицо слишком маленькое: {filename}")
                continue

            embedding = get_face_embedding(face_crop)
            db[name] = embedding
            print(f"Добавлен: {name}")

    print(f"\nБаза готова: {len(db)} человек(а)\n")
    return db