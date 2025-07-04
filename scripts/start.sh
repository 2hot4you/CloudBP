#!/bin/bash

# 项目启动脚本

set -e

echo "正在启动云服务器销售平台..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "错误: Docker未运行，请先启动Docker"
    exit 1
fi

# 生成SSL证书
echo "生成SSL证书..."
./nginx/generate-ssl.sh

# 构建并启动所有服务
echo "构建并启动服务..."
docker-compose up --build -d

echo "等待服务启动..."
sleep 10

# 检查服务状态
echo "检查服务状态..."
docker-compose ps

echo ""
echo "服务启动完成！"
echo "前端地址: http://localhost"
echo "后端API: http://localhost:8080"
echo "API文档: http://localhost:8080/swagger/index.html"
echo "RabbitMQ管理: http://localhost:15672 (用户名: cloudbp, 密码: cloudbp123)"
echo ""
echo "使用 'docker-compose logs -f' 查看日志"
echo "使用 'docker-compose down' 停止服务"