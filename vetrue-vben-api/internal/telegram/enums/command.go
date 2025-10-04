package enums

import "strings"

// Command Telegram 命令枚举
type Command struct {
	Value       string
	Description string
}

var (
	START   = Command{Value: "/start", Description: "🎉 开始"}
	MENU    = Command{Value: "/menu", Description: "📋 菜单"}
	CENTER  = Command{Value: "/center", Description: "👤 个人中心"}
	CHANNEL = Command{Value: "/channel", Description: "📢 官方频道"}
	ABOUT   = Command{Value: "/about", Description: "bot.bot.start.cmd_about"}
	HELP    = Command{Value: "/help", Description: "帮助 & FAQ"}
)

// AllCommands 所有命令列表
var AllCommands = []Command{START, MENU, CENTER, CHANNEL, ABOUT, HELP}

// FromStringIgnoreCase 通过文本返回当前执行的命令
func FromStringIgnoreCase(text string) *Command {
	text = strings.ToLower(text)
	for _, cmd := range AllCommands {
		if strings.ToLower(cmd.Value) == text {
			return &cmd
		}
	}
	return nil
}

// FromDescription 通过描述返回当前执行的命令
func FromDescription(desc string) *Command {
	desc = strings.ToLower(desc)
	for _, cmd := range AllCommands {
		if strings.ToLower(cmd.Description) == desc {
			return &cmd
		}
	}
	return nil
}
