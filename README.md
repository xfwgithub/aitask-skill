# Task Management Skill 📋

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Python 3.11+](https://img.shields.io/badge/python-3.11+-blue.svg)](https://www.python.org/downloads/)

让 AI 助手帮你高效管理任务的标准化技能包！

**核心特性**：
- ✅ **AI Skill 支持** - 自然语言操作任务
- ✅ **快速创建任务** - 支持自然语言和结构化创建
- ✅ **智能查询任务** - 多维度筛选和搜索
- ✅ **任务统计** - 实时统计和数据分析
- ✅ **本地存储** - SQLite 数据库，安全可靠
- ✅ **Agent 支持** - 可指定 Agent 类型自动执行
- ✅ **RESTful API** - 完整的 API 接口

## 技术栈

- **FastAPI** - 高性能 Web 框架
- **SQLAlchemy** - ORM 框架
- **SQLite** - 轻量级数据库
- **Pydantic** - 数据验证

## 📦 安装

### 方式 1：使用 skills CLI（推荐）

```bash
npx skills add xfwgithub/aitask-skill
```

### 方式 2：手动安装

```bash
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill
pip install -r src/requirements.txt
```

### 方式 3：复制文件

```bash
# 复制到 Claude Skills 目录
cp -r . ~/.claude/skills/aitask-skill/

# 或复制到 Cursor Rules 目录
cp -r . ~/.cursor/rules/aitask-skill/
```

## 🚀 快速开始

### 1. 安装依赖

```bash
cd src
pip install -r requirements.txt
```

### 2. 配置服务（可选）

编辑 `src/.env` 文件：

```bash
# 服务器配置
SERVER_PORT=8000        # 端口号
SERVER_HOST=0.0.0.0     # 监听地址
```

### 3. 初始化数据库

```bash
cd src
python start.py --init-db
```

### 4. 启动服务

```bash
# 启动服务（最常用）
python src/main.py

# 开发模式（热重载）
python src/main.py --reload

# 自定义端口
python src/main.py --port 8080
```

### 5. 访问服务

- **API 文档**: http://localhost:8000/docs
- **健康检查**: http://localhost:8000/health

## 🎯 技能触发条件

当用户提到以下关键词或意图时触发此技能：

- **创建任务**："创建任务"、"新建任务"、"添加一个任务"、"帮我记个事"
- **查询任务**："查看任务"、"查询任务"、"我的任务"、"有什么任务"
- **更新任务**："更新任务"、"完成任务"、"取消任务"、"修改任务"
- **任务统计**："任务统计"、"任务概况"、"有多少任务"
- **任务详情**："任务详情"、"查看某个任务"

## 🤖 AI Skill 使用

### 自然语言示例

直接对 AI 助手说：

```
"帮我创建一个紧急任务，明天要完成项目评审"
"查询我所有待处理的任务"
"把任务'审查文档'标记为已完成"
"统计一下本周完成了多少任务"
```

### Python API 调用

```python
from src.skills.task_skill import TaskManagerSkill

# 初始化
skill = TaskManagerSkill()
await skill.initialize()

# 创建任务
task = await skill.create_task(
    title="审查第一卷第 10 章",
    description="检查时间线、人物、情节连贯性",
    priority=2,
    tags=["审查", "第一卷"],
    assignee_name="张三",
    agent_type="reviewer",
)

# 查询任务
tasks = await skill.query_tasks(status="pending")
print(f"待处理任务：{len(tasks['tasks'])} 个")

# 更新任务状态
await skill.update_task_status(task["uuid"], "done")

# 获取任务详情
detail = await skill.get_task_detail(task["uuid"])

# 获取统计信息
stats = await skill.get_dashboard_stats()
print(f"总任务数：{stats['total']}")
```

### RESTful API

#### 创建任务
```bash
curl -X POST "http://localhost:8000/api/v1/tasks/" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "测试任务",
    "description": "这是一个测试任务",
    "priority": 2,
    "tags": ["测试"]
  }'
```

#### 查询任务列表
```bash
curl "http://localhost:8000/api/v1/tasks/?status=pending&page=1&page_size=20"
```

#### 更新任务状态
```bash
curl -X POST "http://localhost:8000/api/v1/tasks/{uuid}/status?status=done"
```

## 📚 API 文档

完整的 API 文档请访问：
- **Swagger UI**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
