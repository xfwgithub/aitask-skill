# Skill 初始化指南

## 官方 Skills 的初始化方式

官方 Anthropic Skills 采用**文档驱动**架构，技能本身就是一个 `SKILL.md` 文件。

### 标准结构

```
skills/
├── skill-name/           # 技能目录
│   └── SKILL.md         # 技能定义（必需）
├── another-skill/
│   ├── SKILL.md
│   └── scripts/         # 可选：辅助脚本
└── third-skill/
    ├── SKILL.md
    └── references/      # 可选：参考文档
```

### SKILL.md 结构

```markdown
---
name: skill-name
description: 技能描述（触发条件和功能说明）
license: 许可证信息
---

# 技能名称

## 概述
技能的主要用途和触发场景

## 使用场景
- 何时使用此技能
- 什么情况下应该触发

## 工作流程
1. 第一步
2. 第二步
3. 第三步

## 输出格式
说明期望的输出格式

## 示例
### 示例 1
输入：...
输出：...
```

## 本项目的初始化方式

我们的 `aitask-skill` 项目采用**应用驱动**架构，与官方文档型技能不同：

### 项目结构

```
aitask-skill/
├── README.md              # 项目文档
├── SKILL.md               # Claude Skills 定义（触发器）
└── task-management/       # 实际技能实现
    ├── main.go           # 主程序
    ├── server.go         # Web 服务器
    ├── db.go             # 数据库
    └── ...
```

### 初始化步骤

#### 1. 作为 Claude Skill 使用

在 Claude Desktop 或 Claude 应用中配置：

1. 打开 Claude 设置
2. 找到 Skills 配置
3. 添加技能路径：`/path/to/aitask-skill`
4. Claude 会读取 `SKILL.md` 了解何时触发

#### 2. 独立运行（推荐）

**Web 模式：**
```bash
cd task-management

# 安装依赖
go mod download

# 启动服务器
./start.sh

# 浏览器访问 http://localhost:8080
```

**CLI 模式：**
```bash
cd task-management

# 创建任务
echo '{"function": "create_task", "parameters": {"title": "测试任务"}}' | ./task-skill

# 查询任务
echo '{"function": "query_tasks"}' | ./task-skill
```

#### 3. 在 IDE 中集成

**VS Code / Cursor / Trae:**
```bash
# 在 IDE 终端中运行
cd task-management
./start.sh
```

**JetBrains IDEs:**
1. Run → Edit Configurations
2. 添加 Go Build 配置
3. Working directory: `task-management`
4. Program arguments: `--server`

## 两种模式对比

| 特性 | 官方 Skills | aitask-skill |
|------|-----------|-------------|
| **形式** | 文档 (`.md`) | 应用程序 (`.go`) |
| **执行** | Claude 阅读文档 | 独立运行的服务 |
| **触发** | 基于描述自动触发 | 基于描述自动触发 |
| **功能** | 指导 Claude 行为 | 提供实际功能 |
| **依赖** | 无 | Go 运行时 |
| **交互** | 文本对话 | Web UI + API |

## 最佳实践

### 1. 作为 Claude Skill 部署

确保 `SKILL.md` 包含：
- 清晰的触发条件
- 明确的功能说明
- 使用示例

### 2. 作为独立应用部署

- 提供 Web 界面
- 支持 CLI 调用
- 完整的 API 文档

### 3. 混合模式（推荐）

结合两者优势：
- `SKILL.md` 作为触发器
- 实际功能由应用提供
- Claude 作为中间层调用应用

## 环境变量配置

```bash
# 配置服务端口
export TASK_SKILL_PORT=3000

# 配置数据库路径
export TASK_SKILL_DB_PATH=/path/to/tasks.db
```

## 测试技能

### 测试 CLI 模式
```bash
cd task-management

# 测试创建任务
echo '{"function": "create_task", "parameters": {"title": "测试", "priority": 1}}' | ./task-skill

# 测试查询任务
echo '{"function": "query_tasks", "parameters": {"status": "pending"}}' | ./task-skill

# 测试统计
echo '{"function": "get_dashboard_stats"}' | ./task-skill
```

### 测试 Web 模式
```bash
cd task-management
./start.sh

# 浏览器访问 http://localhost:8080
# 或使用 curl 测试 API
curl http://localhost:8080/api/stats
```

## 故障排除

### 技能不触发
检查 `SKILL.md` 的 `description` 字段，确保包含：
- 明确的触发场景
- 关键词匹配

### 服务无法启动
```bash
# 检查端口占用
lsof -i :8080

# 使用其他端口
TASK_SKILL_PORT=3000 ./start.sh
```

### 编译失败
```bash
# 清理依赖
cd task-management
go mod tidy
go clean

# 重新编译
go build -o task-skill .
```

## 参考资源

- [Anthropic Skills 官方仓库](https://github.com/anthropics/skills)
- [Echo Framework 文档](https://echo.labstack.com/)
- [HTMX 文档](https://htmx.org/)
