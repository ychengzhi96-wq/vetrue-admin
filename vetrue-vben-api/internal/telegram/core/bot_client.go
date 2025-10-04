package core

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotClient 机器人客户端
type BotClient struct {
	Bot           *tgbotapi.BotAPI
	Handlers      []IHandler
	BotUsername   string
	MiniAppURL    string
	SetGlobalMenu bool
}

// NewBotClient 创建新的 BotClient
func NewBotClient(botToken, botUsername, miniAppURL string, setGlobalMenu bool, handlers []IHandler) (*BotClient, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	return &BotClient{
		Bot:           bot,
		Handlers:      handlers,
		BotUsername:   botUsername,
		MiniAppURL:    miniAppURL,
		SetGlobalMenu: setGlobalMenu,
	}, nil
}

// SetDebugMode 设置调试模式
func (bc *BotClient) SetDebugMode(debug bool) {
	bc.Bot.Debug = debug
}

// InitGlobalMenu 设置全局菜单
func (bc *BotClient) InitGlobalMenu() error {
	if !bc.SetGlobalMenu || bc.MiniAppURL == "" {
		return nil
	}

	// 注意：telegram-bot-api v5 不直接支持设置菜单按钮
	// 如果需要设置菜单按钮，可能需要使用原始 API 调用
	// 这里暂时跳过，或者使用命令列表代替
	log.Println("全局菜单功能暂未实现（需要 Bot API 6.0+）")
	return nil
}

// DeleteWebhook 删除 Webhook（切换为长轮询模式）
func (bc *BotClient) DeleteWebhook() error {
	deleteWebhook := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: true,
	}
	_, err := bc.Bot.Request(deleteWebhook)
	if err != nil {
		log.Printf("删除 Webhook 失败（可忽略）: %v", err)
		return err
	}
	log.Println("已删除 Webhook，切换为长轮询模式")
	return nil
}

// Start 启动机器人，开始接收消息
func (bc *BotClient) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bc.Bot.GetUpdatesChan(u)

	log.Printf("主机器人注册成功 - Username: %s", bc.BotUsername)

	for update := range updates {
		bc.onUpdateReceived(update)
	}
}

// onUpdateReceived 处理接收到的更新
func (bc *BotClient) onUpdateReceived(update tgbotapi.Update) {
	// 查找能处理该消息的 handler
	for _, handler := range bc.Handlers {
		if bc.safeCanHandle(handler, update) {
			bc.processUpdateWithHandler(handler, update)
			return
		}
	}

	log.Println("未找到可处理的处理器，忽略该消息")
}

// safeCanHandle 安全检查是否可以处理
func (bc *BotClient) safeCanHandle(handler IHandler, update tgbotapi.Update) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("处理器 canHandle 调用异常: %v", r)
		}
	}()
	return handler.CanHandle(update, bc.BotUsername)
}

// processUpdateWithHandler 处理消息
func (bc *BotClient) processUpdateWithHandler(handler IHandler, update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("处理器执行异常: %v", r)
		}
	}()

	handler.Handle(update)
	msgRequest := handler.BuildMessage(update)
	if msgRequest != nil {
		bc.executeSafe(msgRequest)
	}
}

// executeSafe 安全执行发送
func (bc *BotClient) executeSafe(msgRequest *SendRequest) {
	if err := msgRequest.ExecuteBy(bc.Bot); err != nil {
		log.Printf("Telegram API error while executing: %v", err)
	}
}
