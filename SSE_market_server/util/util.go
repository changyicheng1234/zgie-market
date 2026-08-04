package util

import (
	"crypto/rand"
	"encoding/hex"
	mathRand "math/rand"
	"time"
)

// 工具函数，若注册时没有传名字，生成随机字符串
func RandomString(n int) string {
	var letters = []byte("fajvhfaufvafbiauAABIFWFIIudfuwfuwdagcqiuetoeh")
	result := make([]byte, n)
	mathRand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[mathRand.Intn(len(letters))]
	}
	return string(result)
}

// 生成随机数字
func GenerateRandomDigits(length int) string {
	digits := "0123456789"
	result := make([]byte, length)
	mathRand.Seed(time.Now().Unix())
	for i := 0; i < length; i++ {
		result[i] = digits[mathRand.Intn(len(digits))]
	}
	return string(result)
}

// 用于OAuth2授权过程生成随机字符串
func GenerateRandomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		// 如果crypto/rand失败，回退到math/rand
		r := mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
		for i := range bytes {
			bytes[i] = byte(r.Intn(256))
		}
	}
	return hex.EncodeToString(bytes)
}
