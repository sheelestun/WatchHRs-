from datetime import date

from fastapi import APIRouter, File, HTTPException, Request, UploadFile
from fastapi.responses import FileResponse, Response

from file_storage import (
    build_screenshots_archive,
    list_screenshots,
    resolve_screenshot_path,
    save_screenshot,
)

router = APIRouter(tags=["Screenshot"])


@router.post("/screenshot/{employee_id}")
async def upload_screenshot(
    employee_id: str,
    request: Request,
    file: UploadFile | None = File(None),
):
    content: bytes
    filename: str | None

    if file is not None:
        content = await file.read()
        filename = file.filename
    else:
        content = await request.body()
        filename = request.headers.get("X-Filename")

    if not content:
        raise HTTPException(400, "empty body")

    user_id = request.headers.get("X-User-Id")
    if user_id is None and filename:
        name = filename.removesuffix(".png") if filename.endswith(".png") else filename
        if "-" in name:
            user_id = name.split("-", 1)[0]

    if user_id is None:
        raise HTTPException(400, "X-User-Id header or userid-photoid.png filename required")

    saved_name = save_screenshot(employee_id, user_id, content)
    return {"employeeId": employee_id, "filename": saved_name}


@router.get("/screenshot/{employee_id}/file/{filename}")
async def download_screenshot(employee_id: str, filename: str):
    path = resolve_screenshot_path(employee_id, filename)
    return FileResponse(path, media_type="image/png", filename=filename)


@router.get("/screenshot/{employee_id}/{day}/archive")
async def download_screenshots_archive(employee_id: str, day: str):
    try:
        parsed_day = date.fromisoformat(day)
    except ValueError as exc:
        raise HTTPException(400, "invalid date, use YYYY-MM-DD") from exc

    content, archive_name = build_screenshots_archive(employee_id, parsed_day)
    headers = {"Content-Disposition": f'attachment; filename="{archive_name}"'}
    return Response(content=content, media_type="application/zip", headers=headers)


@router.get("/screenshot/{employee_id}/{day}")
async def get_screenshots(employee_id: str, day: str):
    try:
        parsed_day = date.fromisoformat(day)
    except ValueError as exc:
        raise HTTPException(400, "invalid date, use YYYY-MM-DD") from exc

    filenames = list_screenshots(employee_id, parsed_day)
    return {
        "screenshots": [
            {
                "filename": name,
                "url": f"/screenshot/{employee_id}/file/{name}",
            }
            for name in filenames
        ]
    }
