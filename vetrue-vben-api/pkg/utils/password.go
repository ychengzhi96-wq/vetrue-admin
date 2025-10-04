package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// ValidatePasswd 验证密码
func ValidatePasswd(inputPassword, salt, hashedPassword string) bool {
	return EncryptPassword(inputPassword, salt) == hashedPassword
}

// EncryptPassword 加密密码
func EncryptPassword(password, salt string) string {
	hash := md5.New()
	hash.Write([]byte(password + salt))
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateSalt 生成盐值
func GenerateSalt() string {
	// 可以使用更复杂的方式生成salt
	return "default_salt"
}
