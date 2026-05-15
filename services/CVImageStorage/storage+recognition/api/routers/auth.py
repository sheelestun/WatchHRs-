import cv2
import numpy as np
from fastapi import APIRouter, File, HTTPException, Request, UploadFile

from face_store import authenticate

router = APIRouter(tags=["Auth"])


@router.post("/auth")
async def auth_by_photo(request: Request, file: UploadFile | None = File(None)):
    if file is not None:
        content = await file.read()
    else:
        content = await request.body()

    if not content:
        raise HTTPException(400, "empty body")

    nparr = np.frombuffer(content, np.uint8)
    image = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
    if image is None:
        raise HTTPException(400, "invalid image")

    user_id, confidence = authenticate(image)
    if user_id is None:
        raise HTTPException(401, "authentication failed")

    return {
        "userId": user_id,
        "authenticated": True,
        "confidence": confidence,
    }
