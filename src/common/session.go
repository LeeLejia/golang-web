/**
 * 系统用户登录验证模块
 *
 * @author zhangqichao
 * Created on 2017-08-18
 */
package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Session struct {
	AccountID int64
	Account   string
	Token     string
}

var UserSession map[string]*Session

func GetSessionID(account, osType, role string) string {
	return fmt.Sprintf("%s-%s-%s", account, osType, role)
}

func SaveSession(uid int64, account, osType, role string) (token string) {
	if UserSession == nil {
		UserSession = map[string]*Session{}
	}
	sessionKey := GetSessionID(account, osType, role)
	session := new(Session)
	session.AccountID = uid
	session.Account = account

	crutime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime, 10))
	token = fmt.Sprintf("%x", h.Sum(nil))
	session.Token = token
	UserSession[sessionKey] = session
	return
}


func RemoveSession(r *http.Request) {
	account := r.FormValue("account")
	osType := r.FormValue("osType")
	role := r.FormValue("roleType")
	sessionKey := GetSessionID(account, osType, role)
	session := new(Session)
	session.AccountID = 0
	session.Account = ""
	session.Token = ""
	if UserSession == nil {
		UserSession = map[string]*Session{}
	}
	UserSession[sessionKey] = session
}

func CheckSession(r *http.Request) (account string, err error) {
	account = ""
	osType := ""
	role := ""
	token := ""
	if r.Method == "GET" {
		account = r.FormValue("account")
		osType = r.FormValue("osType")
		role = r.FormValue("roleType")
		token = r.FormValue("token")
	} else {
		account = r.PostFormValue("account")
		osType = r.PostFormValue("osType")
		role = r.PostFormValue("roleType")
		token = r.PostFormValue("token")
	}
	if account == "" {
		err = fmt.Errorf("用户账号校验失败，account为空！")
		return
	}
	if osType == "" {
		err = fmt.Errorf("用户session校验失败，osType为空！")
		return
	}

	sessionKey := GetSessionID(account, osType, role)
	session := UserSession[sessionKey]
	if session == nil {
		err = fmt.Errorf("用户session校验失败，session不存在！")
		return
	}
	if session.Account != account {
		err = fmt.Errorf("用户session校验失败，account不匹配，请重新登录！")
		return
	}
	if session.Token != token {
		err = fmt.Errorf("用户session校验失败，token不匹配，请重新登录！")
		return
	}
	return
}
