package core

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// IHandler Handler 接口定义
type IHandler interface {
	// CanHandle 是否能处理该 update
	CanHandle(update tgbotapi.Update, botUsername string) bool

	// Handle 业务处理（如落库等）
	Handle(update tgbotapi.Update)

	// BuildMessage 构建要发送的消息
	BuildMessage(update tgbotapi.Update) *SendRequest
}
