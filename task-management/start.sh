#!/bin/bash

echo "🚀 启动 Task Management Skill..."
echo ""

# 编译
echo "📦 编译中..."
go build -o task-skill .

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译成功"
echo ""

# 启动服务器
echo "🌐 启动 Web 服务器..."
echo "访问地址：http://localhost:8080"
echo "按 Ctrl+C 停止服务器"
echo ""

./task-skill --server
