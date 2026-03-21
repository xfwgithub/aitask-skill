from pydantic import BaseModel, Field
from typing import Optional
from datetime import datetime


class CommentBase(BaseModel):
    """评论基础 schema"""
    content: str = Field(..., min_length=1)


class CommentCreate(CommentBase):
    """创建评论请求"""
    parent_comment_id: Optional[int] = None


class CommentResponse(CommentBase):
    """评论响应"""
    id: int
    uuid: str
    task_id: int
    user_id: Optional[int] = None
    parent_comment_id: Optional[int] = None
    created_at: datetime
    updated_at: datetime
    
    class Config:
        from_attributes = True
