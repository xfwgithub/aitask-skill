# GitHub Marketplace 发布指南 v0.1.0

## ✅ 已完成的工作

### 1. 构建配置

- ✅ **GitHub Actions 工作流** (`.github/workflows/release.yml`)
  - 自动构建 macOS Intel 和 Apple Silicon 版本
  - 自动生成 SHA256 校验和
  - 自动创建 GitHub Release

- ✅ **Dockerfile**
  - 容器化构建环境
  - 可重复的构建过程

- ✅ **action.yml**
  - GitHub Marketplace 元数据
  - 定义技能基本信息

### 2. 版本发布

- ✅ **版本号**: v0.1.0
- ✅ **Git Tag**: 已创建并推送
- ✅ **Release Notes**: `RELEASE_v0.1.0.md`

### 3. 构建产物

已编译的二进制文件：
- `task-skill-darwin-amd64` (20MB) - macOS Intel
- `task-skill-darwin-arm64` (19MB) - macOS Apple Silicon

## 📋 发布步骤

### 步骤 1: 验证 GitHub Actions

推送 tag 后，GitHub Actions 会自动触发：

1. 访问 https://github.com/xfwgithub/aitask-skill/actions
2. 查看 "Build and Release" 工作流
3. 等待构建完成（约 2-5 分钟）
4. 验证 Release 已创建

### 步骤 2: 检查 Release

构建完成后，访问：
https://github.com/xfwgithub/aitask-skill/releases/tag/v0.1.0

应该包含：
- ✅ 发布说明
- ✅ 4 个附件：
  - task-skill-darwin-amd64
  - task-skill-darwin-arm64
  - task-skill-darwin-amd64.sha256
  - task-skill-darwin-arm64.sha256

### 步骤 3: 提交到 GitHub Marketplace

#### 3.1 准备 Marketplace 页面

访问 GitHub Marketplace 提交页面：
https://github.com/marketplace/actions/create

填写信息：

**基本信息:**
```
Name: Task Management Skill
Description: 零依赖、高性能的任务管理技能，支持 Web 界面和 CLI 模式
Category: Productivity
```

**技术细节:**
```
Runs on: macOS
Docker Image: Yes
```

**品牌标识:**
```
Icon: check-square (或上传自定义图标)
Color: blue
```

#### 3.2 提交审核

1. 填写所有必填字段
2. 上传截图（可选但推荐）
3. 提供测试说明
4. 提交审核

审核时间：通常 1-3 个工作日

### 步骤 4: 测试安装

发布后，用户可以通过以下方式安装：

**方式 1: GitHub Actions**
```yaml
- uses: xfwgithub/aitask-skill@v0.1.0
```

**方式 2: 直接下载**
```bash
# macOS Apple Silicon
curl -L https://github.com/xfwgithub/aitask-skill/releases/download/v0.1.0/task-skill-darwin-arm64 -o task-skill
chmod +x task-skill

# macOS Intel
curl -L https://github.com/xfwgithub/aitask-skill/releases/download/v0.1.0/task-skill-darwin-amd64 -o task-skill
chmod +x task-skill
```

## 🔍 验证清单

### 发布前验证

- [x] 代码已提交到 main 分支
- [x] Tag v0.1.0 已创建并推送
- [x] 所有测试通过
- [x] 文档完整
- [x] 版本号正确

### 发布后验证

- [ ] GitHub Actions 构建成功
- [ ] Release 包含所有文件
- [ ] 下载链接有效
- [ ] 校验和匹配
- [ ] 安装说明正确
- [ ] Marketplace 页面显示正常

## 📊 监控和维护

### 监控指标

1. **下载量**
   - GitHub Releases 页面查看
   - 按平台和版本统计

2. **用户反馈**
   - GitHub Issues
   - GitHub Discussions
   - Star 数量

3. **构建状态**
   - GitHub Actions 历史记录
   - 构建成功率

### 更新流程

发布新版本（如 v0.1.1）：

```bash
# 1. 更新代码
git add .
git commit -m "fix: bug fixes"

# 2. 打新标签
git tag -a v0.1.1 -m "Task Management Skill v0.1.1"

# 3. 推送
git push origin main
git push origin v0.1.1

# 4. GitHub Actions 会自动构建和发布
```

## 🐛 故障排除

### 问题 1: GitHub Actions 失败

**检查点:**
1. 查看 Actions 日志
2. 确认 Go 版本正确
3. 验证依赖可访问

**解决方案:**
```bash
# 本地测试构建
cd task-management
go build -o task-skill .
```

### 问题 2: Release 未创建

**可能原因:**
- GitHub Actions 配置错误
- 权限问题

**解决方案:**
1. 检查 `.github/workflows/release.yml`
2. 确认 `GITHUB_TOKEN` 权限
3. 手动触发工作流

### 问题 3: Marketplace 审核失败

**常见原因:**
- 文档不完整
- 缺少测试说明
- 安全问题

**解决方案:**
1. 完善 README
2. 添加详细测试步骤
3. 修复安全问题
4. 重新提交

## 📈 推广建议

### 1. 社交媒体

- Twitter/X: 发布 release 公告
- LinkedIn: 专业网络分享
- Reddit: r/ClaudeAI, r/golang

### 2. 技术社区

- Product Hunt: 产品发布
- Hacker News: Show HN
- 掘金/思否：中文社区

### 3. 文档优化

- 添加使用视频教程
- 提供在线 Demo
- 编写详细教程

## 🔗 相关链接

- **GitHub 仓库**: https://github.com/xfwgithub/aitask-skill
- **Releases**: https://github.com/xfwgithub/aitask-skill/releases
- **Actions**: https://github.com/xfwgithub/aitask-skill/actions
- **Marketplace**: https://github.com/marketplace

## 📞 支持

如有问题，请：
1. 查看文档：README.md
2. 提交 Issue: https://github.com/xfwgithub/aitask-skill/issues
3. 参与讨论：https://github.com/xfwgithub/aitask-skill/discussions

---

**发布状态**: ✅ 已发布 v0.1.0
**下次更新**: v0.1.1 (计划添加 Linux 支持)
