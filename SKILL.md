---
name: "task-management"
description: "AI 任务管理技能。当用户需要创建、查询、更新、删除任务或获取任务统计时触发此技能。"
version: "1.0.0"
author: "xfwgithub"
license: "MIT"
keywords:
  - task-management
  - productivity
  - ai-skill
  - agent
  - go
engines:
  go: ">=1.21"
---

# Task Management Skill 📋

零依赖、高性能的任务管理技能！

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
    "priority": 2,
    "description": "检查完整性"
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
    "status": "pending"
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
获取仪表盘统计

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

## 执行规则

### 规则 1：创建任务流程
1. 当用户要创建任务时，必须确认标题
2. 如果用户没说优先级，默认设为 3（中等）
3. 创建成功后，显示任务 UUID 和基本状态

### 规则 2：查询任务流程
1. 默认显示最近 20 条任务
2. 如果结果超过 10 条，询问用户是否需要筛选
3. 支持按状态、负责人、优先级筛选

### 规则 3：更新状态流程
1. 更新前确认任务存在
2. 验证状态流转的合法性
3. 如果是标记为"完成"，询问是否需要总结

### 规则 4：错误处理
1. 如果任务不存在，明确提示用户
2. 如果参数错误，给出修正建议

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

### 方式 1：下载预编译二进制（推荐）
```bash
# macOS (Apple Silicon)
curl -L https://github.com/xfwgithub/aitask-skill/releases/latest/download/task-skill-darwin-arm64 -o task-skill
chmod +x task-skill
```

### 方式 2：从源码编译
```bash
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill/go-skill
go build -o task-skill
```

### 方式 3：使用 skills CLI
```bash
npx skills add xfwgithub/aitask-skill
```

## 许可证

MIT License

## 支持

- GitHub Issues: https://github.com/xfwgithub/aitask-skill/issues
