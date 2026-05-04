from sqlalchemy.orm import Session
from models import Photo, Screenshot
import uuid

def save_photo_to_db(
    db: Session,
    employee_id: str,
    file_content: bytes,
    filename: str,
    content_type: str = "image/png"
) -> Photo:
    """Сохраняет фото в БД (BYTEA), возвращает объект Photo"""
    photo = Photo(
        id=str(uuid.uuid4()),
        user_id=employee_id,
        filename=filename,
        content_type=content_type,
        data=file_content
    )
    db.add(photo)
    return photo

def save_screenshot_to_db(
    db: Session,
    employee_id: str,
    file_content: bytes,
    filename: str,
    cnt_mouse_clicks: int = 0,
    cnt_keyboard_clicks: int = 0,
    content_type: str = "image/png"
) -> Screenshot:
    """Сохраняет скриншот в БД (BYTEA), возвращает объект Screenshot"""
    screenshot = Screenshot(
        id=str(uuid.uuid4()),
        employee_id=employee_id,
        filename=filename,
        content_type=content_type,
        data=file_content,
        cnt_mouse_clicks=cnt_mouse_clicks,
        cnt_keyboard_clicks=cnt_keyboard_clicks
    )
    db.add(screenshot)
    return screenshot

def get_photo_data(db: Session, photo_id: str) -> bytes | None:
    """Получает байты фото по ID"""
    photo = db.query(Photo).filter(Photo.id == photo_id).first()
    return photo.data if photo else None

def get_screenshot_data(db: Session, screenshot_id: str) -> bytes | None:
    """Получает байты скриншота по ID"""
    screenshot = db.query(Screenshot).filter(Screenshot.id == screenshot_id).first()
    return screenshot.data if screenshot else None