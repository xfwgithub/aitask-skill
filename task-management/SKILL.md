---
name: task-management
description: 零依赖、高性能的任务管理技能。当用户需要创建、查询、更新、删除任务或获取任务统计时使用此技能。
version: 0.2.3
platform: macOS (Apple Silicon)
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
- "更新任务"、"完成任务"、"取消任务"、"修改任务"
- "标记为已完成"、"删除任务"

### 任务统计
- "任务统计"、"任务概况"、"有多少任务"
- "完成情况如何"

### 任务详情
- "任务详情"、"查看某个任务"、"任务信息"

### 回收任务
- "回收任务"、"任务到期"、"重置任务状态"

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

### update_task_status
更新任务状态

**参数**:
- `task_uuid` (string, 必需): 任务 UUID
- `new_status` (string, 必需): 新状态

**任务状态说明**：
- `pending` - 待办（等待 agent 或人处理）
- `agent_working` - Agent 工作中（agent 正在处理）
- `agent_review` - Agent 审核中（agent 完成，等待审核）
- `human_review` - 人工审核中（等待人审核）
- `done` - 完成（审核通过）
- `cancelled` - 取消

**状态流转**：
```
pending → agent_working → agent_review → human_review → done
    ↓           ↓             ↓              ↓
    └───────────┴─────────────┴──────────────┘
                        ↓
                    cancelled
```

- 审核不通过的任务会重新回到 `pending` 状态

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
```
