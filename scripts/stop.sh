#!/bin/bash

# 项目停止脚本

set -e

echo "正在停止云服务器销售平台..."

# 停止所有服务
docker-compose down

echo "服务已停止"

# 询问是否删除数据卷
read -p "是否要删除数据卷？(y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "删除数据卷..."
    docker-compose down -v
    echo "数据卷已删除"
fi

echo "云服务器销售平台已完全停止"