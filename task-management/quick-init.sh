#!/bin/bash

# Task Management Skill 快速初始化脚本
# 用于 Agent 自动执行，无交互式

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "🚀 初始化 Task Management Skill..."
echo ""

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误：未找到 Go 环境"
    echo "请先安装 Go: https://golang.org/dl/"
    exit 1
fi

echo "✓ Go 环境：$(go version)"

# 初始化依赖
echo "📦 安装依赖..."
if [ ! -f "go.mod" ]; then
    go mod init task-skill
fi

go mod download
echo "✓ 依赖安装完成"

# 编译
echo "🔨 编译程序..."
go build -o task-skill .
echo "✓ 编译成功"

# 验证
echo "✅ 验证安装..."
if [ -f "task-skill" ]; then
    echo "✓ 二进制文件已创建"
    
    # 测试 CLI 模式
    RESULT=$(echo '{"function": "get_dashboard_stats"}' | ./task-skill 2>&1)
    if echo "$RESULT" | grep -q "total"; then
        echo "✓ CLI 模式测试通过"
    else
        echo "⚠ CLI 模式测试失败"
    fi
else
    echo "❌ 编译失败"
    exit 1
fi

echo ""
echo "================================"
echo "✅ 初始化完成！"
echo "================================"
echo ""
echo "启动命令："
echo "  ./start.sh"
echo "  或"
echo "  ./task-skill --server"
echo ""
echo "访问地址：http://localhost:8080"
