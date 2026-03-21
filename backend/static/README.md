# 静态资源说明

## 目录结构

```
backend/static/
├── css/
│   └── style.css      # 自定义样式
└── js/
    └── app.js         # 应用脚本
```

## style.css

自定义 CSS 样式文件，包含：

### 主题变量
```css
:root {
    --primary-color: #2196F3;    /* 主色调 - 蓝色 */
    --success-color: #4CAF50;    /* 成功 - 绿色 */
    --warning-color: #FF9800;    /* 警告 - 橙色 */
    --danger-color: #F44336;     /* 危险 - 红色 */
}
```

### 主要样式

1. **布局样式**
   - `.sidebar` - 侧边栏（固定宽度 250px）
   - `.main-content` - 主内容区（左边距 270px）
   - `.nav-link` - 导航链接

2. **组件样式**
   - `.task-card` - 任务卡片
   - `.status-badge` - 状态标签
   - `.priority-badge` - 优先级标签
   - `.stats-grid` - 统计网格
   - `.stat-card` - 统计卡片

3. **状态颜色**
   - `.status-pending` - 待处理（黄色）
   - `.status-agent_working` - Agent 执行中（蓝色）
   - `.status-done` - 已完成（绿色）
   - `.status-cancelled` - 已取消（红色）

4. **优先级颜色**
   - `.priority-1` - 紧急（红色）
   - `.priority-2` - 高（橙色）
   - `.priority-3` - 中（蓝色）
   - `.priority-4` - 低（灰色）

5. **响应式设计**
   - 移动端（<768px）：侧边栏变为顶部导航

## app.js

JavaScript 应用脚本，包含：

### HTMX 配置

```javascript
// 显示加载指示器
document.body.addEventListener('htmx:beforeRequest', function() {
    // 显示加载动画
});

// 隐藏加载指示器
document.body.addEventListener('htmx:afterRequest', function() {
    // 隐藏加载动画
});
```

### 工具函数

1. **`refreshTasks()`**
   - 功能：刷新任务列表
   - 使用 HTMX AJAX 更新 `#task-list`

2. **`viewTask(taskId)`**
   - 功能：查看任务详情
   - 参数：任务 ID

3. **`viewProject(projectId)`**
   - 功能：查看项目详情
   - 参数：项目 ID

4. **`editProject(projectId)`**
   - 功能：编辑项目
   - 参数：项目 ID

5. **`closeModal()`**
   - 功能：关闭模态框
   - 清空 `#modal-container`

## 使用方式

### 在 HTML 模板中引用

```html
<head>
    <!-- CSS -->
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <!-- 内容 -->
    
    <!-- JS -->
    <script src="/static/js/app.js"></script>
</body>
```

### 在 FastAPI 中配置

```python
from fastapi.staticfiles import StaticFiles

app.mount("/static", StaticFiles(directory="static"), name="static")
```

## 扩展开发

### 添加新样式

在 `style.css` 末尾添加：

```css
/* 新功能样式 */
.new-feature {
    /* 样式定义 */
}
```

### 添加新脚本

创建 `js/new-feature.js`：

```javascript
// 新功能脚本
function newFeature() {
    // 功能实现
}
```

在模板中引用：

```html
<script src="/static/js/new-feature.js"></script>
```

## 文件大小

- `style.css`: ~3KB（压缩后更小）
- `app.js`: ~1KB（压缩后更小）

总计：< 5KB，加载速度极快！

## 优势

✅ **分离关注点** - HTML、CSS、JS 分离，易于维护  
✅ **可复用性** - 样式和脚本可在多个页面复用  
✅ **缓存优化** - 浏览器可缓存静态文件  
✅ **易于扩展** - 添加新样式或脚本非常简单  
✅ **开发友好** - 代码结构清晰，便于协作
