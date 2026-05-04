import sys
from pathlib import Path
sys.path.append(str(Path(__file__).parent.parent))

from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.orm import Session
from sqlalchemy import func
from datetime import datetime
from typing import Optional

from database import get_db
from models import Employee, WorkSession

router = APIRouter(prefix="/work-sessions", tags=["Work Sessions"])


@router.post("/start")
async def start_work_session(employee_id: str = Query(...), db: Session = Depends(get_db)):
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    session = WorkSession(employee_id=employee_id, start_time=datetime.utcnow())
    db.add(session)
    db.commit()
    db.refresh(session)

    return {"session_id": session.id, "employee_id": employee_id, "start_time": session.start_time.isoformat(),
            "status": "started"}


@router.post("/{session_id}/end")
async def end_work_session(session_id: str, db: Session = Depends(get_db)):
    session = db.query(WorkSession).filter(WorkSession.id == session_id).first()
    if not session:
        raise HTTPException(404, "Session not found")

    end_time = datetime.utcnow()
    session.end_time = end_time
    session.total_time = int((end_time - session.start_time).total_seconds())
    db.commit()

    return {"session_id": session_id, "end_time": end_time.isoformat(), "total_time_seconds": session.total_time,
            "status": "completed"}


@router.get("/employee/{employee_id}")
async def get_employee_work_sessions(
        employee_id: str,
        date: Optional[str] = Query(None),
        db: Session = Depends(get_db)
):
    employee = db.query(Employee).filter(Employee.id == employee_id).first()
    if not employee:
        raise HTTPException(404, "Employee not found")

    query = db.query(WorkSession).filter(WorkSession.employee_id == employee_id)

    if date:
        try:
            target = datetime.strptime(date, "%Y-%m-%d").date()
            query = query.filter(func.date(WorkSession.start_time) == target)
        except ValueError:
            raise HTTPException(400, "Invalid date format")

    sessions = query.order_by(WorkSession.start_time.desc()).all()

    return {
        "employee_id": employee_id,
        "count": len(sessions),
        "sessions": [
            {"id": s.id, "start_time": s.start_time.isoformat(),
             "end_time": s.end_time.isoformat() if s.end_time else None, "total_time_seconds": s.total_time}
            for s in sessions
        ]
    }