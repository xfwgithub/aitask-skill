# Task Management Skill

零依赖、高性能的任务管理技能，使用 Go 语言实现。

## 安装

### 下载预编译二进制

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
```

### 从源码编译

```bash
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill/go-skill
go build -o task-skill
```

## 使用

直接对 AI 助手说：
- "帮我创建一个紧急任务，明天要完成项目评审"
- "查询我所有待处理的任务"
- "把任务'审查文档'标记为已完成"
- "统计一下本周完成了多少任务"

## 项目结构

```
.
├── SKILL.md          # 技能定义（必需）
├── README.md         # 使用文档
└── go-skill/         # Go 实现
    ├── main.go
    ├── utils.go
    ├── go.mod
    └── go.sum
```

## 许可证

MIT License
