package telegram

import (
	"log"
	"sync"
	"vetrue-vben-api/internal/config"
	"vetrue-vben-api/internal/telegram/core"
	"vetrue-vben-api/internal/telegram/handlers"
)

var (
	botClient *core.BotClient
	once      sync.Once
)

// InitTelegram 初始化 Telegram Bot
func InitTelegram() error {
	var err error
	once.Do(func() {
		cfg := config.AppConfig.Telegram

		// 检查是否启用
		if !cfg.Enabled {
			log.Println("Telegram bot 未启用")
			return
		}

		// 检查必要配置
		if cfg.BotToken == "" || cfg.BotUsername == "" {
			log.Println("Telegram bot 配置不完整，跳过初始化")
			return
		}

		// 注册所有 handlers
		allHandlers := []core.IHandler{
			handlers.NewStartHandler(),
			handlers.NewInviteMenuHandler(),
			handlers.NewSupportMenuHandler(),
		}

		// 创建 BotClient
		botClient, err = core.NewBotClient(
			cfg.BotToken,
			cfg.BotUsername,
			cfg.MiniAppURL,
			cfg.SetGlobalMenu,
			allHandlers,
		)
		if err != nil {
			log.Printf("创建 Telegram Bot 失败: %v", err)
			return
		}

		// 设置调试模式
		botClient.SetDebugMode(cfg.Debug)

		// 设置全局 BotClient（供 handlers 使用）
		handlers.SetBotClient(botClient)

		// 删除 Webhook（切换为长轮询模式）
		if err = botClient.DeleteWebhook(); err != nil {
			log.Printf("删除 Webhook 失败: %v", err)
		}

		// 设置全局菜单
		if err = botClient.InitGlobalMenu(); err != nil {
			log.Printf("设置全局菜单失败: %v", err)
		}

		log.Println("Telegram bot 初始化成功")
	})

	return err
}

// StartTelegram 启动 Telegram Bot（阻塞）
func StartTelegram() {
	if botClient == nil {
		log.Println("Telegram bot 未初始化，无法启动")
		return
	}

	log.Println("正在启动 Telegram bot...")
	botClient.Start()
}

// GetBotClient 获取 BotClient 实例
func GetBotClient() *core.BotClient {
	return botClient
}
