# Telegram Bot 模块

这是一个独立的 Telegram Bot 模块，用于处理 Telegram 机器人相关的功能。

## 目录结构

```
internal/telegram/
├── core/                   # 核心功能
│   ├── bot_client.go      # Bot 客户端
│   ├── command_util.go    # 命令工具
│   ├── constants.go       # 常量定义
│   ├── handler.go         # Handler 接口
│   └── send_request.go    # 发送请求封装
├── enums/                 # 枚举类型
│   ├── command.go         # 命令枚举
│   └── send_method.go     # 发送方法枚举
├── handlers/              # 消息处理器
│   ├── start_handler.go   # 开始命令处理器
│   ├── invite_handler.go  # 邀请好友处理器
│   └── support_handler.go # 客服支持处理器
└── telegram.go            # 模块初始化入口
```

## 配置说明

在 `configs/config.yaml` 中配置 Telegram Bot：

```yaml
telegram:
    enabled: true                   # 是否启用 Telegram Bot
    botToken: "YOUR_BOT_TOKEN"     # Bot Token，从 @BotFather 获取
    botUsername: "your_bot"        # Bot 用户名（不含 @）
    miniAppUrl: ""                 # 小程序 URL（可选）
    setGlobalMenu: false           # 是否设置全局菜单
    debug: true                    # 是否开启调试模式
```

## 获取 Bot Token

1. 在 Telegram 中搜索 `@BotFather`
2. 发送 `/newbot` 创建新机器人
3. 按提示设置机器人名称和用户名
4. 获取 Bot Token 并填入配置文件

## 功能特性

### 已实现功能

- ✅ 长轮询模式接收消息
- ✅ 命令处理（/start 等）
- ✅ 内联键盘（InlineKeyboard）
- ✅ 回调查询（CallbackQuery）处理
- ✅ 多 Handler 架构
- ✅ 图片消息发送
- ✅ 文本消息发送

### 主要功能模块

1. **开始页面** (`/start`)
   - 欢迎消息
   - 游戏入口
   - 邀请好友/客服支持按钮
   - 官方群组/频道链接
   - 语言切换

2. **邀请好友页面** (`invite`)
   - 邀请链接生成
   - 返佣说明

3. **客服支持页面** (`support`)
   - 人工客服链接
   - 返回主菜单

## 使用方式

### 添加新的 Handler

1. 在 `handlers/` 目录下创建新的 handler 文件
2. 实现 `core.IHandler` 接口：

```go
package handlers

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "vetrue-vben-api/internal/telegram/core"
)

type MyHandler struct{}

func NewMyHandler() *MyHandler {
    return &MyHandler{}
}

// CanHandle 判断是否能处理该消息
func (h *MyHandler) CanHandle(update tgbotapi.Update, botUsername string) bool {
    // 实现判断逻辑
    return false
}

// Handle 业务处理（如数据库操作等）
func (h *MyHandler) Handle(update tgbotapi.Update) {
    // 实现业务逻辑
}

// BuildMessage 构建要发送的消息
func (h *MyHandler) BuildMessage(update tgbotapi.Update) *core.SendRequest {
    // 构建并返回消息
    return nil
}
```

3. 在 `telegram.go` 的 `InitTelegram()` 中注册：

```go
allHandlers := []core.IHandler{
    handlers.NewStartHandler(),
    handlers.NewMyHandler(),  // 添加新 handler
    // ...
}
```

### 添加新的命令

1. 在 `enums/command.go` 中添加命令定义：

```go
var (
    MYCOMMAND = Command{Value: "/mycommand", Description: "我的命令"}
)

// 更新 AllCommands
var AllCommands = []Command{START, MENU, MYCOMMAND, ...}
```

2. 创建对应的 handler 处理该命令

### 添加新的回调常量

在 `core/constants.go` 中添加：

```go
const (
    MyCallback = "my_callback"
)
```

## 扩展建议

### 后续可添加的功能模块

1. **用户管理模块**
   - 用户注册/登录
   - 用户信息管理
   - 用户级别/权限

2. **USDT 充值/提现模块**
   - Tron 钱包集成
   - USDT 充值检测
   - 提现请求处理

3. **游戏积分模块**
   - 积分查询
   - 积分兑换
   - 积分流水

4. **通知模块**
   - 系统通知推送
   - 活动通知
   - 个人消息

5. **统计模块**
   - 用户数据统计
   - 邀请数据统计
   - 交易数据统计

## 注意事项

1. **线程安全**：BotClient 在独立的 goroutine 中运行，注意数据竞争
2. **错误处理**：所有 handler 都有 panic 恢复机制，不会导致整个 bot 崩溃
3. **日志记录**：建议在 handler 中添加详细的日志，便于调试
4. **配置管理**：敏感信息（如 Bot Token）不要提交到代码仓库

## 依赖

- `github.com/go-telegram-bot-api/telegram-bot-api/v5` - Telegram Bot SDK

## 相关资源

- [Telegram Bot API 文档](https://core.telegram.org/bots/api)
- [telegram-bot-api SDK 文档](https://pkg.go.dev/github.com/go-telegram-bot-api/telegram-bot-api/v5)
- [BotFather 官方文档](https://core.telegram.org/bots#6-botfather)
