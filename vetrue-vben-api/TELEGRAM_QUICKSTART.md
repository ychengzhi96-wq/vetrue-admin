# Telegram Bot 快速开始指南

## 📦 已完成的工作

✅ 集成 Telegram Bot SDK
✅ 创建模块化架构（`internal/telegram/`）
✅ 实现核心功能（消息接收、发送、回调处理）
✅ 实现 3 个基础 Handler（Start、Invite、Support）
✅ 配置管理（`configs/config.yaml`）
✅ 主程序集成（`main.go`）

## 🚀 快速开始（3 步）

### 第 1 步：获取 Bot Token

1. 在 Telegram 搜索 `@BotFather`
2. 发送 `/newbot` 创建机器人
3. 按提示设置名称和用户名（必须以 `bot` 结尾）
4. 复制获得的 Token（格式：`123456789:ABCdefGHI...`）

### 第 2 步：配置 Bot

编辑 `configs/config.yaml`，找到 `telegram` 配置项：

```yaml
telegram:
    enabled: true                                    # 改为 true 启用
    botToken: "YOUR_BOT_TOKEN_HERE"                 # 粘贴你的 Token
    botUsername: "your_bot_username"                # 填写用户名（不含 @）
    miniAppUrl: ""                                  # 可选
    setGlobalMenu: false                            # 可选
    debug: true                                     # 建议开启
```

### 第 3 步：启动并测试

```bash
# 编译运行
go build -o vetrue-vben-api.exe .
./vetrue-vben-api.exe

# 或直接运行
go run main.go
```

在 Telegram 中搜索你的机器人并发送 `/start`，你应该会收到欢迎消息！

## 📁 项目结构

```
internal/telegram/          # Telegram Bot 模块
├── core/                   # 核心功能
│   ├── bot_client.go      # Bot 客户端（主控制器）
│   ├── command_util.go    # 命令解析工具
│   ├── constants.go       # 回调常量定义
│   ├── handler.go         # Handler 接口定义
│   └── send_request.go    # 消息发送封装
├── enums/                 # 枚举类型
│   ├── command.go         # 命令枚举（/start、/help 等）
│   └── send_method.go     # 发送方法枚举
├── handlers/              # 消息处理器
│   ├── start_handler.go   # 处理 /start 命令和主菜单
│   ├── invite_handler.go  # 处理邀请好友功能
│   └── support_handler.go # 处理客服支持功能
├── telegram.go            # 模块入口（初始化和启动）
└── README.md              # 详细文档
```

## 🎯 功能说明

### 主菜单（/start）
- 欢迎消息 + 图片
- 游戏入口按钮
- 邀请好友 / 客服支持
- 官方群组 / 频道链接
- 语言切换

### 邀请好友
- 显示返佣说明
- 提供推荐链接
- 返回主菜单

### 客服支持
- 人工客服链接
- 返回主菜单

## 🛠️ 自定义修改

### 修改欢迎文案

编辑 `internal/telegram/handlers/start_handler.go` 第 47-57 行：

```go
caption := `欢迎光临 人人娱乐综合游戏平台
🔥 活动：注册送体验金，支持 USDT、人民币等多币种兑换
...`
```

### 修改链接地址

**游戏链接** - `start_handler.go` 第 88 行：
```go
tgbotapi.NewInlineKeyboardButtonURL("立即进入游戏🎮", "https://your-domain.example/game")
```

**邀请链接** - `invite_handler.go` 第 56 行：
```go
tgbotapi.NewInlineKeyboardButtonURL("复制推荐链接", "https://t.me/RR666_bot?start=129156")
```

**客服链接** - `support_handler.go` 第 56 行：
```go
tgbotapi.NewInlineKeyboardButtonURL("人工客服", "https://t.me/rryl666")
```

### 修改图片

默认使用 `picsum.photos` 随机图片服务，可改为自己的图片：

```go
// 使用在线图片
photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL("https://your-cdn.com/image.jpg"))

// 使用本地图片
photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("./static/images/welcome.jpg"))
```

## 📝 添加新功能

### 示例：添加"个人中心"功能

#### 1. 创建 Handler

创建 `internal/telegram/handlers/center_handler.go`：

```go
package handlers

import (
    "log"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "vetrue-vben-api/internal/telegram/core"
    "vetrue-vben-api/internal/telegram/enums"
)

type CenterHandler struct{}

func NewCenterHandler() *CenterHandler {
    return &CenterHandler{}
}

func (h *CenterHandler) CanHandle(update tgbotapi.Update, botUsername string) bool {
    return core.IsCommand(update, enums.CENTER, botUsername)
}

func (h *CenterHandler) Handle(update tgbotapi.Update) {
    log.Println("CenterHandler - 处理个人中心")
}

func (h *CenterHandler) BuildMessage(update tgbotapi.Update) *core.SendRequest {
    chatID := update.Message.Chat.ID

    msg := tgbotapi.NewMessage(chatID, "个人中心功能开发中...")
    return core.NewSendMessage(msg)
}
```

#### 2. 注册 Handler

编辑 `internal/telegram/telegram.go`，在 `InitTelegram()` 函数中添加：

```go
allHandlers := []core.IHandler{
    handlers.NewStartHandler(),
    handlers.NewCenterHandler(),  // 添加这行
    handlers.NewInviteMenuHandler(),
    handlers.NewSupportMenuHandler(),
}
```

#### 3. 测试

重新编译运行，发送 `/center` 命令测试。

## 🔍 调试技巧

### 查看日志

启用 debug 模式后，会输出详细日志：

```yaml
telegram:
    debug: true
```

日志会显示：
- 接收到的消息
- API 请求详情
- Handler 处理流程

### 常见问题

**问题 1**：Bot 无响应
- 检查 Token 是否正确
- 检查网络连接
- 查看控制台错误信息

**问题 2**：图片不显示
- 检查图片 URL 是否可访问
- 图片大小不超过 10MB
- 使用支持的格式（JPG、PNG）

**问题 3**：按钮点击无反应
- 检查 callback data 是否正确
- 查看是否有对应的 Handler
- 查看日志错误信息

## 📚 更多文档

- **模块详细文档**：`internal/telegram/README.md`
- **配置指南**：`docs/TELEGRAM_SETUP.md`
- **架构说明**：`docs/MODULE_ARCHITECTURE.md`

## 🎉 下一步

现在你已经有了一个可运行的 Telegram Bot！接下来可以：

1. ✅ 集成数据库，存储用户信息
2. ✅ 实现 USDT 充值功能（Tron 模块）
3. ✅ 添加用户认证和权限管理
4. ✅ 实现更多游戏相关功能
5. ✅ 添加管理员后台

## 💡 技术特点

- ✅ **模块化设计**：独立的 telegram 模块，易于维护
- ✅ **可插拔**：通过配置轻松启用/禁用
- ✅ **低耦合**：Handler 模式，易于扩展
- ✅ **并发安全**：使用 goroutine，支持高并发
- ✅ **错误恢复**：内置 panic 恢复机制
- ✅ **日志完善**：详细的日志记录

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**祝你开发愉快！** 🚀
