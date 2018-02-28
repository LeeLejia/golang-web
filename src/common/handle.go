package common

import (
	"net/http"
)

/**
初始化路由
 */
func SetRouters(routers []BH){
	for i:=0; i< len(routers) ; i++{
		http.HandleFunc(routers[i].Url, routers[i].BsHandle)
	}
}

/**
处理基类
Url：   路由
Check： 校验用户合法性
Handle: 处理函数
 */
type BH struct{
	Url    string
	Check  bool
	Handle func(*Session,http.ResponseWriter, *http.Request)
}
/**
处理
 */
func (hd *BH) BsHandle(w http.ResponseWriter, r *http.Request){
	hd.SetContent(w,r)
	if !hd.Check{
		// 不需要校验Session
		hd.Handle(nil,w,r)
		return
	}
	if err:=CheckSession(w,r);err!=nil{
		ReturnEFormat(w, CODE_NEET_LOGIN_AGAIN, err.Error())
		return
	}
	// 需要校验Session
	sessionId:=""
	if cookie0,err0 := r.Cookie("sessionId");err0==nil{
		sessionId= cookie0.Value
	}else if sessionId = r.FormValue("sessionId");sessionId==""{
		sessionId = r.PostFormValue("sessionId")
	}
	session,_:=UserSessions.Load(sessionId)
	hd.Handle(session.(*Session),w,r)
}

/**
设置跨域访问
 */
func (hd *BH) SetContent(w http.ResponseWriter, r *http.Request) {
	origin:=r.Header.Get("Origin")
	w.Header().Set("Access-Control-Allow-Origin", origin)             						//允许访问所有域
	w.Header().Set("Access-Control-Allow-Credentials","true")
	w.Header().Add("Access-Control-Allow-Headers", "x-requested-with,content-type") 	//header的类型
}