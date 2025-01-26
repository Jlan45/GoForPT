package tools

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func GenerateToken() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
func MD5(data []byte) string {
	hashValue := md5.Sum(data)

	// 将哈希值转换为十六进制字符串
	return hex.EncodeToString(hashValue[:])
}
func SHA1(data []byte) [20]byte {
	hashValue := sha1.Sum(data)

	// 将哈希值转换为十六进制字符串
	return hashValue
}
func RemoveEmptyAndDuplicates(data []string) []string {
	// 创建一个空的字符串切片
	var result []string
	for _, str := range data {
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}
