from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select
from typing import List

from app.db.database import get_db
from app.models import Project
from app.schemas.project import ProjectCreate, ProjectUpdate, ProjectResponse

router = APIRouter()


@router.post("/", response_model=ProjectResponse)
async def create_project(project: ProjectCreate, db: AsyncSession = Depends(get_db)):
    """创建项目"""
    db_project = Project(
        name=project.name,
        description=project.description,
        color=project.color,
    )
    
    db.add(db_project)
    await db.commit()
    await db.refresh(db_project)
    
    return db_project


@router.get("/", response_model=List[ProjectResponse])
async def get_projects(db: AsyncSession = Depends(get_db)):
    """获取项目列表"""
    query = select(Project)
    result = await db.execute(query)
    projects = result.scalars().all()
    
    return projects


@router.get("/{project_uuid}", response_model=ProjectResponse)
async def get_project(project_uuid: str, db: AsyncSession = Depends(get_db)):
    """获取项目详情"""
    query = select(Project).where(Project.uuid == project_uuid)
    result = await db.execute(query)
    project = result.scalar_one_or_none()
    
    if not project:
        raise HTTPException(status_code=404, detail="项目不存在")
    
    return project


@router.put("/{project_uuid}", response_model=ProjectResponse)
async def update_project(project_uuid: str, project_update: ProjectUpdate, db: AsyncSession = Depends(get_db)):
    """更新项目"""
    query = select(Project).where(Project.uuid == project_uuid)
    result = await db.execute(query)
    project = result.scalar_one_or_none()
    
    if not project:
        raise HTTPException(status_code=404, detail="项目不存在")
    
    update_data = project_update.model_dump(exclude_unset=True)
    for field, value in update_data.items():
        if value is not None:
            setattr(project, field, value)
    
    await db.commit()
    await db.refresh(project)
    
    return project


@router.delete("/{project_uuid}")
async def delete_project(project_uuid: str, db: AsyncSession = Depends(get_db)):
    """删除项目"""
    query = select(Project).where(Project.uuid == project_uuid)
    result = await db.execute(query)
    project = result.scalar_one_or_none()
    
    if not project:
        raise HTTPException(status_code=404, detail="项目不存在")
    
    await db.delete(project)
    await db.commit()
    
    return {"message": "删除成功"}


@router.get("/{project_uuid}/stats")
async def get_project_stats(project_uuid: str, db: AsyncSession = Depends(get_db)):
    """获取项目统计信息"""
    from sqlalchemy import func
    from app.models.models import Task
    
    query = select(Project).where(Project.uuid == project_uuid)
    result = await db.execute(query)
    project = result.scalar_one_or_none()
    
    if not project:
        raise HTTPException(status_code=404, detail="项目不存在")
    
    # 统计任务数
    total_query = select(func.count()).select_from(Task).where(Task.project_id == project.id, Task.is_deleted == False)
    total_result = await db.execute(total_query)
    total = total_result.scalar()
    
    # 统计完成任务数
    done_query = select(func.count()).select_from(Task).where(
        Task.project_id == project.id,
        Task.status == 'done',
        Task.is_deleted == False,
    )
    done_result = await db.execute(done_query)
    done = done_result.scalar()
    
    completion_rate = (done / total * 100) if total > 0 else 0
    
    return {
        "project_id": project.uuid,
        "total_tasks": total,
        "done_tasks": done,
        "completion_rate": round(completion_rate, 2),
    }
