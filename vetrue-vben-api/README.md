# Vetrue-Vben-API

<div align="center">
  <h3>基于 Go + Gin + GORM 构建的现代化后端 API 服务</h3>
  <p>为 Vben Admin 前端框架提供完整的后端接口支持</p>
</div>

## 📋 项目简介

Vetrue-Vben-API 是一个功能完善的企业级后端 API 服务，专为 [Vben Admin](https://github.com/vbenjs/vue-vben-admin) 前端框架设计。项目采用 Go 语言开发，使用 Gin 作为 Web 框架，GORM 作为 ORM 框架，提供了完整的权限管理、用户认证、系统管理等企业级功能。

## ✨ 特性

- 🚀 **高性能**：基于 Go 1.25 开发，原生支持高并发
- 🔐 **JWT 认证**：完整的 JWT token 认证机制
- 🔑 **RBAC 权限**：基于 Casbin 的角色访问控制，支持动态权限配置
- 📊 **多数据库支持**：支持 MySQL、PostgreSQL、SQLite、SQL Server
- 🎯 **RESTful API**：标准的 RESTful API 设计
- 🔧 **配置灵活**：支持 YAML 配置文件
- 🛡️ **优雅关闭**：支持服务优雅启动和关闭

## 🛠️ 技术栈

- **语言**：Go 1.25+
- **Web 框架**：Gin v1.10.1
- **ORM**：GORM v1.30.0
- **数据库**：MySQL / PostgreSQL / SQLite / SQL Server
- **认证**：JWT (golang-jwt/jwt v5.2.3)
- **权限**：Casbin v2.109.0 + GORM Adapter v3.33.0
- **配置管理**：Viper v1.20.1

## 📦 核心功能模块

### 系统管理
- **用户管理**：用户的增删改查、角色分配、状态管理
- **角色管理**：角色权限分配、菜单权限、API 权限
- **菜单管理**：动态菜单配置、菜单权限控制
- **API 管理**：API 接口权限管理、接口分组
- **操作日志**：用户操作记录、登录日志

### 基础功能
- **用户认证**：登录、登出、Token 刷新
- **个人中心**：个人信息修改、密码修改
- **权限控制**：基于 Casbin 的动态权限验证

## 🚀 快速开始

### 环境要求

- Go 1.25 或更高版本
- MySQL 5.7+ / PostgreSQL 12+ / SQLite 3 / SQL Server 2017+
- Git

### 安装步骤

1. **克隆项目**

   ```bash
   git clone git@github.com:ychengzhi96-wq/vetrue-admin.git
   cd vetrue-vben-api
   ```

2. **安装依赖**

   ```bash
   go mod tidy
   ```

3. **配置数据库**

   修改 `configs/config.yaml` 文件中的数据库配置：

   ```yaml
   database:
     driver: mysql  # 可选: mysql, postgres, sqlite, sqlserver
     host: localhost
     port: 3306
     username: root
     password: your_password
     database: vetrue_vben
     charset: utf8mb4
   ```

4. **初始化数据库**

   执行 SQL 脚本初始化数据库结构和默认数据：

   ```bash
   # 导入初始化脚本
   mysql -u root -p vetrue_vben < docs/sql/init.sql
   ```

5. **启动服务**

   ```bash
   # 直接运行
   go run main.go

   # 或编译后运行
   go build -o vetrue-vben-api
   ./vetrue-vben-api  # Linux/Mac
   vetrue-vben-api.exe  # Windows
   ```

   服务默认运行在 `http://localhost:8080`

## 📁 项目结构

```
vetrue-vben-api/
├── build/                    # 构建相关
│   ├── docker/              # Docker 配置文件
│   └── scripts/             # 构建和启动脚本
│
├── configs/                  # 配置文件
│   ├── config.yaml          # 主配置文件
│   └── rbac_model.conf      # RBAC 模型配置
│
├── docs/                     # 文档资源
│   └── sql/                 # 数据库脚本
│
├── internal/                 # 内部包（核心业务）
│   ├── config/              # 配置加载
│   ├── database/            # 数据库连接和初始化
│   ├── dto/                 # 数据传输对象
│   ├── handler/             # HTTP 处理器（控制器）
│   ├── logic/               # 业务逻辑层
│   ├── router/              # 路由定义
│   └── service/             # 服务层
│
├── models/                   # 数据模型
│   ├── model.go             # 基础模型
│   ├── sys_user.go          # 用户模型
│   ├── sys_role.go          # 角色模型
│   ├── sys_menu.go          # 菜单模型
│   ├── sys_api.go           # API 模型
│   ├── sys_dict.go          # 字典模型
│   ├── sys_record.go        # 操作记录模型
│   ├── sys_role_api.go      # 角色-API 关联
│   └── sys_role_auth.go     # 角色权限
│
├── pkg/                      # 外部包（可复用）
│   ├── middleware/          # 中间件
│   │   ├── auth.go          # 认证中间件
│   │   ├── casbin.go        # Casbin 权限中间件
│   │   ├── cors.go          # 跨域中间件
│   │   └── logger.go        # 日志中间件
│   ├── response/            # 统一响应格式
│   └── utils/               # 工具函数
│       ├── common.go        # 通用工具
│       └── jwt.go           # JWT 工具
│
├── static/                   # 静态资源
│   └── storage/             # 文件存储
│
├── .gitignore               # Git 忽略配置
├── go.mod                   # Go 模块定义
├── go.sum                   # 依赖版本锁定
├── main.go                  # 程序入口
└── README.md                # 项目说明文档
```

## 🔧 配置说明

### 主配置文件 (configs/config.yaml)

```yaml
app:
  name: "Vetrue Vben API"
  addr: ":8080"
  mode: "debug"  # debug, release, test
  timeout: 30    # 请求超时时间（秒）

database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: password
  database: vetrue_vben
  charset: utf8mb4
  max_idle_conn: 10
  max_open_conn: 100
  log_level: 4  # 1: Silent, 2: Error, 3: Warn, 4: Info

jwt:
  secret: "your-secret-key"
  expire: 7200  # Token 过期时间（秒）
  issuer: "vetrue-vben-api"

casbin:
  model_path: "configs/rbac_model.conf"
```

## 🔨 开发指南

### API 开发流程

1. **定义模型** - 在 `models/` 目录创建数据模型
2. **创建 DTO** - 在 `internal/dto/` 定义请求和响应结构
3. **实现逻辑** - 在 `internal/logic/` 编写业务逻辑
4. **添加处理器** - 在 `internal/handler/` 创建 HTTP 处理器
5. **配置路由** - 在 `internal/router/` 注册路由

### 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/logic

# 查看测试覆盖率
go test -cover ./...
```

### 编译构建

```bash
# 开发环境编译
go build -o vetrue-vben-api

# 生产环境编译（缩小体积）
go build -ldflags="-s -w" -o vetrue-vben-api

# 交叉编译 Linux 版本
GOOS=linux GOARCH=amd64 go build -o vetrue-vben-api-linux

# 交叉编译 Windows 版本
GOOS=windows GOARCH=amd64 go build -o vetrue-vben-api.exe
```

## 📚 API 接口说明

### 认证相关

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/login | 用户登录 |
| POST | /api/auth/logout | 用户登出 |
| POST | /api/auth/refresh | 刷新 Token |
| GET | /api/auth/profile | 获取当前用户信息 |

### 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/users | 获取用户列表 |
| GET | /api/users/:id | 获取用户详情 |
| POST | /api/users | 创建用户 |
| PUT | /api/users/:id | 更新用户 |
| DELETE | /api/users/:id | 删除用户 |

### 角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/roles | 获取角色列表 |
| GET | /api/roles/:id | 获取角色详情 |
| POST | /api/roles | 创建角色 |
| PUT | /api/roles/:id | 更新角色 |
| DELETE | /api/roles/:id | 删除角色 |

### 菜单管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/menus | 获取菜单列表 |
| GET | /api/menus/tree | 获取菜单树 |
| POST | /api/menus | 创建菜单 |
| PUT | /api/menus/:id | 更新菜单 |
| DELETE | /api/menus/:id | 删除菜单 |

## 🐳 Docker 部署

### 使用 Dockerfile

```bash
# 构建镜像
docker build -t vetrue-vben-api .

# 运行容器
docker run -d \
  --name vetrue-api \
  -p 8080:8080 \
  -v ./configs:/app/configs \
  vetrue-vben-api
```

### 使用 Docker Compose

```bash
# 启动服务
docker-compose up -d

# 停止服务
docker-compose down
```

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 📞 联系方式

- 作者：ychengzhi96-wq
- GitHub：[https://github.com/ychengzhi96-wq](https://github.com/ychengzhi96-wq)
- 项目地址：[https://github.com/ychengzhi96-wq/vetrue-admin](https://github.com/ychengzhi96-wq/vetrue-admin)

## 🙏 致谢

- [Gin Web Framework](https://github.com/gin-gonic/gin) - 高性能的 Go Web 框架
- [GORM](https://gorm.io/) - 功能强大的 Go ORM 库
- [Vben Admin](https://github.com/vbenjs/vue-vben-admin) - 现代化的 Vue 3 中后台管理框架
- [Casbin](https://casbin.org/) - 强大的访问控制库
- [Viper](https://github.com/spf13/viper) - Go 配置管理库

---

<div align="center">
  如果这个项目对你有帮助，请给一个 ⭐️ Star！
</div>
