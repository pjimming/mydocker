package randx

import (
	"math/rand"
	"time"
)

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
