# 发布到技能市场指南

## 发布渠道

### 1. GitHub 开源发布（推荐）

将你的技能开源到 GitHub，让用户可以直接安装使用。

#### 步骤：

**1.1 准备项目**
```bash
# 确保项目结构完整
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill

# 验证必要文件
ls -la
# 应该包含：
# - README.md
# - SKILL.md
# - task-management/
```

**1.2 创建 GitHub Release**
```bash
# 打标签
git tag -a v1.0.0 -m "Task Management Skill v1.0.0"
git push origin v1.0.0

# 或使用 GitHub UI
# https://github.com/xfwgithub/aitask-skill/releases/new
```

**1.3 发布到 GitHub Marketplace**
- 访问：https://github.com/marketplace/actions
- 点击 "Create a new action"
- 填写技能信息
- 提交审核

### 2. Claude Code 插件市场

用户可以通过以下命令安装你的技能：

```bash
# 用户安装命令
/plugin marketplace add xfwgithub/aitask-skill

# 或直接安装
/plugin install task-management@xfwgithub/aitask-skill
```

#### 注册为插件市场：

**2.1 创建 .claude 配置**
在项目根目录创建 `.claude/commands/install.md`：

```markdown
---
name: task-management
description: 任务管理技能，用于创建、查询、更新任务
version: 1.0.0
---

# 安装 Task Management Skill

## 自动安装

```bash
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill/task-management
./quick-init.sh
```

## 手动安装

```bash
git clone https://github.com/xfwgithub/aitask-skill.git
cd aitask-skill/task-management
go mod download
go build -o task-skill .
./task-skill --server
```

## 使用

安装完成后，可以直接对 Claude 说：
- "创建一个任务，明天要完成项目评审"
- "查询我所有待处理的任务"
- "把任务标记为已完成"
```

**2.2 提交到 Claude Code 插件目录**
- 访问：https://github.com/anthropics/skills
- 提交 PR 将你的技能添加到示例列表
- 或参考官方文档：https://skills.anthropic.com/

### 3. 社区技能市场

#### 3.1 Skills Marketplace (skills-mp.com)

**提交步骤：**
1. 访问 https://skills-mp.com/
2. 注册开发者账号
3. 提交技能信息：
   - 技能名称：Task Management
   - GitHub 仓库：https://github.com/xfwgithub/aitask-skill
   - 技能描述：任务管理技能
   - 分类：Productivity / Development

#### 3.2 Claude Skills 中文社区

**提交步骤：**
1. 访问 https://claude.cn/
2. 进入 Skills 专区
3. 提交技能教程和链接

### 4. 官方 Anthropic Skills 仓库

提交到官方示例仓库（需要审核）：

**4.1 准备提交材料**
- 完整的 SKILL.md 文件
- 清晰的文档说明
- 示例用例
- 测试脚本

**4.2 提交 Pull Request**
```bash
# Fork 官方仓库
git clone https://github.com/anthropics/skills.git

# 添加你的技能
cp -r /path/to/aitask-skill skills/skills/task-management

# 提交 PR
git add skills/task-management
git commit -m "Add task-management skill"
git push origin main

# 在 GitHub 上创建 Pull Request
# https://github.com/anthropics/skills/pulls
```

## 发布前检查清单

### 必要文件
- [ ] `README.md` - 项目说明
- [ ] `SKILL.md` - 技能定义
- [ ] `task-management/` - 技能实现
- [ ] `LICENSE` - 开源许可证
- [ ] `.gitignore` - Git 忽略文件

### 文档完整性
- [ ] 安装说明
- [ ] 使用示例
- [ ] API 文档
- [ ] 故障排除
- [ ] 贡献指南

### 代码质量
- [ ] 通过所有测试
- [ ] 无安全漏洞
- [ ] 代码注释完整
- [ ] 遵循 Go 语言规范

### 测试验证
- [ ] CLI 模式测试通过
- [ ] Web 模式测试通过
- [ ] 数据库持久化正常
- [ ] 跨平台测试（macOS, Linux, Windows）

## 开源许可证选择

### 推荐许可证

**MIT License**（当前使用）
- 最宽松，允许商业使用
- 只需保留版权声明
- 适合技能市场

**Apache 2.0**
- 包含专利授权
- 更正式的法律保护
- 适合企业级项目

**GPL v3**
- 要求衍生作品也开源
- 保护开源生态
- 适合社区项目

## 版本管理

### Semantic Versioning

遵循语义化版本规范：

```
主版本号。次版本号.修订号
Major.Minor.Patch
```

**示例：**
- `v1.0.0` - 初始发布
- `v1.1.0` - 新增功能
- `v1.1.1` - Bug 修复
- `v2.0.0` - 重大变更

### 发布流程

```bash
# 1. 更新版本号
# 在 skill.config.json 中更新 version

# 2. 更新 CHANGELOG.md
# 记录变更内容

# 3. 提交并打标签
git add .
git commit -m "release: v1.1.0 - Add new features"
git tag -a v1.1.0 -m "Version 1.1.0"
git push origin main
git push origin v1.1.0

# 4. 创建 GitHub Release
# https://github.com/xfwgithub/aitask-skill/releases/new
```

## 推广技能

### 1. 社交媒体
- Twitter/X: 发布技能介绍
- LinkedIn: 专业网络分享
- Reddit: r/ClaudeAI, r/programming
- 微博/知乎：中文社区

### 2. 技术社区
- Product Hunt: 产品发布
- Hacker News: 技术讨论
- Indie Hackers: 独立开发者社区

### 3. 中文社区
- 掘金：技术文章
- 思否：教程分享
- V2EX: 创意工作者社区

### 4. 文档优化
- 添加演示视频
- 提供在线 Demo
- 编写详细教程
- 制作使用案例

## 维护和更新

### 定期更新
- 修复 Bug
- 添加新功能
- 响应用户反馈
- 保持依赖更新

### 用户支持
- GitHub Issues: 问题反馈
- Discussions: 讨论区
- Email: 技术支持
- Discord/Slack: 社区交流

### 性能监控
- 收集使用数据
- 分析用户行为
- 优化性能瓶颈
- 改进用户体验

## 商业化选项

### 1. 开源 + 付费支持
- 基础功能免费
- 企业支持收费
- 定制开发收费

### 2. SaaS 版本
- 自托管免费
- 云服务收费
- 提供托管方案

### 3. 高级功能
- 基础版免费
- 专业版收费
- 企业版定制

## 参考资源

- [Anthropic Skills 官方文档](https://skills.anthropic.com/)
- [GitHub Marketplace](https://github.com/marketplace)
- [Claude Code 插件](https://github.com/anthropics/claude-code)
- [Skills Marketplace](https://skills-mp.com/)
- [开源许可证选择](https://choosealicense.com/)

## 联系支持

如有问题，请通过以下方式联系：
- GitHub Issues: https://github.com/xfwgithub/aitask-skill/issues
- Email: [你的邮箱]
- 文档：https://github.com/xfwgithub/aitask-skill/wiki
