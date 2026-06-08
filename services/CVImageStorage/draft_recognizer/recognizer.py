import numpy as np

from config import MATCH_THRESHOLD
from detector import detect_faces
from db import get_face_embedding


def identify_faces(image, db):
    faces = detect_faces(image)
    if not faces:
        return None

    x1, y1, x2, y2 = faces[0]
    face_crop = image[y1:y2, x1:x2]

    if face_crop.size == 0:
        return None

    current_embedding = get_face_embedding(face_crop)

    best_name = None
    best_distance = float('inf')

    for name, db_embedding in db.items():
        distance = np.linalg.norm(current_embedding - db_embedding)
        if distance < best_distance:
            best_distance = distance
            best_name = name

    return best_name if best_distance < MATCH_THRESHOLD else None