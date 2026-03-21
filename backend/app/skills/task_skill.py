"""
任务管理系统的 MCP Skill 接口

提供自然语言操作任务的能力
"""
from typing import Optional, List, Dict, Any
from datetime import datetime
import asyncio

from app.db.database import init_db
from app.models import Task, Project
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select


class TaskManagerSkill:
    """任务管理 Skill"""
    
    def __init__(self):
        self.db: Optional[AsyncSession] = None
    
    async def initialize(self):
        """初始化 Skill"""
        await init_db()
    
    async def _get_db(self) -> AsyncSession:
        """获取数据库会话"""
        from app.db.database import AsyncSessionLocal
        db = AsyncSessionLocal()
        try:
            yield db
        finally:
            await db.close()
    
    async def create_task(
        self,
        title: str,
        description: Optional[str] = None,
        priority: int = 3,
        project_name: Optional[str] = None,
        assignee_name: Optional[str] = None,
        agent_type: Optional[str] = None,
        agent_model: Optional[str] = None,
        tags: Optional[List[str]] = None,
        due_date: Optional[str] = None,
    ) -> Dict[str, Any]:
        """
        创建任务
        
        Args:
            title: 任务标题
            description: 任务描述
            priority: 优先级 (1-Critical/2-High/3-Medium/4-Low)
            project_name: 所属项目名称
            assignee_name: 负责人名称
            agent_type: Agent 类型 (writer/reviewer/researcher)
            agent_model: 使用的模型
            tags: 标签列表
            due_date: 截止日期 (YYYY-MM-DD 格式)
        
        Returns:
            创建的任务信息
        """
        from app.db.database import AsyncSessionLocal
        
        async with AsyncSessionLocal() as db:
            try:
                # 查询项目 ID
                project_id = None
                if project_name:
                    query = select(Project).where(Project.name == project_name)
                    result = await db.execute(query)
                    project = result.scalar_one_or_none()
                    if project:
                        project_id = project.id
                
                # 解析截止日期
                due_date_dt = None
                if due_date:
                    try:
                        due_date_dt = datetime.fromisoformat(due_date)
                    except ValueError:
                        pass
                
                # 创建任务
                task = Task(
                    title=title,
                    description=description,
                    priority=priority,
                    project_id=project_id,
                    assignee_name=assignee_name,
                    agent_type=agent_type,
                    agent_model=agent_model,
                    tags=str(tags) if tags else None,
                    due_date=due_date_dt,
                )
                
                db.add(task)
                await db.commit()
                await db.refresh(task)
                
                return {
                    "id": task.id,
                    "uuid": task.uuid,
                    "title": task.title,
                    "status": task.status,
                    "message": f"任务已创建：{task.title}",
                }
            finally:
                await db.close()
    
    async def query_tasks(
        self,
        status: Optional[str] = None,
        priority: Optional[int] = None,
        project_name: Optional[str] = None,
        assignee_name: Optional[str] = None,
        agent_type: Optional[str] = None,
        keyword: Optional[str] = None,
        limit: int = 20,
    ) -> Dict[str, Any]:
        """
        查询任务
        
        Args:
            status: 任务状态
            priority: 优先级
            project_name: 项目名称
            assignee_name: 负责人名称
            agent_type: Agent 类型
            keyword: 搜索关键词
            limit: 返回数量限制
        
        Returns:
            任务列表
        """
        from app.db.database import AsyncSessionLocal
        
        async with AsyncSessionLocal() as db:
            try:
                # 构建查询
                query = select(Task).where(Task.is_deleted == False)
                
                if status:
                    query = query.where(Task.status == status)
                if priority:
                    query = query.where(Task.priority == priority)
                if project_name:
                    project_query = select(Project.id).where(Project.name == project_name)
                    project_result = await db.execute(project_query)
                    project_id = project_result.scalar_one_or_none()
                    if project_id:
                        query = query.where(Task.project_id == project_id)
                if assignee_name:
                    query = query.where(Task.assignee_name == assignee_name)
                if agent_type:
                    query = query.where(Task.agent_type == agent_type)
                if keyword:
                    query = query.where(Task.title.ilike(f"%{keyword}%"))
                
                query = query.limit(limit)
                result = await db.execute(query)
                tasks = result.scalars().all()
                
                return {
                    "total": len(tasks),
                    "tasks": [
                        {
                            "id": t.id,
                            "uuid": t.uuid,
                            "title": t.title,
                            "status": t.status,
                            "priority": t.priority,
                        }
                        for t in tasks
                    ],
                }
            finally:
                await db.close()
    
    async def update_task_status(
        self,
        task_uuid: str,
        new_status: str,
    ) -> Dict[str, Any]:
        """
        更新任务状态
        
        Args:
            task_uuid: 任务 UUID
            new_status: 新状态
        
        Returns:
            更新结果
        """
        from app.db.database import AsyncSessionLocal
        
        async with AsyncSessionLocal() as db:
            try:
                # 查询任务
                query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
                result = await db.execute(query)
                task = result.scalar_one_or_none()
                
                if not task:
                    return {"error": "任务不存在"}
                
                # 验证状态流转
                valid_transitions = {
                    'pending': ['agent_working', 'cancelled'],
                    'agent_working': ['pending_review', 'done', 'cancelled'],
                    'pending_review': ['reviewing', 'cancelled'],
                    'reviewing': ['agent_working', 'done', 'cancelled'],
                    'done': ['agent_working'],
                    'cancelled': ['pending'],
                }
                
                if new_status not in valid_transitions.get(task.status, []):
                    return {
                        "error": f"不允许的状态流转：{task.status} -> {new_status}",
                    }
                
                task.status = new_status
                await db.commit()
                
                return {
                    "uuid": task.uuid,
                    "old_status": task.status,
                    "new_status": new_status,
                    "message": f"任务状态已更新为 {new_status}",
                }
            finally:
                await db.close()
    
    async def get_task_detail(self, task_uuid: str) -> Dict[str, Any]:
        """
        获取任务详情
        
        Args:
            task_uuid: 任务 UUID
        
        Returns:
            任务详情
        """
        from app.db.database import AsyncSessionLocal
        
        async with AsyncSessionLocal() as db:
            try:
                query = select(Task).where(Task.uuid == task_uuid, Task.is_deleted == False)
                result = await db.execute(query)
                task = result.scalar_one_or_none()
                
                if not task:
                    return {"error": "任务不存在"}
                
                return {
                    "id": task.id,
                    "uuid": task.uuid,
                    "title": task.title,
                    "description": task.description,
                    "status": task.status,
                    "priority": task.priority,
                    "project_id": task.project_id,
                    "assignee_id": task.assignee_id,
                    "agent_type": task.agent_type,
                    "agent_model": task.agent_model,
                    "tags": task.tags,
                    "created_at": str(task.created_at),
                    "updated_at": str(task.updated_at),
                }
            finally:
                await db.close()


# 全局 Skill 实例
skill = TaskManagerSkill()
