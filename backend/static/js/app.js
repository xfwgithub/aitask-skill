// HTMX 全局配置

// 显示加载指示器
document.body.addEventListener('htmx:beforeRequest', function() {
    const loadingIndicator = document.getElementById('htmx-loading');
    if (loadingIndicator) {
        loadingIndicator.style.display = 'block';
    }
});

// 隐藏加载指示器
document.body.addEventListener('htmx:afterRequest', function() {
    const loadingIndicator = document.getElementById('htmx-loading');
    if (loadingIndicator) {
        setTimeout(() => {
            loadingIndicator.style.display = 'none';
        }, 200);
    }
});

// 自动刷新任务列表
function refreshTasks() {
    if (typeof htmx !== 'undefined') {
        htmx.ajax('GET', '/api/html/tasks', '#task-list');
    }
}

// 查看任务详情
function viewTask(taskId) {
    // TODO: 实现任务详情查看
    console.log('查看任务:', taskId);
}

// 查看项目详情
function viewProject(projectId) {
    // TODO: 实现项目详情查看
    console.log('查看项目:', projectId);
}

// 编辑项目
function editProject(projectId) {
    // TODO: 实现项目编辑
    console.log('编辑项目:', projectId);
}

// 关闭模态框
function closeModal() {
    const modalContainer = document.getElementById('modal-container');
    if (modalContainer) {
        modalContainer.innerHTML = '';
    }
}
