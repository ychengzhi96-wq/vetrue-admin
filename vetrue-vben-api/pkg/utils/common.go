package utils

import (
	"math/rand"
	"time"
)

// IsValidNumber 验证数字是否有效
func IsValidNumber(num int64) bool {
	return num > 0
}

// RandString 生成随机字符串
func RandString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// MakePasswd 生成密码哈希
func MakePasswd(password, salt string) string {
	return EncryptPassword(password, salt)
}
