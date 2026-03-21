# Skill 发布指南

本文档说明如何将 aitask-skill 发布到 skill 市场。

## 📦 发布前准备

### 1. 确保文件完整

检查以下文件是否存在且内容正确：

- ✅ `SKILL.md` - 技能定义文件（包含 frontmatter）
- ✅ `skill.json` - 技能元数据
- ✅ `README.md` - 项目说明文档
- ✅ `.gitignore` - Git 忽略文件
- ✅ `backend/requirements.txt` - Python 依赖

### 2. 更新技能信息

在 `SKILL.md` 和 `skill.json` 中更新：

- 作者信息（name, email）
- GitHub 仓库地址
- 版本号

### 3. 测试技能

确保技能可以正常工作：

```bash
# 安装依赖
cd backend
pip install -r requirements.txt

# 初始化数据库
python start.py --init-db

# 启动服务
python start.py

# 访问 Web 前端
# http://localhost:8000/
```

## 🚀 发布步骤

### 步骤 1：创建 GitHub 仓库

```bash
# 在 GitHub 上创建新仓库
# 仓库名：aitask-skill
# 可见性：Public
```

### 步骤 2：推送代码到 GitHub

```bash
# 初始化 Git（如果还没有）
git init

# 添加所有文件
git add .

# 提交
git commit -m "Initial commit: Task Management Skill v1.0.0"

# 添加远程仓库（替换为你的仓库地址）
git remote add origin https://github.com/YOUR_USERNAME/aitask-skill.git

# 推送到 GitHub
git push -u origin main
```

### 步骤 3：添加 Git 标签

```bash
# 创建版本标签
git tag -a v1.0.0 -m "Release version 1.0.0"

# 推送标签
git push origin v1.0.0
```

### 步骤 4：发布到 Skill 市场

根据不同平台的要求进行发布：

#### 方式 1：Claude Skills

1. 访问 [Claude Skills Marketplace](https://claude.ai/skills)
2. 点击 "Submit a Skill"
3. 填写技能信息：
   - Name: Task Management
   - Description: AI 助手任务管理技能 - 创建、查询、更新任务的标准化技能包
   - Repository URL: https://github.com/YOUR_USERNAME/aitask-skill
   - SKILL.md 文件路径：./SKILL.md
4. 提交审核

#### 方式 2：Cursor Rules

1. 访问 [Cursor Rules Marketplace](https://cursor.sh/rules)
2. 提交你的技能包
3. 提供 GitHub 仓库地址

#### 方式 3：其他平台

根据具体平台的要求进行发布。

## 📝 技能包结构

发布后的技能包结构：

```
aitask-skill/
├── SKILL.md              # 技能定义（核心文件）
├── skill.json            # 技能元数据
├── README.md             # 使用文档
├── .gitignore            # Git 忽略文件
├── backend/
│   ├── app/
│   │   ├── skills/       # Skill 实现
│   │   │   └── task_skill.py
│   │   ├── models/       # 数据模型
│   │   ├── db/           # 数据库连接
│   │   └── ...
│   ├── requirements.txt  # Python 依赖
│   └── ...
└── examples/             # 使用示例（可选）
```

## ✅ 发布检查清单

发布前请确认：

- [ ] SKILL.md 文件格式正确（包含 frontmatter）
- [ ] skill.json 元数据完整
- [ ] README.md 文档清晰
- [ ] 代码可以正常运行
- [ ] 所有依赖已列出
- [ ] GitHub 仓库已创建
- [ ] 版本号已更新
- [ ] 作者信息已更新
- [ ] 许可证文件（可选）

## 🔄 更新技能

当需要更新技能时：

```bash
# 1. 更新 SKILL.md 和 skill.json 中的版本号
# 2. 提交更改
git add .
git commit -m "Update to version 1.0.1"

# 3. 创建新标签
git tag -a v1.0.1 -m "Release version 1.0.1"

# 4. 推送
git push origin main
git push origin v1.0.1

# 5. 在 Skill 市场更新版本
```

## 📞 支持

如有问题，请查看：

- [Skill 创建文档](https://github.com/anthropics/claude-skills)
- [GitHub Issues](https://github.com/YOUR_USERNAME/aitask-skill/issues)

## 📄 许可证

MIT License - 详见 LICENSE 文件
