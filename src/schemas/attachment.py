from pydantic import BaseModel
from typing import Optional
from datetime import datetime


class AttachmentResponse(BaseModel):
    """附件响应"""
    id: int
    uuid: str
    task_id: int
    file_name: str
    file_path: str
    file_size: Optional[int] = None
    mime_type: Optional[str] = None
    uploaded_by: Optional[int] = None
    created_at: datetime
    
    class Config:
        from_attributes = True
