# Task Management Skill

零依赖、高性能的任务管理技能，使用 Go + HTMX 实现。支持 pip 安装，提供直观的命令行界面。

## 平台支持

- **macOS** (Apple Silicon - ARM64) ✅

## 安装

### 方式 1：通过 pip 安装（推荐）

```bash
# 从 GitHub 安装最新版本
pip install git+https://github.com/xfwgithub/aitask-skill.git

# 验证安装
task-skill --version
task-skill --help
```

### 方式 2：下载 Release

```bash
# 下载最新版本
curl -L -o /tmp/task-skill-v0.4.0.zip https://github.com/xfwgithub/aitask-skill/releases/download/v0.4.0/task-skill-v0.4.0.zip

# 解压
unzip /tmp/task-skill-v0.4.0.zip -d /tmp/

# 移动到技能目录
mkdir -p ~/.agents/skills
mv /tmp/task-skill-v0.4.0 ~/.agents/skills/task-management

# 添加到 PATH
export PATH=$PATH:~/.agents/skills/task-management

# 验证
~/.agents/skills/task-management/task-skill --version
```

### 方式 3：从源码安装

```bash
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill
pip install -e .
```

## 更新

### pip 安装方式更新

```bash
# 更新到最新版本
pip install --upgrade git+https://github.com/xfwgithub/aitask-skill.git

# 验证更新
task-skill --version
```

### Release 下载方式更新

```bash
# 删除旧版本
rm -rf ~/.agents/skills/task-management

# 重新下载安装（参考上方"方式 2"）
```

## 使用方式

### 命令行界面（推荐）

```bash
# 查看帮助
task-skill --help

# 创建任务
task-skill create-task --title "我的任务" --project "myproject"
task-skill create-task --title "紧急任务" --project "demo" --priority 1

# 创建子任务
task-skill create-task --title "子任务" --project "demo" --parent <父任务UUID>

# 列出任务
task-skill list-tasks
task-skill ls --status pending
task-skill ls --project demo --limit 10

# 获取任务详情
task-skill get-task <uuid>

# 任务状态操作
task-skill claim-task <uuid>          # 领取任务
task-skill submit-review <uuid>       # 提交初审
task-skill review-task <uuid>         # 提交人工审核
task-skill approve-task <uuid>        # 审核通过
task-skill cancel-task <uuid>         # 取消任务
task-skill delete-task <uuid>         # 物理删除（彻底删除）

# 统计信息
task-skill stats

# 回收任务
task-skill recycle-tasks --due-date 2026-03-22
```

### Web 服务模式（可选）

```bash
task-skill --server
```

访问 http://localhost:8080

**说明**：
- ℹ️ Web UI 仅供人类用户通过浏览器管理任务
- ℹ️ Agent 调用技能**不需要**启动 Web 服务
- ℹ️ 如需要 Web UI，可在后台运行：`task-skill --server &`

## 功能特性

### 任务管理
- ✅ 创建任务（支持优先级、项目、标签、负责人）
- ✅ 创建子任务（支持多级任务结构）
- ✅ 查询任务（支持状态、项目、关键词筛选）
- ✅ 更新任务状态
- ✅ 任务详情查看
- ✅ 物理删除任务（彻底删除）

### 主子任务支持
```bash
# 创建父任务
task-skill create-task --title "项目规划" --project "demo"
# 返回: uuid: abc-123

# 创建子任务
task-skill create-task --title "需求分析" --project "demo" --parent abc-123
task-skill create-task --title "技术设计" --project "demo" --parent abc-123

# 查看子任务详情（包含 parent_uuid 字段）
task-skill get-task <子任务uuid>
```

### 状态流转
```
pending → agent_working → agent_review → human_review → done
    ↓           ↓             ↓              ↓
    └───────────┴─────────────┴──────────────┴──→ cancelled
```

### 命令列表
| 命令 | 说明 |
|------|------|
| `create-task` | 创建新任务 |
| `list-tasks`, `ls` | 列出任务 |
| `get-task` | 获取任务详情 |
| `claim-task` | 领取任务 |
| `submit-review` | 提交初审 |
| `review-task` | 提交人工审核 |
| `approve-task` | 审核通过 |
| `cancel-task` | 取消任务 |
| `delete-task` | 物理删除任务 |
| `recycle-tasks` | 回收到期任务 |
| `stats` | 统计信息 |

## 技术栈

- **后端**: Go + Echo Framework
- **数据库**: SQLite
- **前端**: HTMX + 原生 JavaScript
- **CLI**: 子命令风格命令行界面

## 版本历史

- **v0.4.0** - 添加子任务支持，改进 CLI 界面，支持 pip 安装
- **v0.3.1** - 添加物理删除功能
- **v0.3.0** - 初始版本
