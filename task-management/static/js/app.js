// 显示创建模态框
function showCreateModal() {
    document.getElementById('createModal').style.display = 'flex';
}

// 关闭创建模态框
function closeCreateModal() {
    document.getElementById('createModal').style.display = 'none';
    document.getElementById('createForm').reset();
}

// 创建任务
function createTask(event) {
    event.preventDefault();
    
    const data = {
        title: document.getElementById('title').value,
        description: document.getElementById('description').value,
        priority: parseInt(document.getElementById('priority').value),
        project: document.getElementById('project').value,
        assignee_name: document.getElementById('assignee').value,
        tags: document.getElementById('tags').value.split(',').map(t => t.trim()).filter(t => t)
    };
    
    fetch('/api/tasks', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(data)
    })
    .then(r => r.json())
    .then(result => {
        if (result.error) {
            alert('创建失败：' + result.error);
        } else {
            alert('任务创建成功！');
            closeCreateModal();
            window.location.reload();
        }
    })
    .catch(err => {
        alert('创建失败：' + err);
    });
}

// 更新任务状态
function updateStatus(uuid, status) {
    fetch(`/api/tasks/${uuid}/status`, {
        method: 'PUT',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({new_status: status})
    })
    .then(() => {
        window.location.reload();
    });
}

// 工具函数
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleString('zh-CN');
}

// 点击模态框外部关闭
window.onclick = function(event) {
    const modal = document.getElementById('createModal');
    if (event.target === modal) {
        closeCreateModal();
    }
}
