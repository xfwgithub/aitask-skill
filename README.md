# Task Management Skill 📋

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go 1.21+](https://img.shields.io/badge/go-1.21+-blue.svg)](https://go.dev/)

零依赖、高性能的任务管理技能！

**核心特性**：
- ✅ **零依赖** - 无需 Python/Node.js，单一二进制文件
- ✅ **AI Skill 支持** - 自然语言操作任务
- ✅ **快速创建任务** - 支持自然语言和结构化创建
- ✅ **智能查询任务** - 多维度筛选和搜索
- ✅ **任务统计** - 实时统计和数据分析
- ✅ **跨平台** - Windows、macOS、Linux 全支持

## 📦 安装

### 方式 1：下载预编译二进制（推荐）

```bash
# macOS (Apple Silicon)
curl -L https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill-darwin-arm64 -o task-skill
chmod +x task-skill

# macOS (Intel)
curl -L https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill-darwin-amd64 -o task-skill
chmod +x task-skill

# Linux
curl -L https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill-linux-amd64 -o task-skill
chmod +x task-skill

# Windows (PowerShell)
curl.exe -L https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill-windows-amd64.exe -o task-skill.exe
```

### 方式 2：从源码编译

```bash
# 需要安装 Go 1.21+
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill/go-skill
go build -o task-skill
```

### 方式 3：使用 skills CLI

```bash
npx skills add xfwgithub/aitask-skill
```

## 🚀 快速开始

### 使用自然语言

直接对 AI 助手说：

```
"帮我创建一个紧急任务，明天要完成项目评审"
"查询我所有待处理的任务"
"把任务'审查文档'标记为已完成"
"统计一下本周完成了多少任务"
```

### 命令行调用

```bash
# 创建任务
echo '{"function": "create_task", "parameters": {"title": "审查文档", "priority": 2}}' | ./task-skill

# 查询任务
echo '{"function": "query_tasks", "parameters": {"status": "pending"}}' | ./task-skill

# 更新状态
echo '{"function": "update_task_status", "parameters": {"task_uuid": "xxx", "new_status": "done"}}' | ./task-skill

# 获取统计
echo '{"function": "get_dashboard_stats"}' | ./task-skill
```

## 🎯 技能触发条件

当用户提到以下关键词或意图时触发此技能：

- **创建任务**："创建任务"、"新建任务"、"添加一个任务"、"帮我记个事"
- **查询任务**："查看任务"、"查询任务"、"我的任务"、"有什么任务"
- **更新任务**："更新任务"、"完成任务"、"取消任务"、"修改任务"
- **任务统计**："任务统计"、"任务概况"、"有多少任务"
- **任务详情**："任务详情"、"查看某个任务"

## 📖 API 文档

### create_task
创建新任务

**参数**:
- `title` (string, 必需): 任务标题
- `description` (string, 可选): 任务描述
- `priority` (int, 可选): 优先级 1-4（1=Critical/2=High/3=Medium/4=Low），默认 3
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
    "description": "检查完整性",
    "priority": 2,
    "tags": ["工作", "紧急"],
    "assignee_name": "张三"
  }
}
```

### query_tasks
查询任务列表

**参数**:
- `status` (string, 可选): 状态筛选（pending/agent_working/done/cancelled）
- `priority` (int, 可选): 优先级筛选
- `assignee_name` (string, 可选): 负责人筛选
- `keyword` (string, 可选): 关键词搜索
- `limit` (int, 可选): 返回数量限制，默认 20

**示例**:
```json
{
  "function": "query_tasks",
  "parameters": {
    "status": "pending",
    "limit": 10
  }
}
```

### update_task_status
更新任务状态

**参数**:
- `task_uuid` (string, 必需): 任务 UUID
- `new_status` (string, 必需): 新状态（pending/agent_working/done/cancelled）

**示例**:
```json
{
  "function": "update_task_status",
  "parameters": {
    "task_uuid": "abc-123",
    "new_status": "done"
  }
}
```

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

### get_dashboard_stats
获取统计信息

**示例**:
```json
{
  "function": "get_dashboard_stats"
}
```

**返回**:
```json
{
  "total": 100,
  "pending": 20,
  "agent_working": 5,
  "done": 70,
  "cancelled": 5
}
```

## 🏗️ 项目结构

```
aitask-skill/
├── SKILL.md              # 技能定义
├── skill.json            # 技能元数据
├── README.md             # 使用文档
├── .gitignore            # Git 忽略
└── go-skill/             # Go 实现
    ├── main.go          # 主程序
    ├── utils.go         # 工具函数
    ├── go.mod           # Go 模块
    └── README.md        # Go 版本文档
```

## 🔧 开发

### 编译不同平台

```bash
cd go-skill

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o task-skill-darwin-amd64

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o task-skill-darwin-arm64

# Linux
GOOS=linux GOARCH=amd64 go build -o task-skill-linux-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o task-skill-windows-amd64.exe
```

### 运行测试

```bash
cd go-skill
go test -v ./...
```

## � 性能对比

| 特性 | Go 版本 | Python 版本 |
|------|---------|-------------|
| 依赖 | 无 | Python 3.11+ |
| 启动速度 | < 10ms | ~100ms |
| 内存占用 | ~5MB | ~50MB |
| 部署难度 | 复制文件 | 安装依赖 |
| 性能 | 高 | 中等 |

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
