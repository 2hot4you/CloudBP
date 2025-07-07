# 云服务器销售平台 (CloudBP)

> 🚀 一个功能完整的多厂商云服务器销售平台，支持用户管理、产品销售、订单处理和服务器管理等核心业务功能。

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.3+-green.svg)](https://vuejs.org)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## ✨ 项目特色

- 🏗️ **微服务架构** - 前后端分离，模块化设计
- 🔐 **完整认证系统** - JWT认证，密码加密，权限控制  
- 🌐 **多厂商支持** - 抽象云厂商接口，支持腾讯云、阿里云等
- 📊 **实时监控** - 服务器状态监控，资源使用统计
- 💰 **订单管理** - 完整的购买流程，支付系统，余额管理
- 🛡️ **安全可靠** - 数据加密，SQL注入防护，XSS防护

## 🏛️ 技术架构

### 后端技术栈
- **语言**: Go 1.21 - 高性能并发处理
- **框架**: Gin - 轻量级Web框架
- **数据库**: PostgreSQL 15 - 关系型数据库
- **缓存**: Redis 7.0 - 高性能缓存
- **消息队列**: RabbitMQ 3.12 - 异步任务处理
- **ORM**: GORM - 强大的Go ORM
- **认证**: JWT - 无状态身份验证
- **配置**: Viper - 灵活的配置管理
- **日志**: Zap - 结构化高性能日志

### 前端技术栈
- **框架**: Vue 3.3 + TypeScript - 现代前端框架
- **UI组件**: TDesign Vue Next - 企业级设计语言
- **状态管理**: Pinia - Vue 3官方推荐
- **路由**: Vue Router 4 - 单页应用路由
- **HTTP**: Axios - Promise风格HTTP客户端
- **构建**: Vite - 极速前端构建工具

## 📁 项目结构

```
cloudbp/
├── backend/                    # Go后端服务
│   ├── cmd/
│   │   ├── api/               # API服务入口
│   │   └── migrate/           # 数据库迁移工具
│   ├── internal/              # 内部包
│   │   ├── config/            # 配置管理
│   │   ├── handler/           # HTTP处理器
│   │   │   ├── auth.go        # 认证处理器
│   │   │   ├── server.go      # 服务器管理处理器
│   │   │   ├── admin.go       # 管理员处理器
│   │   │   └── routes.go      # 路由配置
│   │   ├── middleware/        # 中间件
│   │   ├── model/             # 数据模型
│   │   │   ├── user.go        # 用户模型
│   │   │   ├── server.go      # 服务器模型
│   │   │   ├── order.go       # 订单模型
│   │   │   ├── provider.go    # 厂商模型
│   │   │   └── system.go      # 系统模型
│   │   └── service/           # 业务逻辑
│   │       ├── user.go        # 用户服务
│   │       ├── server.go      # 服务器服务
│   │       ├── provider.go    # 厂商服务
│   │       └── admin.go       # 管理员服务
│   ├── pkg/                   # 公共包
│   │   ├── auth/              # 认证工具
│   │   │   ├── jwt.go         # JWT管理
│   │   │   └── password.go    # 密码工具
│   │   ├── database/          # 数据库管理
│   │   │   ├── database.go    # 数据库连接
│   │   │   ├── migration.go   # 迁移管理器
│   │   │   └── init_data.go   # 种子数据
│   │   ├── provider/          # 云厂商集成
│   │   │   ├── interface.go   # 厂商接口
│   │   │   └── tencent.go     # 腾讯云适配器
│   │   ├── cache/             # 缓存
│   │   ├── logger/            # 日志
│   │   └── queue/             # 消息队列
│   ├── migrations/            # 数据库迁移文件
│   │   ├── 001_create_tables.sql  # 建表SQL
│   │   └── 002_seed_data.sql      # 种子数据
│   └── config.yaml            # 配置文件
├── frontend/                  # Vue前端应用
│   ├── src/
│   │   ├── components/        # 组件
│   │   ├── views/             # 页面
│   │   │   ├── auth/          # 认证页面
│   │   │   ├── user/          # 用户页面
│   │   │   └── admin/         # 管理页面
│   │   ├── router/            # 路由
│   │   ├── stores/            # 状态管理
│   │   ├── api/               # API接口
│   │   └── types/             # 类型定义
│   └── public/                # 静态资源
├── nginx/                     # 反向代理配置
│   ├── nginx.conf             # Nginx配置
│   └── generate-ssl.sh        # SSL证书生成
├── scripts/                   # 部署脚本
│   ├── start.sh               # 启动脚本
│   └── stop.sh                # 停止脚本
├── docker-compose.yml         # Docker编排
└── README.md                  # 项目文档
```

## 🎯 核心功能

### ✅ 已完成功能

#### 🔐 用户认证系统
- [x] 用户注册/登录
- [x] JWT Token生成、验证、刷新
- [x] 密码加密存储 (bcrypt)
- [x] 用户信息管理
- [x] 权限中间件 (用户/管理员)

#### 🖥️ 服务器管理
- [x] 产品列表查询
- [x] 服务器购买流程
- [x] 用户服务器列表
- [x] 服务器详情查看
- [x] 服务器操作 (启动/停止/重启)

#### 🏢 多厂商支持  
- [x] 抽象云厂商接口设计
- [x] 腾讯云Lighthouse适配器
- [x] 厂商配置管理
- [x] 产品信息同步

#### 💰 订单系统
- [x] 订单创建和管理
- [x] 支付流程 (余额支付)
- [x] 订单状态跟踪
- [x] 支付记录

#### 👥 管理后台
- [x] 仪表板统计数据
- [x] 用户管理
- [x] 订单管理  
- [x] 产品管理

#### 🗄️ 数据管理
- [x] 数据库自动迁移
- [x] 种子数据初始化
- [x] 迁移版本控制
- [x] 数据备份恢复

### 🚧 规划中功能

#### 💳 支付系统增强
- [ ] 微信支付集成
- [ ] 支付宝支付集成
- [ ] 银行卡支付
- [ ] 优惠券系统

#### 📊 监控系统
- [ ] 实时性能监控
- [ ] 资源使用统计
- [ ] 告警通知
- [ ] 日志分析

#### 🌐 更多云厂商
- [ ] 阿里云ECS集成
- [ ] AWS EC2集成
- [ ] 华为云集成

## 🚀 快速开始

### 环境要求
- 🐳 Docker & Docker Compose
- 🐹 Go 1.21+ (开发环境)
- 📦 Node.js 18+ (开发环境)

### 一键启动
```bash
# 克隆项目
git clone https://github.com/2hot4you/CloudBP.git
cd CloudBP

# 启动所有服务 (数据库、Redis、后端、前端)
./scripts/start.sh

# 停止所有服务
./scripts/stop.sh
```

### 分步启动

#### 1️⃣ 启动基础服务
```bash
docker compose up -d postgres redis rabbitmq
```

#### 2️⃣ 启动后端
```bash
cd backend
go mod tidy
go run cmd/api/main.go
```

#### 3️⃣ 启动前端
```bash
cd frontend
npm install
npm run dev
```

## 🌐 服务地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 🖥️ 前端应用 | http://localhost:3000 | Vue前端界面 |
| 🔌 后端API | http://localhost:8081 | RESTful API |
| 📚 API文档 | http://localhost:8081/swagger/index.html | Swagger文档 |
| ❤️ 健康检查 | http://localhost:8081/health | 服务状态 |
| 🐰 RabbitMQ | http://localhost:15672 | 消息队列管理 |
| 🗄️ PostgreSQL | localhost:5434 | 数据库连接 |
| 📦 Redis | localhost:6381 | 缓存连接 |

## 🧪 API测试

### 用户注册
```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com", 
    "password": "Test123456",
    "real_name": "测试用户"
  }'
```

### 用户登录
```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Test123456"
  }'
```

### 获取产品列表
```bash
# 需要先登录获取token
curl -X GET http://localhost:8081/api/v1/server/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## ⚙️ 配置说明

主要配置在 `backend/config.yaml`:

```yaml
server:
  port: "8081"              # 服务端口
  mode: "debug"             # 运行模式

database:
  host: "localhost"         # 数据库地址
  port: "5434"             # 数据库端口  
  user: "cloudbp"          # 数据库用户
  password: "cloudbp123"    # 数据库密码
  dbname: "cloudbp"        # 数据库名
  sslmode: "disable"       # SSL模式

redis:
  host: "localhost"         # Redis地址
  port: "6381"             # Redis端口
  password: ""             # Redis密码
  db: 0                    # Redis数据库

jwt:
  secret: "your-secret-key-change-in-production"  # JWT密钥
  expire_time: 86400       # Token过期时间(秒)

log:
  level: "info"            # 日志级别
```

## 🛠️ 开发指南

### 数据库迁移
```bash
# 查看迁移状态
go run cmd/migrate/main.go -action=status

# 执行迁移
go run cmd/migrate/main.go -action=up

# 创建新迁移
go run cmd/migrate/main.go -action=create -name=add_new_table
```

### 添加新的云厂商
1. 在 `backend/pkg/provider/` 创建新的适配器文件
2. 实现 `CloudProvider` 接口的所有方法
3. 在 `provider.go` 中注册新厂商
4. 在数据库中添加厂商配置

### 添加新的API接口
1. 在对应的handler文件中添加处理函数
2. 添加Swagger注释
3. 在 `routes.go` 中注册路由
4. 编写对应的service层逻辑

## 🐳 Docker部署

### 完整部署
```bash
# 构建并启动所有服务
docker compose up --build -d

# 查看服务状态
docker compose ps

# 查看日志
docker compose logs -f

# 停止服务
docker compose down

# 完全清理(包括数据)
docker compose down -v
```

### 生产环境配置
1. 🔐 修改默认密码和JWT密钥
2. 🛡️ 配置正式SSL证书
3. 🌍 设置环境变量
4. 📊 配置监控和日志
5. 💾 定期数据库备份

## 📊 项目统计

- **代码行数**: 7,800+ 行
- **文件数量**: 33+ 个文件
- **API接口**: 15+ 个接口
- **数据表**: 9张核心表
- **支持厂商**: 腾讯云 (更多厂商开发中)

## 🤝 贡献指南

1. 🍴 Fork 项目
2. 🌿 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 📝 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 📤 推送分支 (`git push origin feature/AmazingFeature`)
5. 🔀 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 🙏 致谢

- [Gin](https://gin-gonic.com/) - Go Web框架
- [Vue.js](https://vuejs.org/) - 渐进式JavaScript框架  
- [TDesign](https://tdesign.tencent.com/) - 企业级设计语言
- [GORM](https://gorm.io/) - Go ORM库
- [PostgreSQL](https://www.postgresql.org/) - 开源关系数据库

---

⭐ 如果这个项目对你有帮助，请给它一个星星！

📧 有问题或建议？欢迎提交 [Issue](https://github.com/2hot4you/CloudBP/issues)