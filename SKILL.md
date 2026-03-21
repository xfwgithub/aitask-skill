---
name: "task-management"
description: "AI 任务管理技能。当用户需要创建、查询、更新、删除任务或获取任务统计时触发此技能。"
version: "1.0.0"
author: "Your Name <your.email@example.com>"
license: "MIT"
keywords:
  - task-management
  - productivity
  - ai-skill
  - agent
engines:
  python: ">=3.11"
---

# Task Management Skill 📋

让 AI 助手帮你高效管理任务的标准化技能包！

## 触发条件

当用户提到以下关键词或意图时触发此技能：
- 创建任务："创建任务"、"新建任务"、"添加一个任务"、"帮我记个事"
- 查询任务："查看任务"、"查询任务"、"我的任务"、"有什么任务"
- 更新任务："更新任务"、"完成任务"、"取消任务"、"修改任务"
- 任务统计："任务统计"、"任务概况"、"有多少任务"
- 任务详情："任务详情"、"查看某个任务"

## 工具函数

### create_task
创建新任务

**参数**:
- `title` (string, 必需): 任务标题
- `description` (string, 可选): 任务描述
- `priority` (int, 可选): 优先级 1-4（1=Critical/2=High/3=Medium/4=Low），默认 3
- `project_name` (string, 可选): 所属项目名称
- `assignee_name` (string, 可选): 负责人姓名
- `agent_type` (string, 可选): Agent 类型（writer/reviewer/researcher）
- `agent_model` (string, 可选): Agent 模型名称
- `tags` (list, 可选): 标签列表
- `due_date` (string, 可选): 截止日期（YYYY-MM-DD 格式）

**示例**:
```python
create_task(title="审查文档", priority=2, description="检查完整性")
```

### query_tasks
查询任务列表

**参数**:
- `status` (string, 可选): 状态筛选（pending/agent_working/pending_review/reviewing/done/cancelled）
- `priority` (int, 可选): 优先级筛选
- `project_name` (string, 可选): 项目名称筛选
- `assignee_name` (string, 可选): 负责人筛选
- `agent_type` (string, 可选): Agent 类型筛选
- `keyword` (string, 可选): 关键词搜索
- `limit` (int, 可选): 返回数量限制，默认 20

**示例**:
```python
query_tasks(status="pending")
query_tasks(assignee_name="张三", limit=10)
```

### update_task_status
更新任务状态

**参数**:
- `task_uuid` (string, 必需): 任务 UUID
- `new_status` (string, 必需): 新状态（pending/agent_working/pending_review/reviewing/done/cancelled）

**示例**:
```python
update_task_status(task_uuid="abc123", new_status="done")
```

### get_task_detail
获取任务详情

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**示例**:
```python
get_task_detail(task_uuid="abc123")
```

## 执行规则

### 规则 1：创建任务流程
1. 当用户要创建任务时，必须确认标题
2. 如果用户没说优先级，默认设为 3（中等）
3. 如果用户没说截止日期，不设置截止日期
4. 创建成功后，显示任务 UUID 和基本状态

### 规则 2：查询任务流程
1. 默认显示最近 20 条任务
2. 如果结果超过 10 条，询问用户是否需要筛选
3. 支持按状态、负责人、优先级、项目筛选
4. 查询结果按创建时间倒序排列

### 规则 3：更新状态流程
1. 更新前确认任务存在
2. 验证状态流转的合法性
3. 如果是标记为"完成"，询问是否需要总结
4. 如果是"取消"，询问取消原因（可选）

### 规则 4：错误处理
1. 如果任务不存在，明确提示用户
2. 如果参数错误，给出修正建议
3. 如果数据库错误，重试一次后提示用户

## 使用示例

### 示例 1：创建任务
```
用户：帮我创建一个紧急任务，明天要完成项目评审
AI：好的，我来创建一个紧急任务（优先级 1）
    - 标题：项目评审
    - 截止时间：明天
    - 优先级：1（紧急）
    确认创建吗？
```

### 示例 2：查询任务
```
用户：查看我所有待处理的任务
AI：好的，正在查询您的待处理任务...
    找到 5 个待处理任务：
    1. [P2] 审查文档 - 待处理
    2. [P1] 项目评审 - 待处理
    3. [P3] 更新测试 - 待处理
    ...
```

### 示例 3：更新状态
```
用户：把"审查文档"标记为已完成
AI：正在更新任务状态...
    ✓ 任务"审查文档"已标记为完成
    需要我总结一下这个任务的完成情况吗？
```

## 安装说明

### 方式 1：使用 skills CLI（推荐）
```bash
npx skills add your-username/aitask-skill
```

### 方式 2：手动安装
```bash
git clone https://github.com/your-username/aitask-skill.git
cd aitask-skill
pip install -r backend/requirements.txt
```

### 方式 3：复制文件
```bash
# 复制到 Claude Skills 目录
cp -r . ~/.claude/skills/aitask-skill/

# 或复制到 Cursor Rules 目录
cp -r . ~/.cursor/rules/aitask-skill/
```

## 快速开始

### 1. 初始化数据库
```bash
cd backend
python start.py --init-db
```

### 2. 使用自然语言
直接对 AI 助手说：
```
"帮我创建一个紧急任务，明天要完成项目评审"
"查询我所有待处理的任务"
"把任务'审查文档'标记为已完成"
"统计一下本周完成了多少任务"
```

### 3. 使用 Python 调用
```python
from backend.app.skills.task_skill import TaskManagerSkill

# 初始化
skill = TaskManagerSkill()
await skill.initialize()

# 创建任务
task = await skill.create_task(
    title="审查第一卷第 10 章",
    description="检查时间线、人物、情节连贯性",
    priority=2,
    assignee_name="张三"
)

# 查询任务
tasks = await skill.query_tasks(status="pending")
print(f"待处理任务：{len(tasks['tasks'])} 个")

# 更新状态
await skill.update_task_status(task["uuid"], "done")
```

## 配置说明

### 环境变量
创建 `backend/.env` 文件：
```bash
# 数据库配置
DATABASE_URL=sqlite:///./data/tasks.db

# 服务器配置
SERVER_PORT=8000
SERVER_HOST=0.0.0.0
```

### 数据库位置
默认：`backend/data/tasks.db`

## 依赖

- Python >= 3.11
- FastAPI >= 0.109.0
- SQLAlchemy >= 2.0.25
- Pydantic >= 2.11.0
- aiosqlite >= 0.19.0

## 许可证

MIT License

## 支持

- GitHub Issues: https://github.com/your-username/aitask-skill/issues
- 邮箱：your.email@example.com
