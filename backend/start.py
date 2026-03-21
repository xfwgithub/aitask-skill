#!/usr/bin/env python3
"""
任务管理系统 - 启动脚本

用法:
    python start.py              # 启动 Web 服务
    python start.py --skill      # 仅初始化 Skill
    python start.py --init-db    # 仅初始化数据库
"""
import asyncio
import argparse
import uvicorn
import sys
from pathlib import Path

from app.core.config import settings

# 添加项目根目录到 Python 路径
sys.path.insert(0, str(Path(__file__).parent))


async def init_database():
    """初始化数据库"""
    from app.db.database import init_db
    
    print("正在初始化数据库...")
    await init_db()
    print("数据库初始化完成！")


async def run_skill_demo():
    """运行 Skill 演示"""
    from app.skills import skill
    
    print("正在初始化 Skill...")
    await skill.initialize()
    print("Skill 初始化完成！")
    
    # 演示：创建任务
    print("\n=== 演示：创建任务 ===")
    result = await skill.create_task(
        title="测试任务",
        description="这是一个测试任务",
        priority=2,
        tags=["测试", "演示"],
    )
    print(f"创建结果：{result}")
    
    # 演示：查询任务
    print("\n=== 演示：查询任务 ===")
    result = await skill.query_tasks(limit=5)
    print(f"查询结果：{result}")


def main():
    parser = argparse.ArgumentParser(description="任务管理系统启动脚本")
    parser.add_argument(
        "--host",
        type=str,
        default="0.0.0.0",
        help="监听地址 (默认：0.0.0.0)",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=None,
        help="端口号 (默认：从 .env 文件读取，如未设置则为 8000)",
    )
    parser.add_argument(
        "--reload",
        action="store_true",
        help="启用热重载 (开发模式)",
    )
    parser.add_argument(
        "--init-db",
        action="store_true",
        help="仅初始化数据库",
    )
    parser.add_argument(
        "--skill",
        action="store_true",
        help="运行 Skill 演示",
    )
    parser.add_argument(
        "--frontend-only",
        action="store_true",
        help="仅启动前端开发服务器（需要前端目录存在）",
    )
    
    args = parser.parse_args()
    
    if args.init_db:
        # 仅初始化数据库
        asyncio.run(init_database())
        return
    
    if args.skill:
        # 运行 Skill 演示
        asyncio.run(run_skill_demo())
        return
    
    if args.frontend_only:
        # 仅启动前端开发服务器
        frontend_dir = Path(__file__).parent.parent / "frontend"
        if not frontend_dir.exists():
            print("错误：前端目录不存在")
            return
        
        print(f"启动前端开发服务器...")
        print(f"访问地址：http://{args.host}:{args.port}")
        print(f"按 Ctrl+C 停止服务\n")
        
        import subprocess
        subprocess.run(["npm", "run", "dev", "--", "--port", str(args.port)], cwd=frontend_dir)
        return
    
    # 确定端口（命令行参数 > .env 文件 > 默认值）
    port = args.port if args.port is not None else settings.SERVER_PORT
    host = args.host if args.host != "0.0.0.0" else settings.SERVER_HOST
    
    # 检查前端是否存在
    frontend_dir = Path(__file__).parent.parent / "frontend"
    has_frontend = frontend_dir.exists()
    
    # 启动 Web 服务
    print(f"启动任务管理系统 Web 服务...")
    print(f"访问地址：http://{host}:{port}")
    if has_frontend:
        print(f"前端界面：http://{host}:{port}/")
    else:
        print(f"API 文档：http://{host}:{port}/docs")
        print(f"（前端目录不存在，仅启动 API 服务）")
    print(f"按 Ctrl+C 停止服务\n")
    
    uvicorn.run(
        "main:app",
        host=host,
        port=port,
        reload=args.reload,
    )


if __name__ == "__main__":
    main()
