from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime


class ProjectBase(BaseModel):
    """项目基础 schema"""
    name: str = Field(..., min_length=1, max_length=200)
    description: Optional[str] = None
    color: Optional[str] = None


class ProjectCreate(ProjectBase):
    """创建项目请求"""
    pass


class ProjectUpdate(BaseModel):
    """更新项目请求"""
    name: Optional[str] = Field(None, min_length=1, max_length=200)
    description: Optional[str] = None
    status: Optional[str] = None
    color: Optional[str] = None


class ProjectResponse(ProjectBase):
    """项目响应"""
    id: int
    uuid: str
    status: str
    created_at: datetime
    
    class Config:
        from_attributes = True
