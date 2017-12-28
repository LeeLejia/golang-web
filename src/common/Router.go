package common

import (
	"net/http"
	"../model"
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
	User   model.T_user
	Handle func(http.ResponseWriter, *http.Request)
	Handle2 func(http.ResponseWriter, *http.Request, *model.T_user)
}
/**
处理基方法
 */
func (hd *BH) BsHandle(w http.ResponseWriter, r *http.Request){
	SetContent(w)
	if hd.Check{
		if err:=CheckSession(w,r);err!=nil{
			ReturnEFormat(w, CODE_NEET_LOGIN_AGAIN, err.Error())
			return
		}
	}
	if(hd.Handle!=nil){
		hd.Handle(w,r)
	}else{
		sessionId:=""
		cookie0,err0 := r.Cookie("sessionId")
		if err0==nil{
			sessionId= cookie0.Value
		}
		if sessionId==""{
			sessionId = r.FormValue("sessionId")
		}
		if sessionId==""{
			sessionId = r.PostFormValue("sessionId")
		}
		session:=UserSessions[sessionId]
		if session!=nil{
			hd.User=session.User
		}
		hd.Handle2(w,r,&hd.User)
	}
}

/**
设置跨域访问
 */
func SetContent(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:88")             //允许访问所有域
	w.Header().Set("Access-Control-Allow-Credentials","true");
	w.Header().Add("Access-Control-Allow-Headers", "x-requested-with,content-type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
}