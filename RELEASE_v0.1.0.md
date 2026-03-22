# Task Management Skill v0.1.0

## 🎉 初始预览版发布

这是 Task Management Skill 的第一个预览版本，支持 macOS (Intel 和 Apple Silicon)。

## ✨ 功能特性

- 📊 **Web 界面** - 直观的仪表盘和任务管理
- 💻 **CLI 模式** - 命令行快速操作
- 💾 **数据持久化** - SQLite 数据库存储
- 🚀 **零依赖** - 编译为单一二进制文件
- ⚡ **高性能** - Go 语言原生支持

## 📦 安装

### macOS (Apple Silicon M1/M2/M3)

```bash
# 下载
curl -L https://github.com/xfwgithub/aitask-skill/releases/download/v0.1.0/task-skill-darwin-arm64 -o task-skill

# 添加执行权限
chmod +x task-skill

# 移动到家目录
mv task-skill ~/bin/

# 验证安装
task-skill --version
```

### macOS (Intel)

```bash
# 下载
curl -L https://github.com/xfwgithub/aitask-skill/releases/download/v0.1.0/task-skill-darwin-amd64 -o task-skill

# 添加执行权限
chmod +x task-skill

# 移动到家目录
mv task-skill ~/bin/

# 验证安装
task-skill --version
```

## 🚀 快速开始

### Web 模式

```bash
# 启动 Web 服务器
task-skill --server

# 浏览器访问
# http://localhost:8080
```

### CLI 模式

```bash
# 创建任务
echo '{"function": "create_task", "parameters": {"title": "我的任务"}}' | task-skill

# 查询任务
echo '{"function": "query_tasks"}' | task-skill

# 获取统计
echo '{"function": "get_dashboard_stats"}' | task-skill
```

## 📋 校验

### 验证文件完整性

**macOS ARM64:**
```bash
shasum -a 256 task-skill-darwin-arm64
# 应该匹配 task-skill-darwin-arm64.sha256 中的值
```

**macOS AMD64:**
```bash
shasum -a 256 task-skill-darwin-amd64
# 应该匹配 task-skill-darwin-amd64.sha256 中的值
```

## 🔧 配置

### 环境变量

```bash
# 配置端口（默认：8080）
export TASK_SKILL_PORT=3000

# 配置数据库路径（默认：tasks.db）
export TASK_SKILL_DB_PATH=/path/to/tasks.db
```

## 📖 文档

- [完整文档](README.md)
- [发布指南](PUBLISHING.md)
- [技能包配置](skill-package.yaml)

## 🐛 已知问题

这是预览版本，可能存在以下限制：

- 仅支持 macOS
- 部分功能仍在开发中
- 文档可能不完整

## 🗺️ 路线图

- [ ] Linux 支持
- [ ] Windows 支持
- [ ] 更多任务状态
- [ ] 任务优先级排序
- [ ] 任务标签管理
- [ ] 数据导出功能

## 📝 变更日志

### v0.1.0 (2026-03-22)

**新增功能:**
- ✅ 任务创建、查询、更新
- ✅ Web 界面和仪表盘
- ✅ CLI 模式
- ✅ SQLite 数据库持久化
- ✅ 端口配置支持

**技术栈:**
- Go 1.21+
- Echo Framework
- SQLite (modernc.org/sqlite)
- HTMX

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

- 报告 Bug: https://github.com/xfwgithub/aitask-skill/issues
- 功能请求：https://github.com/xfwgithub/aitask-skill/discussions

## 📄 许可证

MIT License - 详见 [LICENSE](LICENSE)

## 🙏 致谢

感谢所有贡献者和用户！

---

**下载地址:**
- [macOS ARM64](https://github.com/xfwgithub/aitask-skill/releases/download/v0.1.0/task-skill-darwin-arm64)
- [macOS AMD64](https://github.com/xfwgithub/aitask-skill/releases/download/v0.1.0/task-skill-darwin-amd64)
