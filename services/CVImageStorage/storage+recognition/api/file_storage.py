import io
import re
import uuid
import zipfile
from datetime import date, datetime, timezone
from pathlib import Path

from fastapi import HTTPException

from config import PHOTOS_DIR, SCREENSHOTS_DIR

PHOTO_NAME_RE = re.compile(r"^[0-9a-fA-F-]{36}\.png$")
SCREENSHOT_NAME_RE = re.compile(
    r"^[0-9a-fA-F-]{36}-[0-9a-fA-F-]{36}\.png$"
)


def ensure_dirs() -> None:
    PHOTOS_DIR.mkdir(parents=True, exist_ok=True)
    SCREENSHOTS_DIR.mkdir(parents=True, exist_ok=True)


def photo_path(user_id: str) -> Path:
    return PHOTOS_DIR / f"{user_id}.png"


def screenshot_dir(employee_id: str) -> Path:
    return SCREENSHOTS_DIR / employee_id


def screenshot_path(employee_id: str, screenshot_id: str) -> Path:
    return screenshot_dir(employee_id) / f"{employee_id}-{screenshot_id}.png"


def save_photo(user_id: str, content: bytes) -> str:
    if not re.fullmatch(r"[0-9a-fA-F-]{36}", user_id):
        raise HTTPException(400, "invalid userId")
    path = photo_path(user_id)
    path.write_bytes(content)
    return path.name


def delete_photo(user_id: str) -> None:
    path = photo_path(user_id)
    if not path.exists():
        raise HTTPException(404, "photo not found")
    path.unlink()


def save_screenshot(employee_id: str, user_id: str, content: bytes) -> str:
    if not re.fullmatch(r"[0-9a-fA-F-]{36}", employee_id):
        raise HTTPException(400, "invalid employeeId")
    if not re.fullmatch(r"[0-9a-fA-F-]{36}", user_id):
        raise HTTPException(400, "invalid userId")

    screenshot_id = str(uuid.uuid4())
    directory = screenshot_dir(employee_id)
    directory.mkdir(parents=True, exist_ok=True)
    filename = f"{employee_id}-{screenshot_id}.png"
    (directory / filename).write_bytes(content)
    return filename


def resolve_screenshot_path(employee_id: str, filename: str) -> Path:
    if not re.fullmatch(r"[0-9a-fA-F-]{36}", employee_id):
        raise HTTPException(400, "invalid employeeId")
    if not SCREENSHOT_NAME_RE.fullmatch(filename):
        raise HTTPException(400, "invalid filename")
    if not filename.startswith(f"{employee_id}-"):
        raise HTTPException(400, "filename does not belong to employee")

    path = screenshot_dir(employee_id) / filename
    if not path.exists():
        raise HTTPException(404, "screenshot not found")
    return path


def build_screenshots_archive(employee_id: str, day: date) -> tuple[bytes, str]:
    filenames = list_screenshots(employee_id, day)
    if not filenames:
        raise HTTPException(404, "no screenshots for this date")

    buffer = io.BytesIO()
    with zipfile.ZipFile(buffer, "w", zipfile.ZIP_DEFLATED) as archive:
        for name in filenames:
            path = screenshot_dir(employee_id) / name
            archive.writestr(name, path.read_bytes())

    archive_name = f"{employee_id}-{day.isoformat()}.zip"
    return buffer.getvalue(), archive_name


def list_screenshots(employee_id: str, day: date) -> list[str]:
    directory = screenshot_dir(employee_id)
    if not directory.exists():
        return []

    result: list[str] = []
    for path in directory.glob(f"{employee_id}-*.png"):
        if not SCREENSHOT_NAME_RE.fullmatch(path.name):
            continue
        modified = datetime.fromtimestamp(path.stat().st_mtime, tz=timezone.utc).date()
        if modified == day:
            result.append(path.name)
    return sorted(result)


def parse_user_id_from_filename(filename: str | None) -> str | None:
    if not filename:
        return None
    name = Path(filename).name
    if name.endswith(".png"):
        user_id = name[:-4]
        if re.fullmatch(r"[0-9a-fA-F-]{36}", user_id):
            return user_id
    return None
