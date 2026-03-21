from sqlalchemy import Column, String, Integer, DateTime, ForeignKey, Boolean, Text, Float
from sqlalchemy.orm import relationship
from sqlalchemy.sql import func
from app.db.database import Base
import uuid


def generate_uuid():
    return str(uuid.uuid4())


class Project(Base):
    """项目表"""
    __tablename__ = "projects"
    
    id = Column(Integer, primary_key=True, index=True)
    uuid = Column(String(36), unique=True, index=True, default=generate_uuid)
    name = Column(String(200), nullable=False)
    description = Column(Text, nullable=True)
    status = Column(String(50), default="active")  # active/archived
    color = Column(String(20), nullable=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    
    # 关系
    tasks = relationship("Task", back_populates="project", cascade="all, delete-orphan")


class Task(Base):
    """任务表"""
    __tablename__ = "tasks"
    
    id = Column(Integer, primary_key=True, index=True)
    uuid = Column(String(36), unique=True, index=True, default=generate_uuid)
    
    # 基本信息
    title = Column(String(500), nullable=False, index=True)
    description = Column(Text, nullable=True)
    status = Column(String(50), default="pending", index=True)  # pending/agent_working/pending_review/reviewing/done/cancelled
    priority = Column(Integer, default=3)  # 1-Critical/2-High/3-Medium/4-Low
    
    # 分类
    project_id = Column(Integer, ForeignKey("projects.id"), nullable=True, index=True)
    parent_task_id = Column(Integer, ForeignKey("tasks.id"), nullable=True, index=True)
    
    # 人员（简化：只记录负责人名称，不关联用户表）
    assignee_name = Column(String(100), nullable=True)  # 负责人名称
    
    # Agent 执行（只记录 Agent 信息，不关联 Agent 表）
    agent_type = Column(String(100), nullable=True)  # writer/reviewer/researcher
    agent_model = Column(String(200), nullable=True)  # qwen3.5-plus
    agent_prompt_id = Column(String(100), nullable=True)
    agent_input = Column(String, nullable=True)  # JSON
    agent_output = Column(String, nullable=True)  # JSON
    agent_log = Column(Text, nullable=True)
    agent_started_at = Column(DateTime(timezone=True), nullable=True)
    agent_completed_at = Column(DateTime(timezone=True), nullable=True)
    estimated_tokens = Column(Integer, nullable=True)
    actual_tokens = Column(Integer, nullable=True)
    
    # 人工审核
    review_started_at = Column(DateTime(timezone=True), nullable=True)
    review_completed_at = Column(DateTime(timezone=True), nullable=True)
    review_comments = Column(Text, nullable=True)
    review_score = Column(Integer, nullable=True)  # 1-5
    sla_deadline = Column(DateTime(timezone=True), nullable=True)
    due_date = Column(DateTime(timezone=True), nullable=True)
    
    # 依赖和自定义字段
    tags = Column(String, nullable=True)  # JSON 数组
    dependencies = Column(String, nullable=True)  # JSON 数组
    custom_fields = Column(String, nullable=True)  # JSON
    
    # 系统字段
    created_at = Column(DateTime(timezone=True), server_default=func.now(), index=True)
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())
    version = Column(Integer, default=0)
    is_deleted = Column(Boolean, default=False, index=True)  # 软删除
    
    # 关系
    project = relationship("Project", back_populates="tasks")
    subtasks = relationship("Task", backref="parent", remote_side=[id])
    comments = relationship("Comment", back_populates="task", cascade="all, delete-orphan")
    attachments = relationship("Attachment", back_populates="task", cascade="all, delete-orphan")
    history = relationship("TaskHistory", back_populates="task", cascade="all, delete-orphan")


class TaskHistory(Base):
    """任务历史表"""
    __tablename__ = "task_history"
    
    id = Column(Integer, primary_key=True, index=True)
    task_id = Column(Integer, ForeignKey("tasks.id"), nullable=False, index=True)
    action = Column(String(100), nullable=False)  # create/update/delete/change_status
    old_value = Column(String, nullable=True)  # JSON
    new_value = Column(String, nullable=True)  # JSON
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    
    # 关系
    task = relationship("Task", back_populates="history")


class Comment(Base):
    """评论表"""
    __tablename__ = "comments"
    
    id = Column(Integer, primary_key=True, index=True)
    uuid = Column(String(36), unique=True, index=True, default=generate_uuid)
    task_id = Column(Integer, ForeignKey("tasks.id"), nullable=False, index=True)
    author_name = Column(String(100), nullable=True)  # 评论者名称
    content = Column(Text, nullable=False)
    parent_comment_id = Column(Integer, ForeignKey("comments.id"), nullable=True)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    updated_at = Column(DateTime(timezone=True), server_default=func.now(), onupdate=func.now())
    
    # 关系
    task = relationship("Task", back_populates="comments")
    replies = relationship("Comment", backref="parent", remote_side=[id])


class Attachment(Base):
    """附件表"""
    __tablename__ = "attachments"
    
    id = Column(Integer, primary_key=True, index=True)
    uuid = Column(String(36), unique=True, index=True, default=generate_uuid)
    task_id = Column(Integer, ForeignKey("tasks.id"), nullable=False, index=True)
    file_name = Column(String(500), nullable=False)
    file_path = Column(String(1000), nullable=False)
    file_size = Column(Integer, nullable=True)
    mime_type = Column(String(200), nullable=True)
    uploaded_by_name = Column(String(100), nullable=True)  # 上传者名称
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    
    # 关系
    task = relationship("Task", back_populates="attachments")
