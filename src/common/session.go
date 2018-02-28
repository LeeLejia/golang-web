package common

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"../model"
	"time"
	"sync"
)

type Session struct {
	User   model.T_user
	Token     string
	Time   time.Time
}
const TOKEN_TIME_OUT = 60*10
var UserSessions = sync.Map{}
/**
获取Session键
 */
func GetSessionID(account, osType, role string) string {
	raw:=fmt.Sprintf("%s-%s-%s", account, osType, role)
	h := md5.New()
	io.WriteString(h,raw)
	return fmt.Sprintf("%X", h.Sum(nil))
}
/**
更新Time,Token并返回Token
 */
func (sess *Session)GetToken()string{
	sess.Time = time.Now()
	crutime := sess.Time.Unix()
	h := md5.New()
	io.WriteString(h, sess.User.Role+sess.User.Email+strconv.FormatInt(crutime, 10))
	sess.Token = fmt.Sprintf("%x", h.Sum(nil))
	return sess.Token
}
/**
更新时间
 */
func (sess *Session)RefreshTime(){
	sess.Time = time.Now()
}
/**
保存session
 */
func SaveSession(user model.T_user, osType string) (sessionKey string, sess *Session) {
	sessionKey = GetSessionID(user.Account,osType, user.Role)
	sess = new(Session)
	sess.User = user
	sess.GetToken()
	UserSessions.Store(sessionKey,sess)
	sess.RefreshTime()
	return
}
/**
移除session
 */
func RemoveSession(r *http.Request) {
	sessionId := r.FormValue("sessionId")
	UserSessions.Delete(sessionId)
}
/**
检测session合法性
包括cookies,get,post提交方式,检查超时
 */
func CheckSession(w http.ResponseWriter,r *http.Request) (err error) {
	sessionId:=""
	token:=""
	cookie0,err0 := r.Cookie("sessionId")
	if err0==nil{
		sessionId= cookie0.Value
	}
	cookie1,err1 := r.Cookie("token")
	if err1==nil{
		token= cookie1.Value
	}
	if sessionId!="" && token!=""{
		//使用cookies包含了登录信息
	}else if r.Method == "GET" {
		sessionId = r.FormValue("sessionId")
		token = r.FormValue("token")
	} else {
		sessionId = r.PostFormValue("sessionId")
		token = r.PostFormValue("token")
	}
	session,ok:=UserSessions.Load(sessionId)
	if !ok || session==nil{
		err = fmt.Errorf("用户session校验失败，session不存在！")
		return
	}
	sess:= session.(*Session)
	if sess.Token!=token{
		err = fmt.Errorf("用户session校验失败，token不匹配，请重新登录！")
		return
	}
	during:=time.Now().Unix()-sess.Time.Unix()
	if during>TOKEN_TIME_OUT{
		err = fmt.Errorf("用户session校验失败，登录超时，请重新登录！")
		return
	}
	// 更新登录时间
	sess.RefreshTime()
	http.SetCookie(w,&http.Cookie{Name:"sessionId",Value:sessionId,Path:"/"})
	http.SetCookie(w,&http.Cookie{Name:"token",Value:token,Path:"/"})
	return nil
}
