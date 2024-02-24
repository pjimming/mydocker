package stringx

import (
	"encoding/json"
	"math/rand"
	"time"
)

// ToJsonString 把结构体转换为json string格式
func ToJsonString(v any) (string, error) {
	dataByte, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(dataByte), nil
}

// RandString 生成随机数
func RandString(n int) string {
	letterBytes := "1234567890"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[r.Intn(len(letterBytes))]
	}
	return string(b)
}
