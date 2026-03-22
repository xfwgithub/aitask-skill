---
name: task-management
description: 零依赖、高性能的任务管理技能。当用户需要创建、查询、更新、删除任务或获取任务统计时使用此技能。
version: 0.2.12
platform: macOS (Apple Silicon)
---

## ⚠️ 重要：执行方式

**本技能通过 CLI 模式运行，Agent 通过标准输入/输出 (stdin/stdout) 与技能交互。**

### 如何调用技能

**重要**：本技能通过 **Bash 命令**调用，命令格式如下：

```bash
cd ~/.agents/skills/task-management && echo '{"function": "函数名", "parameters": {...}}' | ./task-skill
```

**注意**：
- 必须先 `cd` 到技能目录
- 通过 `echo` 将 JSON 参数传递给 `./task-skill`
- 技能通过 stdout 返回 JSON 结果

**示例**：
```bash
# 查询待办任务
cd ~/.agents/skills/task-management && echo '{"function": "query_tasks", "parameters": {"status": "pending"}}' | ./task-skill

# 创建任务
cd ~/.agents/skills/task-management && echo '{"function": "create_task", "parameters": {"title": "审查文档", "project": "aitask-skill"}}' | ./task-skill

# 领取任务
cd ~/.agents/skills/task-management && echo '{"function": "claim_task", "parameters": {"task_uuid": "abc-123"}}' | ./task-skill
```

### ⛔ 严格禁止的行为

**无论任何情况，都不得使用以下方式**：

**错误 1**：使用 curl 调用 API（即使已启动 Web 服务）
```bash
curl http://localhost:8080/api/tasks  # ❌ 绝对禁止！
```

**错误 2**：先启动服务器再通过 API 调用
```bash
./task-skill --server  # ❌ 不需要！
```

**错误 3**：启动 Web 服务后改用 HTTP API
```bash
# 即使 Web UI 已启动，也必须继续使用 CLI
echo '{"function": "query_tasks", ...}' | ./task-skill  # ✅ 正确
curl http://localhost:8080/api/tasks  # ❌ 错误！
```

### ℹ️ 关于 --server 模式

`./task-skill --server` 启动 **Web UI**，用于**人机协作**。

**Agent 调用原则**：
- ✅ Agent **始终**通过 CLI (stdin/stdout) 调用技能
- ✅ **无论 Web 服务是否启动**，Agent 都使用 CLI 方式
- ❌ Agent **不得使用** curl 或任何 HTTP 客户端调用 API

**人机协作场景**：
- 如果用户需要查看 Web 界面，Agent 可以启动 Web 服务
- Web UI 启动后，**Agent 继续通过 CLI 与技能交互**
- 人类通过浏览器访问 Web UI，Agent 通过 CLI 访问，两者数据互通

---

## 触发条件

当用户提到以下关键词或意图时触发此技能：

### 创建任务
- "创建任务"、"新建任务"、"添加一个任务"、"帮我记个事"
- "我要做一个..."、"记得做..."、"提醒我..."

### 查询任务
- "查看任务"、"查询任务"、"我的任务"、"有什么任务"
- "待处理的任务"、"未完成的任务"

### 更新任务
- "更新任务"、"提交初审"、"取消任务"、"修改任务"
- "完成初核"
- "标记为已完成"、"删除任务"

### 任务统计
- "任务统计"、"任务概况"、"有多少任务"
- "完成情况如何"

### 任务详情
- "任务详情"、"查看某个任务"、"任务信息"

### 回收任务
- "回收任务"、"任务到期"、"重置任务状态"

### 更新技能
- "更新技能"、"升级技能"、"检查更新"

**注意**: 技能更新应通过下载 GitHub Release 的最新版本完成，不应使用 git pull 源码。

## 工具函数

### create_task
创建新任务

**参数**:
- `title` (string, 必需): 任务标题
- `description` (string, 可选): 任务描述
- `priority` (int, 可选): 优先级 1-4（1=Critical/2=High/3=Medium/4=Low），默认 3
- `project` (string, 必需): 项目名称
- `tags` ([]string, 可选): 标签列表
- `assignee_name` (string, 可选): 负责人姓名
- `agent_type` (string, 可选): Agent 类型（writer/reviewer/researcher）
- `agent_model` (string, 可选): Agent 模型名称

**示例**:
```json
{
  "function": "create_task",
  "parameters": {
    "title": "审查文档",
    "project": "aitask-skill",
    "priority": 2,
    "description": "检查完整性"
  }
}
```

### query_tasks
查询任务列表

**参数**:
- `status` (string, 可选): 状态筛选（pending/agent_working/agent_review/human_review/done/cancelled）
- `priority` (int, 可选): 优先级筛选
- `project` (string, 可选): 项目筛选
- `assignee_name` (string, 可选): 负责人筛选
- `keyword` (string, 可选): 关键词搜索
- `limit` (int, 可选): 返回数量限制，默认 20

**示例**:
```json
{
  "function": "query_tasks",
  "parameters": {
    "status": "pending"
  }
}
```

### claim_task
领取任务（pending → agent_working）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**示例**:
```json
{
  "function": "claim_task",
  "parameters": {
    "task_uuid": "abc-123"
  }
}
```

### submit_initial_review
提交初审（agent_working → agent_review）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**示例**:
```json
{
  "function": "submit_initial_review",
  "parameters": {
    "task_uuid": "abc-123"
  }
}
```

### review_task
提交人工审核（agent_review → human_review）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**示例**:
```json
{
  "function": "review_task",
  "parameters": {
    "task_uuid": "abc-123"
  }
}
```

### approve_task
人工审核通过（human_review → done）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**示例**:
```json
{
  "function": "approve_task",
  "parameters": {
    "task_uuid": "abc-123"
  }
}
```

### cancel_task
取消任务（任意状态 → cancelled）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**示例**:
```json
{
  "function": "cancel_task",
  "parameters": {
    "task_uuid": "abc-123"
  }
}
```

**任务状态说明**：
- `pending` - 待办（等待 agent 领取）
- `agent_working` - Agent 工作中（agent 已领取，正在处理）
- `agent_review` - Agent 审核中（agent 提交初审后，等待 agent 自己审核确认）
- `human_review` - 人工审核中（agent 审核通过后，提交给人工审核）
- `done` - 完成（人工审核通过）
- `cancelled` - 已取消

### get_task_detail
获取任务详情

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**示例**:
```json
{
  "function": "get_task_detail",
  "parameters": {
    "task_uuid": "abc-123"
  }
}
```

### get_version
获取版本号

**参数**: 无

**示例**:
```json
{
  "function": "get_version"
}
```

### get_dashboard_stats
获取仪表盘统计信息

**参数**: 无

**示例**:
```json
{
  "function": "get_dashboard_stats"
}
```

### recycle_tasks
回收到期未完成的任务

**参数**:
- `due_date` (string, 可选): 截止时间，回收此日期前未完成的任务

**示例**:
```json
{
  "function": "recycle_tasks",
  "parameters": {
    "due_date": "2026-03-22"
  }
}

## 安装路径

技能应安装在以下位置：
- **Agent 技能目录**: `~/.agents/skills/task-management/`
- **Claude Code 软链接**: `~/.claude/skills/task-management` → `~/.agents/skills/task-management`

## 初始化

首次使用或数据库不存在时：

```bash
# 1. 确保在正确的目录
cd ~/.agents/skills/task-management

# 2. 验证安装
./task-skill --version

# 3. 测试技能（会自动创建数据库）
echo '{"function": "get_version"}' | ./task-skill
```

## Web UI（可选）

如果需要 Web 界面管理任务，可以单独启动服务器：

```bash
# 启动 Web 服务器（仅用于 Web UI，非必需）
./task-skill --server

# 访问 http://localhost:8080
```

**注意**：Web UI 是可选的，Agent 调用技能不需要启动服务器。

## 更新技能说明

当用户要求更新技能时：

1. **不要**使用 `git pull` 拉取源码
2. **不要**使用 `git clone` 重新克隆
3. **正确方式**：
   - 访问 https://github.com/xfwgithub/aitask-skill/releases
   - 下载最新版本的 `task-skill-vX.Y.Z.zip`
   - 解压并替换 `~/.agents/skills/task-management/` 中的文件
   - 重启服务（如果使用 `--server` 模式）

4. 更新前可以使用 `get_version` 检查当前版本
```
