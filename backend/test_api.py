#!/usr/bin/env python3
"""
API 测试脚本

测试任务管理系统的各项 API 功能
"""
import requests
import json

BASE_URL = "http://localhost:8000/api/v1"


def test_health():
    """测试健康检查"""
    response = requests.get("http://localhost:8000/health")
    print(f"健康检查：{response.json()}")
    assert response.status_code == 200
    print("✓ 健康检查通过\n")


def test_create_task():
    """测试创建任务"""
    task_data = {
        "title": "测试任务",
        "description": "这是一个 API 测试任务",
        "priority": 2,
        "tags": ["测试", "API"],
    }
    
    response = requests.post(f"{BASE_URL}/tasks/", json=task_data)
    print(f"创建任务响应：{json.dumps(response.json(), indent=2, ensure_ascii=False)}")
    assert response.status_code == 200
    task = response.json()
    print("✓ 任务创建成功\n")
    return task


def test_get_tasks(task_uuid):
    """测试查询任务"""
    # 查询单个任务
    response = requests.get(f"{BASE_URL}/tasks/{task_uuid}")
    print(f"查询任务详情：{json.dumps(response.json(), indent=2, ensure_ascii=False)}")
    assert response.status_code == 200
    
    # 查询任务列表
    response = requests.get(f"{BASE_URL}/tasks/", params={"page": 1, "page_size": 10})
    print(f"查询任务列表：总数={response.json()['total']}")
    assert response.status_code == 200
    print("✓ 任务查询成功\n")


def test_update_task_status(task_uuid):
    """测试更新任务状态"""
    response = requests.post(f"{BASE_URL}/tasks/{task_uuid}/status", params={"status": "agent_working"})
    print(f"更新任务状态：{response.json()}")
    assert response.status_code == 200
    print("✓ 任务状态更新成功\n")


def test_create_project():
    """测试创建项目"""
    project_data = {
        "name": "测试项目",
        "description": "API 测试项目",
        "color": "#4CAF50",
    }
    
    response = requests.post(f"{BASE_URL}/projects/", json=project_data)
    print(f"创建项目：{json.dumps(response.json(), indent=2, ensure_ascii=False)}")
    assert response.status_code == 200
    print("✓ 项目创建成功\n")
    return response.json()


def test_get_projects():
    """测试获取项目列表"""
    response = requests.get(f"{BASE_URL}/projects/")
    print(f"项目列表：{len(response.json())} 个项目")
    assert response.status_code == 200
    print("✓ 项目列表获取成功\n")


def test_get_reports():
    """测试报表 API"""
    # 获取仪表盘统计
    response = requests.get(f"{BASE_URL}/reports/stats")
    print(f"仪表盘统计：{json.dumps(response.json(), indent=2, ensure_ascii=False)}")
    assert response.status_code == 200
    print("✓ 报表获取成功\n")


def main():
    print("=" * 60)
    print("任务管理系统 API 测试")
    print("=" * 60 + "\n")
    
    try:
        # 健康检查
        test_health()
        
        # 任务测试
        task = test_create_task()
        task_uuid = task["uuid"]
        test_get_tasks(task_uuid)
        test_update_task_status(task_uuid)
        
        # 项目测试
        test_create_project()
        test_get_projects()
        
        # 报表测试
        test_get_reports()
        
        print("=" * 60)
        print("所有测试通过！✓")
        print("=" * 60)
        
    except requests.exceptions.ConnectionError:
        print("错误：无法连接到服务器，请确保服务正在运行")
        print("运行：python start.py")
    except AssertionError as e:
        print(f"测试失败：{e}")
    except Exception as e:
        print(f"发生错误：{e}")


if __name__ == "__main__":
    main()
