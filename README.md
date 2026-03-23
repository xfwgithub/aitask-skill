# Task Management Skill

零依赖、高性能的任务管理技能，使用 Go + HTMX 实现。支持 pip 安装，提供直观的命令行界面。

## 平台支持

- **macOS** (Apple Silicon - ARM64) ✅

## 快速安装 (npx skills)

如果你的系统已配置 `npx skills` 环境，你可以直接使用以下命令安装此技能：

```bash
npx skills add xfwgithub/aitask-skill@task-management
```

## 安装

### 方式 1：通过 pip 安装（推荐）

```bash
# 从 GitHub 安装最新版本
pip install git+https://github.com/xfwgithub/aitask-skill.git

# 验证安装（首次运行时会自动下载二进制文件）
task-skill --version
task-skill --help
```

**说明**：
- ✅ 首次运行时会自动从 GitHub Releases 下载预编译的完整包（约 11MB，包含二进制 + 静态资源）
- ✅ 后续运行直接使用缓存，无需重复下载
- ✅ 无需 Go 环境，无需编译

### 方式 2：下载 Release

```bash
# 下载最新版本（推荐）
wget https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill.zip
unzip task-skill.zip
cd task-skill
./task-skill --version

# 或者下载指定版本
# wget https://github.com/xfwgithub/aitask-skill/releases/download/v0.4.3/task-skill-v0.4.3.zip
```

**注意**：此方式下载的是完整包，包含二进制和静态资源，适合直接运行或调试。
如要作为 AI Agent 技能使用，请参考上方"方式 1"安装。

### 方式 3：从源码安装

```bash
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill
pip install -e .
```

**说明**：
- ℹ️ 适合开发和调试
- ℹ️ 需要先编译 Go 二进制（运行 `cd skills/task-management && bash build.sh`）

## 更新

### pip 安装方式更新

```bash
# 更新到最新版本
pip install --upgrade git+https://github.com/xfwgithub/aitask-skill.git

# 验证更新
task-skill --version
```

**说明**：
- ✅ 如果版本号更新，首次运行时会自动下载新版本的二进制文件

### Release 下载方式更新

```bash
# 删除旧版本（替换为你的技能目录）
rm -rf <你的技能目录>/skills/task-management

# 重新下载安装（参考上方"方式 2"）
```

## 使用方式

### 作为 AI Agent 技能（推荐）

**安装步骤**：

1. **安装 Python 包**（提供命令行工具）
   ```bash
   pip install git+https://github.com/xfwgithub/aitask-skill.git
   ```

2. **复制 SKILL.md 到技能目录**
   ```bash
   # 找到项目中的 SKILL.md 文件
   # 复制到你的技能目录（具体目录取决于你的 IDE/Agent 配置）
   cp /path/to/aitask-skill/skills/task-management/SKILL.md <你的技能目录>/task-management/
   ```

3. **首次使用**
   - 当 AI Agent 需要执行任务管理功能时，会自动调用 `task-skill` 命令
   - 首次运行时会自动下载二进制文件和静态资源（约 11MB）

**说明**：
- ✅ 只需要复制 `SKILL.md` 文件到技能目录
- ✅ 不需要复制其他文件（二进制和静态资源会自动下载）
- ✅ `task-skill` 命令会自动添加到 PATH，AI Agent 可以直接调用
- ℹ️ 技能目录位置取决于你使用的 IDE/Agent 配置

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

- **v0.4.3** - 修复由于运行模式导致数据库文件割裂（Split Brain）的问题，统一绝对路径
- **v0.4.2** - 统一版本号管理，修复版本号硬编码导致的多个文件版本不一致问题
- **v0.4.1** - 修复 Web UI 静态资源路径问题，自动下载完整包
- **v0.4.0** - 添加子任务支持，改进 CLI 界面，支持 pip 安装，自动下载二进制
- **v0.3.1** - 添加物理删除功能
- **v0.3.0** - 初始版本
