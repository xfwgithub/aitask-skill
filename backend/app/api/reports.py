from fastapi import APIRouter, Depends
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func, and_
from datetime import datetime, timedelta

from app.db.database import get_db
from app.models import Task

router = APIRouter()


@router.get("/overdue")
async def get_overdue_tasks(db: AsyncSession = Depends(get_db)):
    """获取逾期任务"""
    now = datetime.now()
    
    query = select(Task).where(
        Task.due_date < now,
        Task.status.not_in(['done', 'cancelled']),
        Task.is_deleted == False,
    )
    result = await db.execute(query)
    tasks = result.scalars().all()
    
    return {"total": len(tasks), "items": tasks}


@router.get("/due-soon")
async def get_due_soon_tasks(db: AsyncSession = Depends(get_db)):
    """获取即将到期任务（3 天内）"""
    now = datetime.now()
    three_days_later = now + timedelta(days=3)
    
    query = select(Task).where(
        Task.due_date >= now,
        Task.due_date <= three_days_later,
        Task.status.not_in(['done', 'cancelled']),
        Task.is_deleted == False,
    )
    result = await db.execute(query)
    tasks = result.scalars().all()
    
    return {"total": len(tasks), "items": tasks}


@router.get("/stats")
async def get_dashboard_stats(db: AsyncSession = Depends(get_db)):
    """获取仪表盘统计"""
    # 总任务数
    total_query = select(func.count()).select_from(Task).where(Task.is_deleted == False)
    total_result = await db.execute(total_query)
    total = total_result.scalar()
    
    # 按状态统计
    status_query = select(Task.status, func.count()).group_by(Task.status).where(Task.is_deleted == False)
    status_result = await db.execute(status_query)
    status_stats = {row[0]: row[1] for row in status_result.all()}
    
    # 按优先级统计
    priority_query = select(Task.priority, func.count()).group_by(Task.priority).where(Task.is_deleted == False)
    priority_result = await db.execute(priority_query)
    priority_stats = {row[0]: row[1] for row in priority_result.all()}
    
    return {
        "total": total,
        "by_status": status_stats,
        "by_priority": priority_stats,
    }
