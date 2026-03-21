from pydantic import BaseModel, Field, field_validator
from typing import Optional, List, Dict, Any
from datetime import datetime
import json


def parse_json_field(value):
    """解析 JSON 字段"""
    if value is None:
        return None
    if isinstance(value, str):
        try:
            return json.loads(value)
        except (json.JSONDecodeError, TypeError):
            return None
    return value


class TaskBase(BaseModel):
    """任务基础 schema"""
    title: str = Field(..., min_length=1, max_length=500)
    description: Optional[str] = None
    status: str = Field(default="pending")
    priority: int = Field(default=3, ge=1, le=4)
    project_id: Optional[int] = None
    parent_task_id: Optional[int] = None
    tags: Optional[List[str]] = None
    dependencies: Optional[List[str]] = None
    custom_fields: Optional[Dict[str, Any]] = None


class TaskCreate(TaskBase):
    """创建任务请求"""
    assignee_name: Optional[str] = None  # 负责人名称
    agent_type: Optional[str] = None  # writer/reviewer/researcher
    agent_model: Optional[str] = None  # qwen3.5-plus
    agent_prompt_id: Optional[str] = None
    agent_input: Optional[Dict[str, Any]] = None
    estimated_tokens: Optional[int] = None
    sla_deadline: Optional[datetime] = None
    due_date: Optional[datetime] = None


class TaskUpdate(BaseModel):
    """更新任务请求"""
    title: Optional[str] = Field(None, min_length=1, max_length=500)
    description: Optional[str] = None
    status: Optional[str] = None
    priority: Optional[int] = Field(None, ge=1, le=4)
    project_id: Optional[int] = None
    parent_task_id: Optional[int] = None
    assignee_name: Optional[str] = None
    agent_type: Optional[str] = None
    agent_model: Optional[str] = None
    agent_prompt_id: Optional[str] = None
    agent_input: Optional[Dict[str, Any]] = None
    agent_output: Optional[Dict[str, Any]] = None
    agent_log: Optional[str] = None
    agent_started_at: Optional[datetime] = None
    agent_completed_at: Optional[datetime] = None
    estimated_tokens: Optional[int] = None
    actual_tokens: Optional[int] = None
    review_started_at: Optional[datetime] = None
    review_completed_at: Optional[datetime] = None
    review_comments: Optional[str] = None
    review_score: Optional[int] = Field(None, ge=1, le=5)
    sla_deadline: Optional[datetime] = None
    due_date: Optional[datetime] = None
    tags: Optional[List[str]] = None
    dependencies: Optional[List[str]] = None
    custom_fields: Optional[Dict[str, Any]] = None


class TaskResponse(TaskBase):
    """任务响应"""
    id: int
    uuid: str
    assignee_name: Optional[str] = None
    agent_type: Optional[str] = None
    agent_model: Optional[str] = None
    agent_input: Optional[Dict[str, Any]] = None
    agent_output: Optional[Dict[str, Any]] = None
    agent_log: Optional[str] = None
    agent_started_at: Optional[datetime] = None
    agent_completed_at: Optional[datetime] = None
    estimated_tokens: Optional[int] = None
    actual_tokens: Optional[int] = None
    review_started_at: Optional[datetime] = None
    review_completed_at: Optional[datetime] = None
    review_comments: Optional[str] = None
    review_score: Optional[int] = None
    sla_deadline: Optional[datetime] = None
    due_date: Optional[datetime] = None
    created_at: datetime
    updated_at: datetime
    version: int
    
    @field_validator('tags', 'dependencies', 'custom_fields', 'agent_input', 'agent_output', mode='before')
    @classmethod
    def parse_json_fields(cls, v):
        return parse_json_field(v)
    
    class Config:
        from_attributes = True


class TaskListResponse(BaseModel):
    """任务列表响应"""
    total: int
    items: List[TaskResponse]
    page: int
    page_size: int
