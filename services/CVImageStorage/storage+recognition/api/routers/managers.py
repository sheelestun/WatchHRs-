import sys
from pathlib import Path
sys.path.append(str(Path(__file__).parent.parent))

from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from database import get_db
from models import Manager

router = APIRouter(prefix="/managers", tags=["Managers"])


@router.post("")
async def create_manager(
        name: str = Query(...),
        email: str = Query(...),
        db: Session = Depends(get_db)
):
    existing = db.query(Manager).filter(Manager.email == email).first()
    if existing:
        raise HTTPException(400, "Manager with this email already exists")

    manager = Manager(name=name, email=email)
    db.add(manager)
    db.commit()
    db.refresh(manager)
    return {"id": manager.id, "name": manager.name, "email": manager.email}


@router.get("")
async def get_managers(db: Session = Depends(get_db)):
    managers = db.query(Manager).all()
    return [
        {"id": m.id, "name": m.name, "email": m.email, "employees_count": len(m.employees)}
        for m in managers
    ]