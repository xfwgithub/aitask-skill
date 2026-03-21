# HTMX 前端使用指南

## 🎉 完成！

Web 前端已经开发完成！使用超轻量的 **HTMX** 技术栈。

## 📦 技术栈

- **HTMX** (14KB) - 实现动态交互
- **PicoCSS** - 轻量级 CSS 框架
- **Jinja2** - 模板引擎
- **FastAPI** - 后端服务

**总大小**: < 50KB（无需 npm、无需构建）

## 🚀 访问前端

### 1. 启动服务

```bash
cd backend
python start.py
```

### 2. 访问页面

- **首页/仪表盘**: http://localhost:8000/
- **任务列表**: http://localhost:8000/tasks
- **项目管理**: http://localhost:8000/projects
- **API 文档**: http://localhost:8000/docs

## ✨ 功能特性

### 仪表盘
- 📊 实时统计卡片（总任务数、待处理、执行中、已完成）
- 📝 最近任务列表
- ⏳ 待处理任务

### 任务管理
- ✅ 任务列表（支持搜索、筛选）
- ➕ 创建任务（弹窗表单）
- ✓ 完成任务（一键完成）
- 🏷️ 任务标签
- 🤖 Agent 类型显示

### 项目管理
- 📁 项目列表
- ➕ 创建项目
- 📊 项目任务统计

## 🎯 HTMX 特性

### 实时搜索
```html
<input type="search" 
       hx-get="/api/tasks"
       hx-trigger="input changed delay:300ms"
       hx-target="#task-list">
```
输入时自动搜索，300ms 防抖

### 动态筛选
```html
<select name="status"
        hx-get="/api/tasks"
        hx-trigger="change"
        hx-target="#task-list">
    <option value="">所有状态</option>
    <option value="pending">待处理</option>
    ...
</select>
```
选择时自动刷新列表

### 无刷新提交
```html
<form hx-post="/api/tasks"
      hx-target="#task-list"
      hx-swap="afterbegin">
    <input name="title" required>
    <button type="submit">创建</button>
</form>
```
提交后自动更新列表，无需刷新页面

### 一键操作
```html
<button hx-post="/api/tasks/{uuid}/status?status=done"
        hx-target="#task-{id}"
        hx-swap="outerHTML">
    ✓
</button>
```
点击完成任务，自动替换卡片

## 📁 文件结构

```
backend/
├── templates/           # HTML 模板
│   ├── base.html       # 基础模板
│   ├── index.html      # 仪表盘
│   ├── tasks.html      # 任务列表页
│   ├── task_list.html  # 任务列表片段
│   ├── task_create.html # 创建任务表单
│   ├── projects.html   # 项目列表页
│   └── project_list.html # 项目列表片段
├── static/             # 静态资源
│   ├── css/
│   └── js/
└── app/
    └── api/
        └── htmx_routes.py  # HTMX 路由
```

## 🎨 样式特性

- ✅ 响应式设计（支持手机、平板、桌面）
- ✅ 暗色模式支持
- ✅ 状态标签颜色区分
- ✅ 优先级标签
- ✅ 加载动画
- ✅ 悬停效果

## 💡 优势

### 对比 React/Vue 方案

| 特性 | HTMX | React/Vue |
|------|------|-----------|
| 打包大小 | <50KB | >500KB |
| 构建工具 | 无需 | Webpack/Vite |
| 学习曲线 | 10 分钟 | 数天 |
| 开发速度 | 极快 | 中等 |
| 性能 | 优秀 | 良好 |

### 代码量对比

**HTMX** (10 行):
```html
<button hx-post="/api/tasks/1/complete"
        hx-target="#task-1"
        hx-swap="outerHTML">
  完成
</button>
```

**React** (30+ 行):
```jsx
function TaskButton({ id }) {
  const [loading, setLoading] = useState(false);
  
  const handleClick = async () => {
    setLoading(true);
    await fetch(`/api/tasks/${id}/complete`, { method: 'POST' });
    // 更新状态、重新获取...
  };
  
  return <button onClick={handleClick} disabled={loading}>完成</button>;
}
```

## 🔧 扩展开发

### 添加新页面

1. 在 `templates/` 创建 HTML 文件
2. 在 `main.py` 添加路由
3. 在侧边栏添加导航链接

### 添加新功能

1. 创建 HTML 片段模板
2. 在 `htmx_routes.py` 添加路由
3. 使用 HTMX 属性调用

### 自定义样式

编辑 `templates/base.html` 中的 `<style>` 部分

## 📝 示例

### 添加"删除任务"功能

```html
<!-- 在 task_list.html 中添加 -->
<button hx-delete="/api/tasks/{{ task.uuid }}"
        hx-target="#task-{{ task.id }}"
        hx-swap="outerHTML"
        onclick="return confirm('确定删除？')">
  🗑️ 删除
</button>
```

### 添加批量操作

```html
<!-- 在 tasks.html 中添加 -->
<div class="action-bar">
  <button hx-post="/api/tasks/bulk-delete"
          hx-target="#task-list"
          hx-confirm="确定删除选中的任务？">
    批量删除
  </button>
</div>
```

## 🎯 总结

✅ **开发完成** - 所有核心功能已实现  
✅ **超轻量** - 无需构建，零配置  
✅ **易维护** - 简单直观的代码  
✅ **高性能** - 快速加载和响应  
✅ **易扩展** - 添加新功能非常简单  

现在您可以访问 http://localhost:8000/ 开始使用了！
