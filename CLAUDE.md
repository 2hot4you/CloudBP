# CLAUDE.md
后续此项目中的所有回复跟文档输出、注释等都使用中文。另外需要你严格审视我的需求，如果有不合理的地方请你及时指出，当你完成某个需求时，你最好做一下测试。另外，当前主机中的 docker 容器比较多，不要搞混了。

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

这是一个多厂商云服务器销售平台，采用前后端分离架构：

- **后端**: Go 1.21 + Gin + PostgreSQL + Redis + RabbitMQ
- **前端**: Vue 3.3 + TypeScript + TDesign UI + Pinia + Vite
- **基础设施**: Docker + Nginx + SSL

## 核心功能

### 用户端 (C端)
- 用户认证 (登录/注册/OAuth)
- 服务器管理控制台
- 监控仪表板
- 购买流程
- 个人资料和费用中心
- 安全中心

### 管理端 (B端)
- 仪表板概览
- 用户管理
- 订单管理
- 多厂商产品管理
- 系统配置

### 多厂商支持
- 支持多个云厂商 (腾讯云、阿里云、AWS等)
- 跨厂商统一资源管理
- 产品对比和智能推荐系统
- 自动化厂商API集成

## 开发命令

### 快速启动
```bash
# 启动所有服务
./scripts/start.sh

# 停止所有服务
./scripts/stop.sh
```

### 后端开发
```bash
cd backend
go mod tidy                    # 安装依赖
go run cmd/api/main.go        # 启动开发服务器
go test ./...                 # 运行测试
```

### 前端开发
```bash
cd frontend
npm install                   # 安装依赖
npm run dev                   # 启动开发服务器
npm run build                 # 构建生产版本
npm run lint                  # 代码检查
```

### Docker操作
```bash
docker-compose up --build -d  # 构建并启动所有服务
docker-compose logs -f        # 查看日志
docker-compose down           # 停止服务
docker-compose down -v        # 停止服务并删除数据卷
```

## 架构说明

### 项目结构
```
cloudbp/
├── backend/           # Go后端服务
│   ├── cmd/api/      # 应用入口
│   ├── internal/     # 内部包
│   └── pkg/          # 公共包
├── frontend/         # Vue前端应用
├── nginx/            # 反向代理配置
├── scripts/          # 部署脚本
└── docker-compose.yml
```

### 设计模式
- **适配器模式**: 多厂商云API集成
- **仓储模式**: 数据访问层抽象
- **中间件模式**: 请求处理管道
- **工厂模式**: 云厂商实例创建

### 数据库
- **主数据库**: PostgreSQL 15 (支持JSON字段存储多厂商配置)
- **缓存**: Redis 7.0 (会话管理和缓存)
- **消息队列**: RabbitMQ 3.12 (异步任务处理)

## 服务地址

- 前端应用: http://localhost
- 后端API: http://localhost:8080
- API文档: http://localhost:8080/swagger/index.html
- RabbitMQ管理: http://localhost:15672
- 管理后台: https://localhost:8443

## 开发规范

### 后端开发
- 使用中文注释
- 遵循Go官方代码规范
- 错误处理使用统一的错误响应格式
- 使用Swagger注释生成API文档

### 前端开发
- 使用TypeScript类型定义
- 遵循Vue 3 Composition API
- 使用TDesign组件库
- API接口调用统一使用封装的request方法

### 数据库
- 使用GORM进行ORM操作
- 数据库迁移文件放在migrations目录
- 敏感信息不要硬编码在代码中