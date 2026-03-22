# Task Management Skill

零依赖、高性能的任务管理技能，使用 Go + HTMX 实现。

## 平台支持

- **macOS** (Apple Silicon - ARM64) ✅

## 安装

### 方式 1：下载 Release（推荐）

```bash
# 下载最新版本（以 v0.2.8 为例）
curl -L -o /tmp/task-skill-v0.2.8.zip https://github.com/xfwgithub/aitask-skill/releases/download/v0.2.8/task-skill-v0.2.8.zip

# 解压
unzip /tmp/task-skill-v0.2.8.zip -d /tmp/

# 移动到技能目录
mkdir -p ~/.agents/skills
mv /tmp/task-skill-v0.2.8 ~/.agents/skills/task-management

# 清理
rm /tmp/task-skill-v0.2.8.zip
```

### 方式 2：手动下载

访问 [GitHub Releases](https://github.com/xfwgithub/aitask-skill/releases) 下载 `task-skill-vX.Y.Z.zip`，解压到 `~/.agents/skills/task-management`。

## 使用方式

### 方式 1：CLI 模式（推荐）

```bash
# 创建任务
echo '{"function": "create_task", "parameters": {"title": "我的任务", "project": "myproject"}}' | ./task-skill

# 查询任务
echo '{"function": "query_tasks", "parameters": {"status": "pending"}}' | ./task-skill

# 获取统计
echo '{"function": "get_dashboard_stats"}' | ./task-skill

# 查询版本
echo '{"function": "get_version"}' | ./task-skill
```

### 方式 2：Web 服务模式

```bash
./task-skill --server
```

访问 http://localhost:8080

## 功能特性

### 任务管理
- 创建任务（支持优先级、项目、标签）
- 查询任务（支持状态、项目、关键词筛选）
- 更新任务状态
- 任务详情查看

### 状态流转
```
pending → agent_working → agent_review → human_review → done
    ↓           ↓             ↓              ↓
    └───────────┴─────────────┴──────────────┴──→ cancelled
```

### 原子化操作
- `claim_task` - 领取任务
- `submit_initial_review` - 提交初审
- `review_task` - 提交人工审核
- `approve_task` - 审核通过
- `cancel_task` - 取消任务

## 技术栈

- **后端**: Go + Echo Framework
- **数据库**: SQLite
- **前端**: HTMX + 原生 JavaScript
