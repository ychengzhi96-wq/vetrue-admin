package handlers

import (
	"log"
	"vetrue-vben-api/internal/telegram/core"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// InviteMenuHandler 邀请好友页面：由 callbackData=invite 触发
type InviteMenuHandler struct{}

// NewInviteMenuHandler 创建 InviteMenuHandler
func NewInviteMenuHandler() *InviteMenuHandler {
	return &InviteMenuHandler{}
}

// CanHandle 是否能处理该 update
func (h *InviteMenuHandler) CanHandle(update tgbotapi.Update, botUsername string) bool {
	return update.CallbackQuery != nil &&
		update.CallbackQuery.Data == core.Invite
}

// Handle 业务处理
func (h *InviteMenuHandler) Handle(update tgbotapi.Update) {
	log.Println("InviteHandler.handle - 处理邀请好友回调")
}

// BuildMessage 构建要发送的消息
func (h *InviteMenuHandler) BuildMessage(update tgbotapi.Update) *core.SendRequest {
	chatID := update.CallbackQuery.Message.Chat.ID
	msgID := update.CallbackQuery.Message.MessageID

	caption := `邀请你的朋友，获取返佣最高0.6%

拉好友进机器人，自动绑定上级`

	keyboard := h.buildInviteKeyboard()

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

// buildInviteKeyboard 构建邀请键盘
func (h *InviteMenuHandler) buildInviteKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("复制推荐链接", "https://t.me/RR666_bot?start=129156"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("返回主菜单", core.BackToStart),
		),
	)
}
