from fastapi import APIRouter, File, HTTPException, Request, UploadFile

from face_store import register_photo, remove_user
from file_storage import delete_photo, parse_user_id_from_filename, save_photo

router = APIRouter(tags=["Photo"])


async def _read_body(request: Request, upload: UploadFile | None) -> tuple[bytes, str | None]:
    if upload is not None:
        content = await upload.read()
        return content, upload.filename

    content = await request.body()
    if not content:
        raise HTTPException(400, "empty body")
    return content, request.headers.get("X-Filename")


@router.post("/photo")
async def upload_photo(request: Request, file: UploadFile | None = File(None)):
    content, filename = await _read_body(request, file)
    user_id = parse_user_id_from_filename(filename)
    if user_id is None:
        user_id = request.headers.get("X-User-Id")
    if user_id is None:
        raise HTTPException(400, "userId required in filename (userId.png) or X-User-Id header")

    saved_name = save_photo(user_id, content)
    try:
        register_photo(user_id, content)
    except ValueError as exc:
        raise HTTPException(400, str(exc)) from exc

    return {"userId": user_id, "filename": saved_name}


@router.delete("/photo/{user_id}")
async def remove_photo(user_id: str):
    delete_photo(user_id)
    remove_user(user_id)
    return {"userId": user_id, "deleted": True}
