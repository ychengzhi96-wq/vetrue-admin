package core

import (
	"log"
	"vetrue-vben-api/internal/telegram/enums"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendRequest 发送请求封装
type SendRequest struct {
	Method   enums.SendMethod
	Message  *tgbotapi.MessageConfig
	Photo    *tgbotapi.PhotoConfig
	Document *tgbotapi.DocumentConfig
	Video    *tgbotapi.VideoConfig
	Edit     *tgbotapi.EditMessageMediaConfig
}

// NewSendMessage 创建文本消息请求
func NewSendMessage(msg tgbotapi.MessageConfig) *SendRequest {
	return &SendRequest{
		Method:  enums.SendMessage,
		Message: &msg,
	}
}

// NewSendPhoto 创建图片消息请求
func NewSendPhoto(photo tgbotapi.PhotoConfig) *SendRequest {
	return &SendRequest{
		Method: enums.SendPhoto,
		Photo:  &photo,
	}
}

// NewSendDocument 创建文档消息请求
func NewSendDocument(doc tgbotapi.DocumentConfig) *SendRequest {
	return &SendRequest{
		Method:   enums.SendDocument,
		Document: &doc,
	}
}

// NewSendVideo 创建视频消息请求
func NewSendVideo(video tgbotapi.VideoConfig) *SendRequest {
	return &SendRequest{
		Method: enums.SendVideo,
		Video:  &video,
	}
}

// NewEditMessageMedia 创建编辑媒体消息请求
func NewEditMessageMedia(edit tgbotapi.EditMessageMediaConfig) *SendRequest {
	return &SendRequest{
		Method: enums.EditMessageMedia,
		Edit:   &edit,
	}
}

// ExecuteBy 通过指定的 bot 执行发送
func (sr *SendRequest) ExecuteBy(bot *tgbotapi.BotAPI) error {
	if sr == nil {
		return nil
	}

	var err error
	switch sr.Method {
	case enums.SendMessage:
		if sr.Message != nil {
			_, err = bot.Send(*sr.Message)
		}
	case enums.SendPhoto:
		if sr.Photo != nil {
			_, err = bot.Send(*sr.Photo)
		}
	case enums.SendDocument:
		if sr.Document != nil {
			_, err = bot.Send(*sr.Document)
		}
	case enums.SendVideo:
		if sr.Video != nil {
			_, err = bot.Send(*sr.Video)
		}
	case enums.EditMessageMedia:
		if sr.Edit != nil {
			_, err = bot.Request(*sr.Edit)
		}
	default:
		log.Printf("未知的发送方法: %v", sr.Method)
	}

	return err
}
