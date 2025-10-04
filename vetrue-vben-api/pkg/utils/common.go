package utils

import (
	"math/rand/v2"
)

// IsValidNumber 验证数字是否有效
func IsValidNumber(num int64) bool {
	return num > 0
}

// RandString 生成随机字符串
// 使用 Go 1.25 的 math/rand/v2，不需要手动设置种子
func RandString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

// MakePasswd 生成密码哈希
func MakePasswd(password, salt string) string {
	return EncryptPassword(password, salt)
}
