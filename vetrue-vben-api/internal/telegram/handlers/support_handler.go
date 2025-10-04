package handlers

import (
	"log"
	"vetrue-vben-api/internal/telegram/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SupportMenuHandler 客服支持页面：由 callbackData=support 触发
type SupportMenuHandler struct{}

// NewSupportMenuHandler 创建 SupportMenuHandler
func NewSupportMenuHandler() *SupportMenuHandler {
	return &SupportMenuHandler{}
}

// CanHandle 是否能处理该 update
func (h *SupportMenuHandler) CanHandle(update tgbotapi.Update, botUsername string) bool {
	return update.CallbackQuery != nil &&
		update.CallbackQuery.Data == core.Support
}

// Handle 业务处理
func (h *SupportMenuHandler) Handle(update tgbotapi.Update) {
	log.Println("SupportHandler.handle - 处理客服支持回调")
}

// BuildMessage 构建要发送的消息
func (h *SupportMenuHandler) BuildMessage(update tgbotapi.Update) *core.SendRequest {
	chatID := update.CallbackQuery.Message.Chat.ID
	msgID := update.CallbackQuery.Message.MessageID

	caption := `客服支持中心

请选择下方方式联系人工客服或加入官方社群。`

	keyboard := h.buildSupportKeyboard()

	// 编辑消息媒体和文本
	media := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL("https://picsum.photos/1000/560"))
	media.Caption = caption

	editConfig := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatID,
			MessageID:   msgID,
			ReplyMarkup: &keyboard,
		},
		Media: media,
	}

	return core.NewEditMessageMedia(editConfig)
}

// buildSupportKeyboard 构建客服键盘
func (h *SupportMenuHandler) buildSupportKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("人工客服", "https://t.me/rryl666"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("返回主菜单", core.BackToStart),
		),
	)
}
