# Task Management Skill

零依赖、高性能的任务管理技能，使用 Go + HTMX 实现。

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

## 快速开始

### 启动 Web 服务器

```bash
cd task-management

# 安装依赖
go mod download

# 启动服务器（默认端口 8080）
go run .

# 或使用自定义端口
TASK_SKILL_PORT=3000 go run .

# 打开浏览器访问
# http://localhost:8080
```

### 配置端口

**方式 1：环境变量**
```bash
export TASK_SKILL_PORT=3000
go run .
```

**方式 2：启动脚本**
```bash
TASK_SKILL_PORT=3000 ./start.sh
```

### 命令行使用

```bash
cd task-management

# 创建任务
echo '{"function": "create_task", "parameters": {"title": "我的任务"}}' | go run .

# 查询任务
echo '{"function": "query_tasks"}' | go run .
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

## IDE 集成

### VS Code

**方式 1：使用 Web 界面（推荐）**
```bash
cd task-management
go run .
# 浏览器访问 http://localhost:8080
```

**方式 2：集成终端**
在 VS Code 中打开集成终端：
```bash
cd task-management
./start.sh
```

**方式 3：使用 Code Runner 插件**
安装 Code Runner 插件，然后按 `Ctrl+Alt+N` 运行

### JetBrains IDEs (GoLand, IntelliJ IDEA)

**方式 1：Run Configuration**
1. `Run` → `Edit Configurations...`
2. 点击 `+` → `Go Build`
3. 设置：
   - Working directory: `task-management`
   - Program arguments: `--server`
   - Environment variables: `TASK_SKILL_PORT=8080`
4. 点击 `Run`

**方式 2：Go Run**
右键点击 `main.go` → `Run 'go run main.go --server'`

**方式 3：Terminal**
在 IDE 终端中：
```bash
cd task-management
go run . --server
```

### Cursor

**方式 1：使用内置终端**
```bash
cd task-management
./start.sh
```

**方式 2：配置 Tasks**
创建 `.cursor/tasks.json`：
```json
{
  "tasks": {
    "start-skill": {
      "command": "cd task-management && ./start.sh",
      "name": "Start Task Skill"
    }
  }
}
```

### Trae

**方式 1：使用终端**
```bash
cd task-management
./start.sh
```

**方式 2：配置运行脚本**
在项目根目录创建 `.trae/tasks.json`：
```json
{
  "tasks": {
    "start-skill": {
      "command": "cd task-management && ./start.sh",
      "name": "启动任务管理技能"
    }
  }
}
```

### 通用方式（任何 IDE）

**1. 系统终端**
```bash
cd /path/to/aitask-skill/task-management
./start.sh
```

**2. 直接运行编译后的二进制**
```bash
cd task-management
go build -o task-skill .
./task-skill --server
```

**3. 使用 Docker（可选）**
创建 `Dockerfile`：
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY task-management/ .
RUN go mod download
EXPOSE 8080
CMD ["./task-skill", "--server"]
```

## 初始化方式

### 服务型 Skills 的初始化模式

本项目属于**服务型 Skill**，与官方文档型 Skills 不同：

**文档型 Skills（官方）**：
- 形式：`SKILL.md` 文档
- 作用：指导 Claude 如何完成任务
- 执行：Claude 阅读文档后执行
- 示例：pdf, xlsx, pptx 等

**服务型 Skills（本项目）**：
- 形式：独立运行的服务
- 作用：提供实际的功能接口
- 执行：Claude 调用服务 API
- 示例：mcp-builder, webapp-testing 等

### 初始化步骤

#### 方式 1：Agent 自动初始化（推荐）

让 AI 助手自动完成所有配置：

```bash
# 克隆项目
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill/task-management

# 运行快速初始化脚本
chmod +x quick-init.sh
./quick-init.sh
```

脚本会自动：
- ✅ 检查 Go 环境
- ✅ 安装依赖
- ✅ 编译程序
- ✅ 验证功能
- ✅ 生成配置文件

**使用配置文件初始化：**

AI 助手可以读取 `skill.config.json` 自动执行：

```json
{
  "init": {
    "script": "./quick-init.sh",
    "commands": [
      "go mod download",
      "go build -o task-skill ."
    ],
    "verify": {
      "command": "echo '{\"function\": \"get_dashboard_stats\"}' | ./task-skill",
      "expect": "total"
    }
  }
}
```

#### 方式 2：手动初始化

```bash
# 克隆项目
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill/task-management

# 安装依赖
go mod download

# 编译
go build -o task-skill .

# 启动服务
./task-skill --server
```

#### 方式 3：使用启动脚本

```bash
cd task-management
chmod +x start.sh
./start.sh
```

#### 2. 配置 Claude

在 Claude 的 Skills 配置中添加：
```yaml
name: task-management
description: 任务管理技能，用于创建、查询、更新任务
type: service
endpoint: http://localhost:8080
```

#### 3. 验证服务

```bash
# 测试 API
curl http://localhost:8080/api/stats

# 预期输出
{
  "total": 0,
  "pending": 0,
  "agent_working": 0,
  "done": 0,
  "cancelled": 0
}
```

#### 4. 持久化运行

**使用 systemd（Linux）：**
```bash
sudo systemctl edit --force task-skill
```

内容：
```ini
[Unit]
Description=Task Management Skill
After=network.target

[Service]
Type=simple
User=youruser
WorkingDirectory=/path/to/aitask-skill/task-management
Environment="TASK_SKILL_PORT=8080"
ExecStart=/path/to/task-skill --server
Restart=always

[Install]
WantedBy=multi-user.target
```

启动：
```bash
sudo systemctl enable task-skill
sudo systemctl start task-skill
```

**使用 pm2（跨平台）：**
```bash
npm install -g pm2
cd task-management
pm2 start task-skill --name task-skill -- --server
pm2 save
```

## 许可证

MIT License
