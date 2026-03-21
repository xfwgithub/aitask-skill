#!/usr/bin/env python3
"""
HTMX 前端功能测试脚本
"""
import requests

BASE_URL = "http://localhost:8000"

def test_page(url, expected_title):
    """测试页面访问"""
    response = requests.get(f"{BASE_URL}{url}")
    if response.status_code == 200:
        if expected_title in response.text:
            print(f"✅ {url} - 正常")
            return True
        else:
            print(f"❌ {url} - 内容不匹配")
            return False
    else:
        print(f"❌ {url} - 状态码：{response.status_code}")
        return False

def test_api(url, expected_key):
    """测试 API"""
    response = requests.get(f"{BASE_URL}{url}")
    if response.status_code == 200:
        if expected_key in response.text:
            print(f"✅ {url} - 正常")
            return True
        else:
            print(f"⚠️  {url} - 响应正常但内容不匹配")
            return True
    else:
        print(f"❌ {url} - 状态码：{response.status_code}")
        return False

def main():
    print("=" * 50)
    print("🧪 HTMX 前端功能测试")
    print("=" * 50)
    print()
    
    # 测试页面
    print("📄 测试页面:")
    test_page("/", "仪表盘")
    test_page("/tasks", "任务列表")
    test_page("/projects", "项目管理")
    print()
    
    # 测试 API
    print("🔌 测试 API:")
    test_api("/health", "ok")
    test_api("/api", "任务管理系统")
    test_api("/api/v1/tasks/", "items")
    test_api("/api/v1/projects/", "items")
    print()
    
    # 测试 HTMX 路由
    print("🎯 测试 HTMX 路由:")
    test_api("/api/html/tasks", "task-card")
    test_api("/api/html/projects", "project")
    print()
    
    # 测试静态文件
    print("📦 测试静态资源:")
    response = requests.get(f"{BASE_URL}/static")
    if response.status_code in [200, 403]:  # 403 表示目录存在但禁止访问
        print("✅ /static - 正常")
    else:
        print(f"❌ /static - 状态码：{response.status_code}")
    print()
    
    print("=" * 50)
    print("✨ 测试完成！")
    print("=" * 50)
    print()
    print("访问地址：http://localhost:8000/")
    print("任务列表：http://localhost:8000/tasks")
    print("项目管理：http://localhost:8000/projects")

if __name__ == "__main__":
    main()
