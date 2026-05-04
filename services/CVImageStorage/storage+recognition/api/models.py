from sqlalchemy import Column, String, DateTime, ForeignKey, Index, Text, Integer, LargeBinary
from sqlalchemy.orm import relationship
from database import Base
import uuid
from datetime import datetime


class Manager(Base):
    __tablename__ = "managers"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    name = Column(String, nullable=False)
    email = Column(String, unique=True, nullable=False)
    created_at = Column(DateTime, default=datetime.utcnow)

    employees = relationship("Employee", back_populates="manager", cascade="all, delete-orphan")


class Employee(Base):
    __tablename__ = "employees"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    name = Column(String, nullable=False)
    email = Column(String, unique=True, nullable=False)
    manager_id = Column(String, ForeignKey("managers.id"), nullable=True)
    embedding = Column(Text, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    manager = relationship("Manager", back_populates="employees")
    photos = relationship("Photo", back_populates="employee", cascade="all, delete-orphan")
    screenshots = relationship("Screenshot", back_populates="employee", cascade="all, delete-orphan")
    work_sessions = relationship("WorkSession", back_populates="employee", cascade="all, delete-orphan")


class Photo(Base):
    __tablename__ = "photos"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    user_id = Column(String, ForeignKey("employees.id"), index=True)
    filename = Column(String)
    content_type = Column(String, default="image/png")
    data = Column(LargeBinary)  # ← BYTEA: сами байты файла
    uploaded_at = Column(DateTime, default=datetime.utcnow)

    employee = relationship("Employee", back_populates="photos")


class Screenshot(Base):
    __tablename__ = "screenshots"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    employee_id = Column(String, ForeignKey("employees.id"), index=True)
    filename = Column(String)
    content_type = Column(String, default="image/png")
    data = Column(LargeBinary)  # ← BYTEA: сами байты файла
    cnt_mouse_clicks = Column(Integer, default=0)
    cnt_keyboard_clicks = Column(Integer, default=0)
    timestamp = Column(DateTime, default=datetime.utcnow, index=True)

    employee = relationship("Employee", back_populates="screenshots")

    __table_args__ = (Index('idx_employee_timestamp', 'employee_id', 'timestamp'),)


class WorkSession(Base):
    __tablename__ = "work_sessions"

    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    employee_id = Column(String, ForeignKey("employees.id"), index=True)
    start_time = Column(DateTime, nullable=False)
    end_time = Column(DateTime, nullable=True)
    total_time = Column(Integer, nullable=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    employee = relationship("Employee", back_populates="work_sessions")