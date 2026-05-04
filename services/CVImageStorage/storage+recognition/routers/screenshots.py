import sys
from pathlib import Path
sys.path.append(str(Path(__file__).parent.parent))

from fastapi import APIRouter, Depends, HTTPException, Query, UploadFile, File
from sqlalchemy.orm import Session
from sqlalchemy import func
from datetime import datetime
from typing import Optional
from database import get_db
from models import Employee, Screenshot
from storage import save_screenshot_to_db

router = APIRouter(prefix="/screenshots", tags=["Screenshots"])


@router.post("")
async def upload_screenshot(
        employee_id: str = Query(...),
        file: UploadFile = File(...),
        cnt_mouse_clicks: int = Query(0),
        cnt_keyboard_clicks: int = Query(0),
        db: Session = Depends(get_db)
):
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    if not file.filename.lower().endswith(('.png', '.jpg', '.jpeg')):
        raise HTTPException(400, "Invalid format")

    content = await file.read()
    if len(content) > 10 * 1024 * 1024:
        raise HTTPException(400, "File too large (max 10MB)")

    # === Сохраняем в БД (не на диск!) ===
    screenshot = save_screenshot_to_db(
        db=db,
        employee_id=employee_id,
        file_content=content,
        filename=file.filename,
        cnt_mouse_clicks=cnt_mouse_clicks,
        cnt_keyboard_clicks=cnt_keyboard_clicks,
        content_type=file.content_type or "image/png"
    )
    db.commit()
    db.refresh(screenshot)

    return {
        "id": screenshot.id,
        "employee_id": employee_id,
        "filename": file.filename,
        "cnt_mouse_clicks": cnt_mouse_clicks,
        "cnt_keyboard_clicks": cnt_keyboard_clicks,
        "timestamp": screenshot.timestamp.isoformat(),
        "stored_in": "PostgreSQL (BYTEA)"
    }


@router.get("/employee/{employee_id}")
async def get_employee_screenshots(
        employee_id: str,
        date: Optional[str] = Query(None),
        limit: int = Query(100),
        db: Session = Depends(get_db)
):
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    query = db.query(Screenshot).filter(Screenshot.employee_id == employee_id)

    if date:
        try:
            target = datetime.strptime(date, "%Y-%m-%d").date()
            query = query.filter(func.date(Screenshot.timestamp) == target)
        except ValueError:
            raise HTTPException(400, "Invalid date format")

    screenshots = query.order_by(Screenshot.timestamp.desc()).limit(limit).all()
    return {
        "employee_id": employee_id,
        "count": len(screenshots),
        "screenshots": [
            {
                "id": s.id,
                "filename": s.filename,
                "content_type": s.content_type,
                "size_bytes": len(s.data) if s.data else 0,
                "cnt_mouse_clicks": s.cnt_mouse_clicks,
                "cnt_keyboard_clicks": s.cnt_keyboard_clicks,
                "timestamp": s.timestamp.isoformat()
            }
            for s in screenshots
        ]
    }


@router.get("/screenshot/{screenshot_id}/download")
async def download_screenshot(screenshot_id: str, db: Session = Depends(get_db)):
    """Скачать скриншот как файл"""
    from fastapi.responses import Response
    screenshot = db.query(Screenshot).filter(Screenshot.id == screenshot_id).first()
    if not screenshot:
        raise HTTPException(404, "Screenshot not found")

    return Response(
        content=screenshot.data,
        media_type=screenshot.content_type,
        headers={"Content-Disposition": f"attachment; filename={screenshot.filename}"}
    )