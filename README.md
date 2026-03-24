# Task Management Skill

零依赖、高性能的任务管理技能，使用 Go + HTMX 实现。单文件绿色免安装，提供直观的命令行界面和 Web UI。

## 平台支持

- **macOS** (Apple Silicon - ARM64) ✅

## 安装

### 方式 1：一键安装（推荐）

直接下载预编译二进制文件到系统 PATH 中：

```bash
sudo curl -sL https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill -o /usr/local/bin/task-skill
sudo chmod +x /usr/local/bin/task-skill

# 验证安装
task-skill --version
```

**说明**：
- ✅ 单文件，自带 Web UI 静态资源，无需依赖任何环境
- ✅ 极速运行，没有 Python 包装器的冷启动延迟

### 方式 2：作为 AI Agent 技能配置

1. 确保已按上方步骤安装 `task-skill` 并可通过终端运行。
2. 下载 [SKILL.md](https://raw.githubusercontent.com/xfwgithub/aitask-skill/main/skills/task-management/SKILL.md) 并放入你的技能配置目录中：

```bash
mkdir -p ~/.agents/skills/task-management
curl -sL https://raw.githubusercontent.com/xfwgithub/aitask-skill/main/skills/task-management/SKILL.md -o ~/.agents/skills/task-management/SKILL.md
```

## 更新

```bash
sudo curl -sL https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill -o /usr/local/bin/task-skill
sudo chmod +x /usr/local/bin/task-skill
task-skill --version
```

## 使用方式

### 命令行界面（直接使用）

```bash
# 查看帮助
task-skill --help

# 创建任务
task-skill create-task --title "我的任务" --project "myproject"
task-skill create-task --title "紧急任务" --project "demo" --priority 1

# 创建子任务
task-skill create-task --title "子任务" --project "demo" --parent <父任务 UUID>

# 列出任务
task-skill list-tasks
task-skill ls --status pending
task-skill ls --project demo --limit 10

# 获取任务详情
task-skill get-task <uuid>

# 任务状态操作
task-skill claim-task <uuid> [意见]         # 领取任务
task-skill submit-review <uuid> [意见]      # 提交初审
task-skill review-task <uuid> [意见]        # 提交人工审核
task-skill approve-task <uuid> [意见]       # 审核通过
task-skill reject-task <uuid> [意见]        # 审核不通过并退回
task-skill cancel-task <uuid> [意见]        # 取消任务
task-skill delete-task <uuid>              # 物理删除（彻底删除）

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
# 返回：uuid: abc-123

# 创建子任务
task-skill create-task --title "需求分析" --project "demo" --parent abc-123
task-skill create-task --title "技术设计" --project "demo" --parent abc-123

# 查看子任务详情（包含 parent_uuid 字段）
task-skill get-task <子任务 uuid>
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
| `reject-task` | 审核不通过并退回 |
| `cancel-task` | 取消任务 |
| `delete-task` | 物理删除任务 |
| `recycle-tasks` | 回收到期任务 |
| `stats` | 统计信息 |

## 技术栈

- **后端**: Go + Echo Framework
- **数据库**: SQLite
- **前端**: HTMX + 原生 JavaScript
- **CLI**: 子命令风格命令行界面
- **安装**: Python pip 包装器，自动下载预编译二进制

## 版本历史

- **v1.2.0** - 重大架构升级：移除 Python/pip 依赖，采用 Go `go:embed` 静态编译，实现真正的单文件免安装运行。
- **v1.1.6** - 修复了 CLI 中的版本强校验导致的循环下载 Bug 以及创建任务页面标签字段的 Bug
- **v1.1.0** - 增加任务状态变更时的活动记录和说明评论功能，强化前端操作体验和历史追溯能力
- **v1.0.0** - 第一个正式稳定版本
- **v0.4.3** - 修复由于运行模式导致数据库文件割裂（Split Brain）的问题，统一绝对路径
- **v0.4.2** - 统一版本号管理，修复版本号硬编码导致的多个文件版本不一致问题
- **v0.4.1** - 修复 Web UI 静态资源路径问题，自动下载完整包
- **v0.4.0** - 添加子任务支持，改进 CLI 界面，支持 pip 安装，自动下载二进制
- **v0.3.1** - 添加物理删除功能
- **v0.3.0** - 初始版本
