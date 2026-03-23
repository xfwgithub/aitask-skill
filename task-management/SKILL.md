---
name: task-management
description: 零依赖、高性能的任务管理技能。当用户需要创建、查询、更新、删除任务或获取任务统计时使用此技能。
version: 0.3.1
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
- "更新任务"、"提交初审"、"取消任务"、"修改任务"
- "完成初核"
- "标记为已完成"、"删除任务"

### 任务统计
- "任务统计"、"任务概况"、"有多少任务"
- "完成情况如何"

### 任务详情
- "任务详情"、"查看某个任务"、"任务信息"

### 回收任务
- "回收任务"、"任务到期"、"重置任务状态"

### 删除任务
- "彻底删除"、"物理删除"、"删除任务"

### 更新技能
- "更新技能"、"升级技能"、"检查更新"


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

**调用示例**:
```bash
echo '{"function": "create_task", "parameters": {"title": "审查文档", "project": "aitask-skill", "priority": 2, "description": "检查完整性"}}' | ./task-skill
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

**调用示例**:
```bash
echo '{"function": "query_tasks", "parameters": {"status": "pending"}}' | ./task-skill
```

### claim_task
领取任务（pending → agent_working）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**调用示例**:
```bash
echo '{"function": "claim_task", "parameters": {"task_uuid": "abc-123"}}' | ./task-skill
```

### submit_initial_review
提交初审（agent_working → agent_review）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**调用示例**:
```bash
echo '{"function": "submit_initial_review", "parameters": {"task_uuid": "abc-123"}}' | ./task-skill
```

### review_task
提交人工审核（agent_review → human_review）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**调用示例**:
```bash
echo '{"function": "review_task", "parameters": {"task_uuid": "abc-123"}}' | ./task-skill
```

### approve_task
人工审核通过（human_review → done）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**调用示例**:
```bash
echo '{"function": "approve_task", "parameters": {"task_uuid": "abc-123"}}' | ./task-skill
```

### cancel_task
取消任务（任意状态 → cancelled）

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**调用示例**:
```bash
echo '{"function": "cancel_task", "parameters": {"task_uuid": "abc-123"}}' | ./task-skill
```

**任务状态说明**：
- `pending` - 待办（等待 agent 领取）
- `agent_working` - Agent 工作中（agent 已领取，正在处理）
- `agent_review` - Agent 审核中（agent 提交初审后，等待 agent 自己审核确认）
- `human_review` - 人工审核中（agent 审核通过后，提交给人工审核）
- `done` - 完成（人工审核通过）
- `cancelled` - 已取消

### get_task_detail
获取任务详情

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**调用示例**:
```bash
echo '{"function": "get_task_detail", "parameters": {"task_uuid": "abc-123"}}' | ./task-skill
```

### get_version
获取版本号

**参数**: 无

**调用示例**:
```bash
echo '{"function": "get_version"}' | ./task-skill
```

### get_dashboard_stats
获取仪表盘统计信息

**参数**: 无

**调用示例**:
```bash
echo '{"function": "get_dashboard_stats"}' | ./task-skill
```

### recycle_tasks
回收到期未完成的 Agent 任务

**逻辑**:
- 回收 `due_date` 之前创建的、状态为 `agent_working` 的任务
- 回收后任务状态变为 `pending`（重新进入待办池）

**参数**:
- `due_date` (string, 可选): 截止时间，回收此日期前创建的超时任务

**调用示例**:
```bash
echo '{"function": "recycle_tasks", "parameters": {"due_date": "2026-03-22"}}' | ./task-skill
```

### delete_task
物理删除任务（彻底删除）

**注意**: 此操作不可恢复，请谨慎使用！建议先使用 `cancel_task` 将任务状态改为 cancelled（取消），确认不再需要后再使用此功能彻底删除。

**参数**:
- `task_uuid` (string, 必需): 任务 UUID

**调用示例**:
```bash
echo '{"function": "delete_task", "parameters": {"task_uuid": "abc-123"}}' | ./task-skill
```

## 初始化

首次使用或数据库不存在时，在技能目录中执行：

```bash
# 1. 验证安装
./task-skill --version

# 2. 启动 Web 服务器（用于人机协作界面）
./task-skill --server
```