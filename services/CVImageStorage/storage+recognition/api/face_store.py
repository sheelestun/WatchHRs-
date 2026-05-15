import cv2
import numpy as np

from config import EMBEDDINGS_DIR, MATCH_THRESHOLD, PHOTOS_DIR
from detector import detect_faces
from face_embedding import get_face_embedding

_embeddings: dict[str, np.ndarray] = {}


def ensure_dirs() -> None:
    PHOTOS_DIR.mkdir(parents=True, exist_ok=True)
    EMBEDDINGS_DIR.mkdir(parents=True, exist_ok=True)


def _embedding_path(user_id: str):
    return EMBEDDINGS_DIR / f"{user_id}.npy"


def _load_embedding_from_photo(user_id: str) -> np.ndarray | None:
    photo = PHOTOS_DIR / f"{user_id}.png"
    if not photo.exists():
        return None

    image = cv2.imread(str(photo))
    if image is None:
        return None

    faces = detect_faces(image)
    if not faces:
        return None

    x1, y1, x2, y2 = faces[0]
    face_crop = image[y1:y2, x1:x2]
    if face_crop.size == 0:
        return None

    return get_face_embedding(face_crop)


def load_all_embeddings() -> None:
    global _embeddings
    _embeddings = {}

    for path in PHOTOS_DIR.glob("*.png"):
        user_id = path.stem
        embedding_file = _embedding_path(user_id)
        if embedding_file.exists():
            embedding = np.load(embedding_file)
        else:
            embedding = _load_embedding_from_photo(user_id)
            if embedding is None:
                continue
            np.save(embedding_file, embedding)

        _embeddings[user_id] = embedding

    print(f"Loaded {len(_embeddings)} face embedding(s)")


def register_photo(user_id: str, content: bytes) -> None:
    nparr = np.frombuffer(content, np.uint8)
    image = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    if image is None:
        raise ValueError("invalid image")

    faces = detect_faces(image)
    if not faces:
        raise ValueError("no face detected")

    x1, y1, x2, y2 = faces[0]
    face_crop = image[y1:y2, x1:x2]
    if face_crop.size == 0:
        raise ValueError("face crop is empty")

    embedding = get_face_embedding(face_crop)
    np.save(_embedding_path(user_id), embedding)
    _embeddings[user_id] = embedding


def remove_user(user_id: str) -> None:
    _embeddings.pop(user_id, None)
    embedding_file = _embedding_path(user_id)
    if embedding_file.exists():
        embedding_file.unlink()


def authenticate(image: np.ndarray) -> tuple[str | None, float | None]:
    faces = detect_faces(image)
    if not faces:
        return None, None

    x1, y1, x2, y2 = faces[0]
    face_crop = image[y1:y2, x1:x2]
    if face_crop.size == 0:
        return None, None

    current = get_face_embedding(face_crop)
    best_id = None
    best_distance = float("inf")

    for user_id, embedding in _embeddings.items():
        distance = float(np.linalg.norm(current - embedding))
        if distance < best_distance:
            best_distance = distance
            best_id = user_id

    if best_id is None or best_distance >= MATCH_THRESHOLD:
        return None, None

    confidence = 1.0 - best_distance / MATCH_THRESHOLD
    return best_id, confidence
