package handlers

import (
	"log"
	"vetrue-vben-api/internal/telegram/core"
	"vetrue-vben-api/internal/telegram/enums"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var globalBotClient *core.BotClient

// SetBotClient 设置全局 bot client（用于 handler 中访问）
func SetBotClient(bc *core.BotClient) {
	globalBotClient = bc
}

// StartHandler 开始 handler，接收 /start 命令处理
type StartHandler struct{}

// NewStartHandler 创建 StartHandler
func NewStartHandler() *StartHandler {
	return &StartHandler{}
}

// CanHandle 是否能处理该 update
func (h *StartHandler) CanHandle(update tgbotapi.Update, botUsername string) bool {
	// 支持 /start 命令与从其他页面返回主菜单的回调
	if core.IsCommand(update, enums.START, botUsername) || core.IsDescription(update, enums.START) {
		return true
	}
	return update.CallbackQuery != nil &&
		update.CallbackQuery.Data == core.BackToStart
}

// Handle 业务处理
func (h *StartHandler) Handle(update tgbotapi.Update) {
	log.Println("StartHandler.handle - 处理 /start 命令")
}

// BuildMessage 构建要发送的消息
func (h *StartHandler) BuildMessage(update tgbotapi.Update) *core.SendRequest {
	var chatID int64
	var messageID int

	if update.Message != nil {
		chatID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
		messageID = update.CallbackQuery.Message.MessageID
	}

	caption := `欢迎光临 人人娱乐综合游戏平台
🔥 活动：注册送体验金，支持 USDT、人民币等多币种兑换
PG、PP、BB、电子、捕鱼、棋牌游戏、体育、真人，应有尽有`

	keyboard := h.buildWelcomeKeyboard()

	// 如果是回调返回主菜单，编辑原消息以实现"替换"效果
	if update.CallbackQuery != nil {
		// 编辑消息媒体和文本
		media := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://picsum.photos/1200/675"))
		media.Caption = caption

		editConfig := tgbotapi.EditMessageMediaConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      chatID,
				MessageID:   messageID,
				ReplyMarkup: &keyboard,
			},
			Media: media,
		}
		return core.NewEditMessageMedia(editConfig)
	}

	// 直接使用默认在线图片 URL（无需配置）
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL("https://picsum.photos/1200/675"))
	photo.Caption = caption
	photo.ReplyMarkup = keyboard
	return core.NewSendPhoto(photo)
}

// buildWelcomeKeyboard 构建欢迎键盘
func (h *StartHandler) buildWelcomeKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("立即进入游戏🎮", "https://your-domain.example/game"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("邀请好友", core.Invite),
			tgbotapi.NewInlineKeyboardButtonData("客服支持", core.Support),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("官方交流群", "https://t.me/your_group_link"),
			tgbotapi.NewInlineKeyboardButtonURL("官方频道", "https://t.me/your_channel_link"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("中文简体", "lang_zh_cn"),
			tgbotapi.NewInlineKeyboardButtonData("English", "lang_en"),
		),
	)
}
