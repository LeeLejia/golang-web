/*
	通用类
 */
package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
)

type R struct {
	Code int         `json:"code"`
	Data map[string]interface{} `json:"data"`
}
type RE struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
/**
api回应
 */
func ReturnFormat(w http.ResponseWriter, code int, data map[string]interface{}) {
	if data!=nil{
		res := R{Code: code, Data: data}
		omg, _ := json.Marshal(res)
		w.Write(omg)
	}
}
/**
打印错误
 */
func ReturnEFormat(w http.ResponseWriter, code int, msg string) {
	res := RE{Code: code, Msg: msg}
	omg, _ := json.Marshal(res)
	//w.WriteHeader(code)
	w.Write(omg)
}

/**
获取md5
 */
func MD5Password(pwd string) string {
	data := []byte(pwd)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}

/**
获取页面限制条件
 */
func PageParams(r *http.Request) (limitStr string) {
	skip := 0
	limit := 20
	limitStr = ""
	page := r.FormValue("page")
	pageSize := r.FormValue("pageSize")
	pageInt, err := strconv.Atoi(page)
	if err == nil {
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err == nil {
			limit = pageSizeInt
			skip = limit * (pageInt - 1)
		}
	}
	limitStr = fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, skip)
	return
}

/**
匹配手机号
 */
func IsPhone(s string) bool {
	if regexp.MustCompile(`^(13[0-9]|14[0-9]|15[0-9]|17[0-9]|18[0-9])\d{8}$`).MatchString(s) {
		return true
	}
	return false
}

/**
匹配邮箱
 */
func IsEmail(s string) bool {
	if regexp.MustCompile("[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+").MatchString(s) {
		return true
	}
	return false
}

/**
是否包含某个角色
 */
func IsRole(userRoles string, role string) bool {
	roles := strings.Split(userRoles, ",")
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

/**
发送邮件通知
 */
func SendToMail() error {
	// Set up authentication information.
	userName:="baicaiplus@163.com"
	password:="baibaiasdfsssee2"
	host:="mail.163.com"
	auth := smtp.PlainAuth("", userName, password, host)
	to := []string{"1436983000@qq.com","cjwddz@gmail.com"}
	msg := []byte("To: recipient@example.net\r\n" +
		"Subject: discount Gophers!\r\n" +
		"\r\n" +
		"This is the email body.\r\n")
	err := smtp.SendMail("mail.163.com:25", auth, userName, to, msg)
	if err != nil {
		fmt.Printf(err.Error())
	}
	return err
}
