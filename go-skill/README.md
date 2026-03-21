# Task Management Skill - Go 版本 🚀

零依赖、高性能的任务管理技能实现

## ✨ 特性

- ✅ **零依赖** - 编译为单一二进制文件，无需运行时环境
- ✅ **高性能** - Go 语言原生并发支持
- ✅ **跨平台** - Windows、macOS、Linux 全支持
- ✅ **易部署** - 复制即可使用
- ✅ **兼容 Python 版本** - 完全相同的 API 接口

## 📦 安装

### 方式 1：使用预编译二进制（推荐）

```bash
# 下载对应平台的二进制文件
# macOS (Intel)
curl -L https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill-darwin-amd64 -o task-skill

# macOS (Apple Silicon)
curl -L https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill-darwin-arm64 -o task-skill

# Linux
curl -L https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill-linux-amd64 -o task-skill

# Windows
curl -L https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill-windows-amd64.exe -o task-skill.exe

# 添加执行权限
chmod +x task-skill
```

### 方式 2：从源码编译

```bash
# 需要安装 Go 1.21+
cd go-skill
go build -o task-skill
```

## 🚀 使用

### 作为 AI Skill 使用

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

### 编程调用

```go
package main

import (
    "encoding/json"
    "fmt"
    "task-skill"
)

func main() {
    skill := taskskill.NewSkill()
    
    // 创建任务
    result := skill.CreateTask(taskskill.CreateTaskInput{
        Title:    "审查文档",
        Priority: 2,
        Assignee: "张三",
    })
    
    jsonBytes, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(jsonBytes))
    
    // 查询任务
    tasks := skill.QueryTasks(taskskill.QueryTasksInput{
        Status: "pending",
    })
    
    // 更新状态
    skill.UpdateTaskStatus(taskskill.UpdateTaskStatusInput{
        TaskUUID:  "xxx",
        NewStatus: "done",
    })
    
    // 获取统计
    stats := skill.GetDashboardStats()
}
```

## 📖 API 文档

### create_task
创建新任务

**参数**:
- `title` (string, 必需): 任务标题
- `description` (string, 可选): 任务描述
- `priority` (int, 可选): 优先级 1-4，默认 3
- `tags` ([]string, 可选): 标签列表
- `assignee_name` (string, 可选): 负责人
- `agent_type` (string, 可选): Agent 类型
- `agent_model` (string, 可选): Agent 模型

### query_tasks
查询任务列表

**参数**:
- `status` (string, 可选): 状态筛选
- `priority` (int, 可选): 优先级筛选
- `assignee_name` (string, 可选): 负责人筛选
- `keyword` (string, 可选): 关键词搜索
- `limit` (int, 可选): 返回数量限制

### update_task_status
更新任务状态

**参数**:
- `task_uuid` (string, 必需): 任务 UUID
- `new_status` (string, 必需): 新状态

### get_task_detail
获取任务详情

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

### get_dashboard_stats
获取统计信息

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

## 🔧 开发

### 编译不同平台

```bash
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
go test -v ./...
```

## 📊 对比 Python 版本

| 特性 | Go 版本 | Python 版本 |
|------|---------|-------------|
| 依赖 | 无 | Python 3.11+ |
| 启动速度 | < 10ms | ~100ms |
| 内存占用 | ~5MB | ~50MB |
| 部署难度 | 复制文件 | 安装依赖 |
| 性能 | 高 | 中等 |

## 📄 许可证

MIT License
