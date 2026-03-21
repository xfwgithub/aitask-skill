from fastapi import FastAPI, HTTPException, Request, Depends
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from fastapi.responses import FileResponse, HTMLResponse
from fastapi.templating import Jinja2Templates
from contextlib import asynccontextmanager
import uvicorn
import os
from pathlib import Path

from app.core.config import settings
from app.db.database import init_db, get_db
from app.api import tasks, projects, reports
from app.models.models import Task, Project
from sqlalchemy import select
from typing import Optional


@asynccontextmanager
async def lifespan(app: FastAPI):
    """应用生命周期管理"""
    # 启动时初始化数据库
    await init_db()
    print("数据库初始化完成")
    yield
    # 关闭时清理资源
    print("应用关闭")


# 创建 FastAPI 应用
app = FastAPI(
    title=settings.APP_NAME,
    version=settings.APP_VERSION,
    lifespan=lifespan,
)

# 配置 CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # 生产环境应该限制具体域名
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 配置模板和静态文件
templates_dir = Path(__file__).parent / "templates"
templates = Jinja2Templates(directory=str(templates_dir))
app.mount("/static", StaticFiles(directory=str(Path(__file__).parent / "static")), name="static")


# 注册路由
app.include_router(projects.router, prefix=f"{settings.API_PREFIX}/projects", tags=["项目管理"])
app.include_router(reports.router, prefix=f"{settings.API_PREFIX}/reports", tags=["报表管理"])


@app.get("/api")
async def api_root():
    """API 根路径"""
    return {
        "name": settings.APP_NAME,
        "version": settings.APP_VERSION,
        "docs": "/docs",
        "health": "/health",
    }


@app.get("/health")
async def health_check():
    """健康检查"""
    return {"status": "ok"}


# ===== HTMX 前端路由 =====

@app.get("/", response_class=HTMLResponse)
async def home(request: Request):
    """首页 - 仪表盘"""
    return templates.TemplateResponse("index.html", {"request": request})


@app.get("/tasks", response_class=HTMLResponse)
async def tasks_page(request: Request):
    """任务列表页"""
    return templates.TemplateResponse("tasks.html", {"request": request})


@app.get("/projects", response_class=HTMLResponse)
async def projects_page(request: Request):
    """项目列表页"""
    return templates.TemplateResponse("projects.html", {"request": request})


# ===== HTMX HTML 片段路由 =====

@app.get("/api/html/tasks", response_class=HTMLResponse)
async def get_task_list_html(
    request: Request,
    db = Depends(get_db),
    status: Optional[str] = None,
    priority: Optional[int] = None,
    keyword: Optional[str] = None,
    page: int = 1,
    page_size: int = 20,
):
    """获取任务列表 HTML 片段"""
    query = select(Task)
    
    if status:
        query = query.where(Task.status == status)
    if priority:
        query = query.where(Task.priority == priority)
    if keyword:
        query = query.where(Task.title.ilike(f"%{keyword}%"))
    
    offset = (page - 1) * page_size
    query = query.offset(offset).limit(page_size)
    
    result = await db.execute(query)
    tasks = result.scalars().all()
    
    return templates.TemplateResponse("task_list.html", {"request": request, "tasks": tasks})


@app.get("/api/html/projects", response_class=HTMLResponse)
async def get_project_list_html(
    request: Request,
    db = Depends(get_db),
    keyword: Optional[str] = None,
):
    """获取项目列表 HTML 片段"""
    query = select(Project)
    
    if keyword:
        query = query.where(Project.name.ilike(f"%{keyword}%"))
    
    result = await db.execute(query)
    projects = result.scalars().all()
    
    return templates.TemplateResponse("project_list.html", {"request": request, "projects": projects})


if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host=settings.SERVER_HOST,
        port=settings.SERVER_PORT,
        reload=settings.DEBUG,
    )
