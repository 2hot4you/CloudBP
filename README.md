# 云服务器销售平台

一个基于Go后端和Vue前端的多厂商云服务器销售平台。

## 项目架构

### 后端技术栈
- **语言**: Go 1.21
- **框架**: Gin
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7.0
- **消息队列**: RabbitMQ 3.12
- **ORM**: GORM
- **配置管理**: Viper
- **日志**: Zap
- **认证**: JWT

### 前端技术栈
- **框架**: Vue 3.3 + TypeScript
- **UI组件**: TDesign Vue Next
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **图表**: ECharts
- **构建工具**: Vite

## 项目结构

```
cloudbp/
├── backend/               # Go后端服务
│   ├── cmd/api/          # 应用入口
│   ├── internal/         # 内部包
│   │   ├── config/       # 配置管理
│   │   ├── handler/      # 请求处理器
│   │   ├── middleware/   # 中间件
│   │   ├── model/        # 数据模型
│   │   ├── repository/   # 数据层
│   │   └── service/      # 业务逻辑
│   ├── pkg/             # 公共包
│   │   ├── database/    # 数据库
│   │   ├── cache/       # 缓存
│   │   ├── logger/      # 日志
│   │   └── queue/       # 消息队列
│   └── migrations/      # 数据库迁移
├── frontend/            # Vue前端应用
│   ├── src/
│   │   ├── components/  # 组件
│   │   ├── views/       # 页面
│   │   ├── router/      # 路由
│   │   ├── stores/      # 状态管理
│   │   ├── api/         # API接口
│   │   └── types/       # 类型定义
│   └── public/          # 静态资源
├── nginx/               # 反向代理配置
├── docs/                # 项目文档
└── scripts/             # 部署脚本
```

## 功能特性

### 用户端功能
- 🔐 用户认证 (登录/注册/三方登录)
- 🖥️ 服务器管理控制台
- 📊 监控控制台
- 🛒 购买页面
- 👤 个人资料管理
- 💳 费用中心
- 🔒 安全中心

### 管理端功能
- 📊 总览仪表板
- 👥 用户管理
- 📦 订单管理
- 🏢 多厂商产品管理
- ⚙️ 系统配置

### 多厂商支持
- 🔄 支持多个云厂商 (腾讯云、阿里云、AWS等)
- 🔗 统一资源管理
- 📈 产品对比和智能推荐
- 🔧 自动化厂商API集成

## 快速开始

### 环境要求
- Docker & Docker Compose
- Go 1.21+ (开发环境)
- Node.js 18+ (开发环境)

### 一键启动
```bash
# 启动所有服务
./scripts/start.sh

# 停止所有服务
./scripts/stop.sh
```

### 开发环境

#### 后端开发
```bash
cd backend
go mod tidy
go run cmd/api/main.go
```

#### 前端开发
```bash
cd frontend
npm install
npm run dev
```

## 服务地址

- 前端应用: http://localhost
- 后端API: http://localhost:8080
- API文档: http://localhost:8080/swagger/index.html
- RabbitMQ管理: http://localhost:15672
- 管理后台: https://localhost:8443

## 环境变量

主要配置项在 `backend/config.yaml` 中，也可以通过环境变量覆盖：

```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=cloudbp
DB_PASSWORD=cloudbp123
DB_NAME=cloudbp

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRE_TIME=86400
```

## 部署说明

### Docker部署
```bash
# 构建并启动
docker-compose up --build -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 生产环境注意事项
1. 修改默认密码和密钥
2. 配置正式的SSL证书
3. 设置适当的环境变量
4. 配置监控和日志收集
5. 定期备份数据库

## 开发指南

### 添加新的云厂商
1. 在 `backend/pkg/provider/` 下创建新的厂商适配器
2. 实现 `CloudProvider` 接口
3. 在配置文件中注册新厂商
4. 更新前端产品选择组件

### API开发
1. 在 `backend/internal/handler/` 下添加新的处理器
2. 更新路由配置
3. 添加Swagger文档注释
4. 编写单元测试

### 前端开发
1. 使用TDesign组件库
2. 遵循TypeScript类型定义
3. 使用Pinia进行状态管理
4. 添加对应的API类型定义

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 推送到分支
5. 创建 Pull Request

## 许可证

MIT License