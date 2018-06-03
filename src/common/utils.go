package common

import (
	"math/rand"
	"time"
	"strconv"
	"bytes"
)

/**
获取随机字符串
 */
func  GetRandomString(l int) string {
	str := "0123456789ASDFGHJKLPOIUYTREWQZXCVBNM"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
/**
获取随机数字
 */
func  GetRandomInt(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	result = append(result, bytes[r.Intn(len(bytes)-1)+1])
	for i := 1; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
/**
字符串转ints
 */
func BytesToInt(data []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data {
		s := strconv.FormatInt(int64(b&0xff), 10)
		buffer.WriteString(s)
	}
	return buffer.String()
}
