# Task Management Skill

零依赖、高性能的任务管理技能，使用 Go + HTMX 实现。

## 平台支持

- **macOS** (Apple Silicon - ARM64) ✅
- **版本**: 0.2.1

## 快速开始

### 方式 1：作为技能使用（推荐）

**安装**：
```bash
# 从 GitHub Release 下载编译好的二进制
# https://github.com/xfwgithub/aitask-skill/releases

# 解压到技能目录
unzip task-skill-v0.2.1.zip -d ~/.agents/skills/task-management

# 或者从源码安装（开发用）
cp -r ~/github/aitask-skill/task-management ~/.agents/skills/task-management
```

**使用**：
- 在 Claude Code / Trae IDE 中直接使用
- 说"创建任务"、"查看任务"等自然语言即可

### 方式 2：独立运行 Web 服务

**从源码编译**：
```bash
cd task-management
go build -o task-skill .
./task-skill --server
```

**从 Release 下载**：
```bash
# 下载对应平台的二进制文件
./task-skill --server
```

**访问 Web 界面**：http://localhost:8080

### 方式 3：命令行使用

```bash
# 创建任务
echo '{"function": "create_task", "parameters": {"title": "我的任务"}}' | ./task-skill

# 查询任务
echo '{"function": "query_tasks"}' | ./task-skill

# 获取统计
echo '{"function": "get_dashboard_stats"}' | ./task-skill
```

## 功能

### Web 界面
- 📊 仪表盘 - 任务统计
- 📋 任务列表 - 浏览、筛选、搜索
- ➕ 创建任务 - 快速创建
- 📝 任务详情 - 查看和编辑
- ✅ 状态管理 - 更新任务状态

### API 端点
- `POST /api/tasks` - 创建任务
- `GET /api/tasks` - 查询任务
- `PUT /api/tasks/:uuid/status` - 更新状态
- `GET /api/stats` - 获取统计

## 技术栈

- **后端**: Go + Echo Framework
- **数据库**: SQLite (纯 Go 实现)
- **前端**: HTMX + 原生 JavaScript
- **样式**: 自定义 CSS

## 项目结构

```
.
├── README.md              # 本文件
├── SKILL.md               # 技能定义（Claude Skills 标准格式）
└── task-management/       # 技能实现
    ├── main.go           # 主程序（CLI + Web 服务器）
    ├── db.go             # SQLite 数据库操作
    ├── server.go         # Web 服务器和 API
    ├── utils.go          # 工具函数
    ├── go.mod            # Go 模块依赖
    ├── templates/        # HTML 模板
    │   ├── base.html
    │   ├── index.html
    │   ├── tasks.html
    │   └── task_detail.html
    └── static/           # 静态资源
        ├── css/
        │   └── style.css
        └── js/
            └── app.js
```

## 配置端口

**方式 1：环境变量**
```bash
export TASK_SKILL_PORT=3000
go run .
```

**方式 2：启动脚本**
```bash
cd task-management
chmod +x start.sh
TASK_SKILL_PORT=3000 ./start.sh
```

## 编译二进制

```bash
cd task-management
go build -o task-skill .
./task-skill --server
```

## 常见问题

**Q: 端口被占用怎么办？**
A: 使用环境变量设置其他端口：`TASK_SKILL_PORT=3000 go run .`

**Q: 数据保存在哪里？**
A: SQLite 数据库文件保存在 `task-management/` 目录下。

**Q: 如何停止服务？**

A: 根据运行方式选择：
- **前台运行**：在终端按 `Ctrl+C`
- **后台运行**：`kill $(pgrep -f "task-skill")` 或 `killall task-skill`
- **查看进程**：`ps aux | grep task-skill`

## 许可证

MIT License
