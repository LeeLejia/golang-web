package common

import (
	"math/rand"
	"time"
)

/**
获取随机字符串
 */
func  GetRandomString(l int) string {
	str := "0123456789qwertyuiopasdfghjklzxcvbnm"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
