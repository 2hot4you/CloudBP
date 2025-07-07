#!/bin/bash

# 项目停止脚本

set -e

echo "正在停止云服务器销售平台..."

# 停止所有服务
docker compose down

echo "服务已停止"
echo ""
echo "如需完全清理（包括数据卷），请运行："
echo "docker compose down -v"