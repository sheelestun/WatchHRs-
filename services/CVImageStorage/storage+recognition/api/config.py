import os
from pathlib import Path

STORAGE_ROOT = Path(os.getenv("STORAGE_ROOT", "storage"))
PHOTOS_DIR = STORAGE_ROOT / "photos"
SCREENSHOTS_DIR = STORAGE_ROOT / "screenshots"
EMBEDDINGS_DIR = STORAGE_ROOT / "embeddings"

MATCH_THRESHOLD = float(os.getenv("MATCH_THRESHOLD", "0.9"))
MAX_FILE_SIZE = int(os.getenv("MAX_FILE_SIZE", str(15 * 1024 * 1024)))
HOST = os.getenv("HOST", "0.0.0.0")
PORT = int(os.getenv("PORT", "8081"))
