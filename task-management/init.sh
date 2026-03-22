#!/bin/bash

# Task Management Skill 初始化脚本
# 自动完成环境检查、依赖安装、编译和配置

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "🚀 Task Management Skill 初始化"
echo "================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Go 环境
echo "📋 检查环境..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ 错误：未找到 Go 环境${NC}"
    echo "请先安装 Go: https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓${NC} Go 环境：$GO_VERSION"

# 检查 Go 版本（需要 1.21+）
GO_MAJOR=$(go version | awk '{print $3}' | cut -d. -f1 | sed 's/go//')
GO_MINOR=$(go version | awk '{print $3}' | cut -d. -f2)

if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 21 ]); then
    echo -e "${YELLOW}⚠ 警告：Go 版本可能过旧，建议使用 1.21+${NC}"
fi

# 初始化 Go 模块
echo ""
echo "📦 初始化依赖..."
if [ ! -f "go.mod" ]; then
    echo "创建 go.mod..."
    go mod init task-skill
fi

echo "下载依赖..."
go mod download
echo -e "${GREEN}✓${NC} 依赖安装完成"

# 编译
echo ""
echo "🔨 编译程序..."
go build -o task-skill .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} 编译成功"
else
    echo -e "${RED}❌ 编译失败${NC}"
    exit 1
fi

# 初始化数据库
echo ""
echo "💾 初始化数据库..."
if [ -f "tasks.db" ]; then
    echo -e "${YELLOW}⚠ 数据库已存在，跳过初始化${NC}"
else
    echo -e "${GREEN}✓${NC} 数据库将在首次运行时自动创建"
fi

# 配置检查
echo ""
echo "⚙️  配置检查..."

# 检查端口
PORT=${TASK_SKILL_PORT:-8080}
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo -e "${YELLOW}⚠ 端口 $PORT 已被占用${NC}"
    echo "可通过环境变量修改端口：export TASK_SKILL_PORT=3000"
else
    echo -e "${GREEN}✓${NC} 端口 $PORT 可用"
fi

# 创建 systemd 服务文件（可选）
echo ""
read -p "是否创建 systemd 服务文件？(y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    if [ "$(id -u)" -ne 0 ]; then
        echo -e "${YELLOW}⚠ 需要 sudo 权限创建 systemd 服务${NC}"
        echo "请手动运行以下命令："
        echo "sudo $0 --create-systemd"
    else
        SERVICE_FILE="/etc/systemd/system/task-skill.service"
        cat > "$SERVICE_FILE" << EOF
[Unit]
Description=Task Management Skill
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$SCRIPT_DIR
Environment="TASK_SKILL_PORT=$PORT"
ExecStart=$SCRIPT_DIR/task-skill --server
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
        
        systemctl daemon-reload
        systemctl enable task-skill
        echo -e "${GREEN}✓${NC} systemd 服务已创建"
        echo "启动服务：sudo systemctl start task-skill"
        echo "查看状态：sudo systemctl status task-skill"
    fi
fi

# 完成
echo ""
echo "================================"
echo -e "${GREEN}✅ 初始化完成！${NC}"
echo ""
echo "启动方式："
echo "  1. 直接运行：./task-skill --server"
echo "  2. 使用脚本：./start.sh"
echo "  3. systemd 服务：sudo systemctl start task-skill"
echo ""
echo "访问地址：http://localhost:$PORT"
echo ""

# 询问是否立即启动
read -p "是否立即启动服务？(y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "🌐 启动服务..."
    ./task-skill --server
fi
