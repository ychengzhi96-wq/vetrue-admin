package enums

// SendMethod 发送消息方法枚举
type SendMethod int

const (
	SendMessage SendMethod = iota + 1
	SendPhoto
	SendDocument
	SendVideo
	EditMessageMedia
)

// MethodName 获取方法名称
func (s SendMethod) MethodName() string {
	switch s {
	case SendMessage:
		return "sendMessage"
	case SendPhoto:
		return "sendPhoto"
	case SendDocument:
		return "sendDocument"
	case SendVideo:
		return "sendVideo"
	case EditMessageMedia:
		return "editMessageMedia"
	default:
		return ""
	}
}

// FromMethod 通过方法名获取 SendMethod
func FromMethod(method string) SendMethod {
	switch method {
	case "sendMessage":
		return SendMessage
	case "sendPhoto":
		return SendPhoto
	case "sendDocument":
		return SendDocument
	case "sendVideo":
		return SendVideo
	case "editMessageMedia":
		return EditMessageMedia
	default:
		return 0
	}
}

// FromValue 通过值获取 SendMethod
func FromValue(value int) SendMethod {
	if value >= 1 && value <= 5 {
		return SendMethod(value)
	}
	return 0
}
