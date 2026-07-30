package utils

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

// RandomString 生成指定长度的随机字符串（只含大小写字母和数字）
// 用途：生成随机文件名、验证码、临时 Token 等
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// Md5 计算字符串的 MD5 值
// 用途：缓存键、文件校验（⚠️ 不要用于密码存储）
func Md5(str string) string {
	data := []byte(str)
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}
