# 项目模块化架构说明

## 项目结构

```
vetrue-vben-api/
├── internal/
│   ├── telegram/          # Telegram Bot 模块（独立模块）
│   │   ├── core/          # 核心功能
│   │   ├── enums/         # 枚举类型
│   │   ├── handlers/      # 消息处理器
│   │   ├── telegram.go    # 模块入口
│   │   └── README.md      # 模块文档
│   ├── tron/              # Tron USDT 模块（待实现）
│   ├── config/            # 配置管理
│   ├── database/          # 数据库连接
│   ├── handler/           # HTTP 处理器
│   ├── logic/             # 业务逻辑
│   └── router/            # 路由配置
├── configs/
│   └── config.yaml        # 配置文件
├── docs/
│   ├── MODULE_ARCHITECTURE.md  # 本文档
│   └── TELEGRAM_SETUP.md       # Telegram 配置指南
└── main.go                # 程序入口
```

## 模块化设计原则

### 1. 独立性

每个模块都是相对独立的，具有：

- **独立的目录结构**：每个模块在 `internal/` 下有自己的目录
- **独立的配置项**：在 `config.yaml` 中有专属的配置节
- **独立的初始化函数**：每个模块提供自己的 `Init()` 函数
- **独立的文档**：每个模块有自己的 README.md

### 2. 可插拔性

模块可以轻松地启用或禁用：

```yaml
# config.yaml
telegram:
    enabled: true    # 启用 Telegram 模块

tron:
    enabled: false   # 禁用 Tron 模块
```

### 3. 低耦合

- 模块之间通过接口交互
- 避免直接依赖具体实现
- 使用依赖注入

### 4. 高内聚

- 相关功能集中在同一模块
- 清晰的职责划分
- 便于维护和扩展

## 现有模块

### Telegram Bot 模块

**位置**：`internal/telegram/`

**功能**：
- Telegram 机器人消息接收和发送
- 命令处理
- 内联键盘交互
- 回调查询处理

**配置**：
```yaml
telegram:
    enabled: true
    botToken: ""
    botUsername: ""
    miniAppUrl: ""
    setGlobalMenu: false
    debug: true
```

**初始化**：
```go
// main.go
import "vetrue-vben-api/internal/telegram"

func main() {
    // 初始化 Telegram Bot
    if err := telegram.InitTelegram(); err != nil {
        log.Printf("初始化 Telegram Bot 失败: %v (继续运行)", err)
    }

    // 启动 Telegram Bot（在独立的 goroutine 中）
    go func() {
        telegram.StartTelegram()
    }()
}
```

**文档**：
- 模块文档：`internal/telegram/README.md`
- 配置指南：`docs/TELEGRAM_SETUP.md`

## 待实现模块

### Tron USDT 模块

**位置**：`internal/tron/`（待创建）

**计划功能**：
- Tron 钱包管理
- USDT-TRC20 充值检测
- USDT 提现
- 交易记录查询
- 地址管理

**建议配置**：
```yaml
tron:
    enabled: true
    network: "mainnet"           # mainnet 或 testnet
    apiKey: ""                   # TronGrid API Key
    contractAddress: ""          # USDT-TRC20 合约地址
    hotWallet: ""                # 热钱包地址
    coldWallet: ""               # 冷钱包地址
    confirmations: 19            # 确认数
    minDeposit: 10               # 最小充值金额
```

**建议结构**：
```
internal/tron/
├── core/
│   ├── client.go         # Tron 客户端
│   ├── wallet.go         # 钱包管理
│   └── transaction.go    # 交易处理
├── monitor/
│   ├── deposit.go        # 充值监控
│   └── withdraw.go       # 提现处理
├── tron.go               # 模块入口
└── README.md             # 模块文档
```

## 模块开发指南

### 创建新模块的步骤

#### 1. 创建目录结构

```bash
mkdir -p internal/mymodule/{core,handlers}
```

#### 2. 添加配置定义

在 `internal/config/config.go` 中添加：

```go
type Config struct {
    // ... 其他配置
    MyModule  MyModuleConfig  `mapstructure:"mymodule"`
}

type MyModuleConfig struct {
    Enabled bool   `mapstructure:"enabled"`
    // 其他配置项...
}
```

#### 3. 在配置文件中添加配置

在 `configs/config.yaml` 中添加：

```yaml
mymodule:
    enabled: true
    # 其他配置项...
```

#### 4. 实现模块入口

创建 `internal/mymodule/mymodule.go`：

```go
package mymodule

import (
    "log"
    "vetrue-vben-api/internal/config"
)

var moduleInstance *MyModule

type MyModule struct {
    // 模块字段
}

// InitMyModule 初始化模块
func InitMyModule() error {
    cfg := config.AppConfig.MyModule

    if !cfg.Enabled {
        log.Println("MyModule 未启用")
        return nil
    }

    moduleInstance = &MyModule{}
    // 初始化逻辑...

    log.Println("MyModule 初始化成功")
    return nil
}

// StartMyModule 启动模块（如果需要）
func StartMyModule() {
    if moduleInstance == nil {
        log.Println("MyModule 未初始化")
        return
    }

    // 启动逻辑...
}

// GetInstance 获取模块实例
func GetInstance() *MyModule {
    return moduleInstance
}
```

#### 5. 在 main.go 中集成

```go
import "vetrue-vben-api/internal/mymodule"

func main() {
    // ... 其他初始化

    // 初始化模块
    if err := mymodule.InitMyModule(); err != nil {
        log.Printf("初始化 MyModule 失败: %v", err)
    }

    // 如果需要在后台运行
    go func() {
        mymodule.StartMyModule()
    }()

    // ... 其他代码
}
```

#### 6. 编写文档

创建 `internal/mymodule/README.md`：

```markdown
# MyModule 模块

## 功能说明
...

## 配置说明
...

## 使用方式
...
```

### 模块间通信

#### 方式 1：直接调用（适合简单场景）

```go
// 在模块 A 中调用模块 B
import "vetrue-vben-api/internal/moduleb"

func DoSomething() {
    instance := moduleb.GetInstance()
    if instance != nil {
        instance.SomeMethod()
    }
}
```

#### 方式 2：事件总线（适合复杂场景）

可以考虑引入事件总线模式，实现模块间的松耦合通信。

#### 方式 3：共享服务（推荐）

在 `pkg/` 下创建共享服务：

```go
// pkg/services/user_service.go
package services

type UserService interface {
    GetUser(id int64) (*User, error)
    CreateUser(user *User) error
}
```

各模块实现或使用该接口。

## 最佳实践

### 1. 配置管理

- ✅ 所有配置集中在 `config.yaml`
- ✅ 敏感信息使用环境变量
- ✅ 提供配置验证
- ✅ 提供默认值

### 2. 错误处理

- ✅ 模块初始化失败不应导致整个应用崩溃
- ✅ 使用日志记录错误详情
- ✅ 返回有意义的错误信息

### 3. 并发安全

- ✅ 使用 `sync.Once` 确保单例
- ✅ 使用互斥锁保护共享资源
- ✅ 避免数据竞争

### 4. 测试

- ✅ 为每个模块编写单元测试
- ✅ 使用 mock 隔离依赖
- ✅ 编写集成测试

### 5. 日志

- ✅ 统一日志格式
- ✅ 合理的日志级别
- ✅ 包含足够的上下文信息

## 模块依赖图

```
main.go
  ├── config (配置管理)
  ├── database (数据库)
  ├── telegram (Telegram Bot)
  │   ├── core
  │   ├── handlers
  │   └── enums
  └── tron (Tron USDT，待实现)
      ├── core
      └── monitor
```

## 扩展建议

未来可以添加的模块：

1. **用户认证模块** (`internal/auth/`)
   - JWT 认证
   - OAuth2 集成
   - 权限管理

2. **支付模块** (`internal/payment/`)
   - 多种支付方式
   - 订单管理
   - 退款处理

3. **通知模块** (`internal/notification/`)
   - 邮件通知
   - 短信通知
   - 推送通知

4. **统计分析模块** (`internal/analytics/`)
   - 数据统计
   - 报表生成
   - 数据导出

5. **缓存模块** (`internal/cache/`)
   - Redis 缓存
   - 内存缓存
   - 缓存策略

## 总结

通过模块化设计，项目具有以下优势：

- ✅ **易于维护**：每个模块职责清晰
- ✅ **易于扩展**：添加新功能不影响现有代码
- ✅ **易于测试**：模块间低耦合，便于单元测试
- ✅ **易于协作**：团队成员可以并行开发不同模块
- ✅ **易于部署**：可以按需启用或禁用模块

随着项目的发展，继续保持模块化的设计思路，将使项目更加健壮和可维护。
