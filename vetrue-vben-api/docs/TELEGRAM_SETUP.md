# Telegram Bot 配置与使用指南

## 快速开始

### 1. 获取 Bot Token

1. 在 Telegram 中搜索并打开 `@BotFather`
2. 发送命令 `/newbot`
3. 按提示输入机器人显示名称（例如：`人人娱乐Bot`）
4. 输入机器人用户名（必须以 `bot` 结尾，例如：`RR666_bot`）
5. 创建成功后，BotFather 会返回 Bot Token，格式类似：
   ```
   123456789:ABCdefGHIjklMNOpqrsTUVwxyz1234567890
   ```

### 2. 配置 Bot

编辑 `configs/config.yaml` 文件，添加以下配置：

```yaml
telegram:
    enabled: true                                    # 启用 Telegram Bot
    botToken: "123456789:ABCdefGHIjklMNOpqrsTUVwxyz1234567890"  # 替换为你的 Bot Token
    botUsername: "RR666_bot"                        # 替换为你的 Bot 用户名（不含 @）
    miniAppUrl: ""                                  # 可选：小程序 URL
    setGlobalMenu: false                            # 是否设置全局菜单按钮
    debug: true                                     # 开发环境建议开启
```

### 3. 启动项目

```bash
go run main.go
```

或编译后运行：

```bash
go build -o vetrue-vben-api.exe .
./vetrue-vben-api.exe
```

### 4. 测试 Bot

1. 在 Telegram 中搜索你的 Bot 用户名（例如：`@RR666_bot`）
2. 点击 "Start" 或发送 `/start` 命令
3. 应该会收到欢迎消息和功能菜单

## 功能说明

### 主菜单 (/start)

发送 `/start` 命令后，Bot 会显示：

- **立即进入游戏** - 跳转到游戏网站
- **邀请好友** - 显示邀请链接和返佣信息
- **客服支持** - 联系人工客服
- **官方交流群** - 加入 Telegram 群组
- **官方频道** - 关注官方频道
- **语言切换** - 中文简体/English

### 邀请好友

点击 "邀请好友" 按钮后：

- 显示邀请返佣说明
- 提供推荐链接（可复制）
- 返回主菜单按钮

### 客服支持

点击 "客服支持" 按钮后：

- 显示客服联系方式
- 人工客服链接
- 返回主菜单按钮

## 自定义配置

### 修改欢迎消息

编辑 `internal/telegram/handlers/start_handler.go`：

```go
caption := `欢迎光临 人人娱乐综合游戏平台
🔥 活动：注册送体验金，支持 USDT、人民币等多币种兑换

✅ 免注册，免实名
✅ 支持USDT、人民币等便捷兑换
✅ 提供USDT、汇旺、ABA、借易等充值方式
✅ 5000万现金储备，100%出款

PG、PP、BB、电子、捕鱼、棋牌游戏、体育、真人，应有尽有`
```

### 修改链接地址

在各 handler 文件中修改 URL：

**游戏链接** (`start_handler.go`)：
```go
tgbotapi.NewInlineKeyboardButtonURL("立即进入游戏🎮", "https://your-domain.example/game")
```

**邀请链接** (`invite_handler.go`)：
```go
tgbotapi.NewInlineKeyboardButtonURL("复制推荐链接", "https://t.me/RR666_bot?start=129156")
```

**客服链接** (`support_handler.go`)：
```go
tgbotapi.NewInlineKeyboardButtonURL("人工客服", "https://t.me/rryl666")
```

### 修改图片

默认使用在线随机图片服务，可以修改为自己的图片 URL：

```go
// 原代码
photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL("https://picsum.photos/1200/675"))

// 修改为
photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL("https://your-domain.com/images/welcome.jpg"))

// 或使用本地文件
photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath("./static/images/welcome.jpg"))
```

## 常见问题

### Bot 无法启动

**问题**：运行项目后提示 "Telegram bot 未启用"

**解决**：检查 `configs/config.yaml` 中 `telegram.enabled` 是否为 `true`

---

**问题**：提示 "Telegram bot 配置不完整"

**解决**：确保 `botToken` 和 `botUsername` 都已正确配置

---

**问题**：提示 "401 Unauthorized"

**解决**：检查 Bot Token 是否正确，是否包含完整的 Token（包括冒号）

### Bot 无响应

**问题**：发送 `/start` 没有反应

**解决**：

1. 检查控制台日志，是否有错误信息
2. 确认 Bot Token 正确
3. 检查网络连接（某些地区可能需要代理）
4. 启用 debug 模式查看详细日志：`telegram.debug: true`

### 图片无法显示

**问题**：消息发送成功但图片不显示

**解决**：

1. 检查图片 URL 是否可访问
2. 图片大小不要超过 10MB
3. 使用支持的格式（JPG、PNG、WebP）

## Bot 命令列表设置（可选）

为了让用户更方便地使用命令，可以在 BotFather 中设置命令列表：

1. 向 `@BotFather` 发送 `/setcommands`
2. 选择你的 Bot
3. 发送以下命令列表：

```
start - 🎉 开始
menu - 📋 菜单
center - 👤 个人中心
channel - 📢 官方频道
help - 帮助 & FAQ
```

## 安全建议

1. **不要泄露 Bot Token**
   - Token 相当于密码，不要提交到公开代码仓库
   - 使用环境变量或配置文件管理
   - 定期更换 Token（在 BotFather 中使用 `/token` 命令）

2. **验证用户身份**
   - 在处理敏感操作前验证用户
   - 记录用户操作日志
   - 实现限流机制防止滥用

3. **数据安全**
   - 不要在 Telegram 消息中传输敏感信息
   - 使用 HTTPS 链接
   - 加密存储用户数据

## 下一步

- [ ] 集成数据库存储用户信息
- [ ] 实现 USDT 充值功能（Tron 模块）
- [ ] 添加用户认证和权限管理
- [ ] 实现游戏积分系统
- [ ] 添加管理员后台功能

更多详细信息请参考：`internal/telegram/README.md`
