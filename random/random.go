package random

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

// MD5 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// String 生成随机字符串
func String(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// Number 生成随机数字
func Number(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}

	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

// Username 随机生成一个 用户名
func Username() string {
	return String(8)
}
