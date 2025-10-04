package core

import (
	"strings"
	"vetrue-vben-api/internal/telegram/enums"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// IsCommand 匹配 /cmd、/cmd 参数、或 /cmd@BotName
func IsCommand(update tgbotapi.Update, command enums.Command, botUsername string) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := update.Message.Text
	cmd := command.Value
	cmdLen := len(cmd)

	if !strings.HasPrefix(strings.ToLower(text), strings.ToLower(cmd)) {
		return false
	}

	// 完整 /cmd
	if len(text) == cmdLen {
		return true
	}

	if len(text) <= cmdLen {
		return false
	}

	next := text[cmdLen]

	// /cmd xxx
	if next == ' ' {
		return true
	}

	// /cmd@YourBotName
	if next == '@' && botUsername != "" {
		tail := text[cmdLen+1:]
		atName := strings.Split(tail, " ")[0]
		return strings.EqualFold(atName, botUsername)
	}

	return false
}

// IsDescription 是否描述匹配
func IsDescription(update tgbotapi.Update, command enums.Command) bool {
	if update.Message == nil || update.Message.Text == "" {
		return false
	}

	text := update.Message.Text
	return strings.EqualFold(text, command.Description)
}
