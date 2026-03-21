from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func, and_
from typing import List, Optional
from datetime import datetime

from app.db.database import get_db
from app.models import Task, Project
from app.schemas.task import TaskCreate, TaskUpdate, TaskResponse, TaskListResponse
from app.utils import json_converter

router = APIRouter()


@router.post("/", response_model=TaskResponse)
async def create_task(task: TaskCreate, db: AsyncSession = Depends(get_db)):
    """创建任务"""
    db_task = Task(
        title=task.title,
        description=task.description,
        status=task.status,
        priority=task.priority,
        project_id=task.project_id,
        parent_task_id=task.parent_task_id,
        assignee_name=task.assignee_name,
        agent_type=task.agent_type,
        agent_model=task.agent_model,
        agent_prompt_id=task.agent_prompt_id,
        agent_input=json_converter.dumps(task.agent_input),
        estimated_tokens=task.estimated_tokens,
        sla_deadline=task.sla_deadline,
        due_date=task.due_date,
        tags=json_converter.dumps(task.tags),
        dependencies=json_converter.dumps(task.dependencies),
        custom_fields=json_converter.dumps(task.custom_fields),
    )
    
    db.add(db_task)
    await db.commit()
    await db.refresh(db_task)
    
    return db_task


@router.get("/", response_model=TaskListResponse)
async def get_tasks(
    page: int = Query(1, ge=1),
    page_size: int = Query(20, ge=1, le=100),
    status: Optional[str] = None,
    priority: Optional[int] = None,
    project_id: Optional[int] = None,
    assignee_name: Optional[str] = None,
    agent_type: Optional[str] = None,
    keyword: Optional[str] = None,
    db: AsyncSession = Depends(get_db),
):
    """查询任务列表（支持多维度筛选）"""
    # 构建查询条件
    conditions = [Task.is_deleted == False]
    
    if status:
        conditions.append(Task.status == status)
    if priority:
        conditions.append(Task.priority == priority)
    if project_id:
        conditions.append(Task.project_id == project_id)
    if assignee_name:
        conditions.append(Task.assignee_name == assignee_name)
    if agent_type:
        conditions.append(Task.agent_type == agent_type)
    if keyword:
        conditions.append(Task.title.ilike(f"%{keyword}%"))
    
    # 查询总数
    total_query = select(func.count()).select_from(Task).where(and_(*conditions))
    total_result = await db.execute(total_query)
    total = total_result.scalar()
    
    # 分页查询
    query = select(Task).where(and_(*conditions)).offset((page - 1) * page_size).limit(page_size)
    result = await db.execute(query)
    tasks = result.scalars().all()
    
    return TaskListResponse(
        total=total,
        items=tasks,
        page=page,
        page_size=page_size,
    )


@router.get("/{task_uuid}", response_model=TaskResponse)
async def get_task(task_uuid: str, db: AsyncSession = Depends(get_db)):
    """获取任务详情"""
    query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
    result = await db.execute(query)
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(status_code=404, detail="任务不存在")
    
    return task


@router.put("/{task_uuid}", response_model=TaskResponse)
async def update_task(task_uuid: str, task_update: TaskUpdate, db: AsyncSession = Depends(get_db)):
    """更新任务"""
    query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
    result = await db.execute(query)
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(status_code=404, detail="任务不存在")
    
    # 更新字段
    update_data = task_update.model_dump(exclude_unset=True)
    for field, value in update_data.items():
        if value is not None:
            if field in ["agent_input", "agent_output", "tags", "dependencies", "custom_fields"]:
                setattr(task, field, json_converter.dumps(value))
            else:
                setattr(task, field, value)
    
    task.version += 1  # 版本号 +1（乐观锁）
    
    await db.commit()
    await db.refresh(task)
    
    return task


@router.delete("/{task_uuid}")
async def delete_task(task_uuid: str, hard: bool = Query(False), db: AsyncSession = Depends(get_db)):
    """删除任务（支持软删除和硬删除）"""
    query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
    result = await db.execute(query)
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(status_code=404, detail="任务不存在")
    
    if hard:
        # 硬删除
        await db.delete(task)
    else:
        # 软删除
        task.is_deleted = True
    
    await db.commit()
    
    return {"message": "删除成功"}


@router.post("/{task_uuid}/status")
async def update_task_status(
    task_uuid: str,
    status: str = Query(..., description="目标状态"),
    db: AsyncSession = Depends(get_db),
):
    """更新任务状态"""
    # 验证状态流转
    valid_transitions = {
        'pending': ['agent_working', 'cancelled'],
        'agent_working': ['pending_review', 'done', 'cancelled'],
        'pending_review': ['reviewing', 'cancelled'],
        'reviewing': ['agent_working', 'done', 'cancelled'],
        'done': ['agent_working'],
        'cancelled': ['pending'],
    }
    
    query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
    result = await db.execute(query)
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(status_code=404, detail="任务不存在")
    
    if status not in valid_transitions.get(task.status, []):
        raise HTTPException(
            status_code=400,
            detail=f"不允许的状态流转：{task.status} -> {status}",
        )
    
    task.status = status
    task.version += 1
    
    # 根据状态设置时间字段
    now = datetime.now()
    if status == 'agent_working':
        task.agent_started_at = now
    elif status == 'pending_review':
        task.agent_completed_at = now
    elif status == 'reviewing':
        task.review_started_at = now
    elif status == 'done':
        task.review_completed_at = now
    
    await db.commit()
    await db.refresh(task)
    
    return {"message": "状态更新成功", "task": task}


@router.get("/{task_uuid}/subtasks")
async def get_subtasks(task_uuid: str, db: AsyncSession = Depends(get_db)):
    """获取子任务列表"""
    query = select(Task).where(Task.parent_task_id == Task.id, Task.uuid == task_uuid)
    result = await db.execute(query)
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(status_code=404, detail="任务不存在")
    
    subtasks_query = select(Task).where(
        Task.parent_task_id == task.id,
        Task.is_deleted == False,
    )
    subtasks_result = await db.execute(subtasks_query)
    subtasks = subtasks_result.scalars().all()
    
    return {"total": len(subtasks), "items": subtasks}


@router.get("/{task_uuid}/comments")
async def get_comments(task_uuid: str, db: AsyncSession = Depends(get_db)):
    """获取评论列表"""
    from app.models.models import Comment
    
    query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
    result = await db.execute(query)
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(status_code=404, detail="任务不存在")
    
    comments_query = select(Comment).where(Comment.task_id == task.id).order_by(Comment.created_at)
    comments_result = await db.execute(comments_query)
    comments = comments_result.scalars().all()
    
    return {"total": len(comments), "items": comments}


@router.post("/{task_uuid}/comments")
async def add_comment(
    task_uuid: str,
    content: str = Query(...),
    parent_comment_id: Optional[int] = None,
    db: AsyncSession = Depends(get_db),
):
    """添加评论"""
    from app.models.models import Comment
    
    query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
    result = await db.execute(query)
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(status_code=404, detail="任务不存在")
    
    comment = Comment(
        task_id=task.id,
        user_id=1,  # TODO: 从当前用户获取
        content=content,
        parent_comment_id=parent_comment_id,
    )
    
    db.add(comment)
    await db.commit()
    await db.refresh(comment)
    
    return comment


@router.get("/{task_uuid}/history")
async def get_task_history(task_uuid: str, db: AsyncSession = Depends(get_db)):
    """获取任务历史"""
    from app.models.models import TaskHistory
    
    query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
    result = await db.execute(query)
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(status_code=404, detail="任务不存在")
    
    history_query = select(TaskHistory).where(TaskHistory.task_id == task.id).order_by(TaskHistory.created_at.desc())
    history_result = await db.execute(history_query)
    history = history_result.scalars().all()
    
    return {"total": len(history), "items": history}


@router.get("/{task_uuid}/attachments")
async def get_attachments(task_uuid: str, db: AsyncSession = Depends(get_db)):
    """获取附件列表"""
    from app.models.models import Attachment
    
    query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
    result = await db.execute(query)
    task = result.scalar_one_or_none()
    
    if not task:
        raise HTTPException(status_code=404, detail="任务不存在")
    
    attachments_query = select(Attachment).where(Attachment.task_id == task.id)
    attachments_result = await db.execute(attachments_query)
    attachments = attachments_result.scalars().all()
    
    return {"total": len(attachments), "items": attachments}
