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
		if user,err:=CheckSession(w,r);err!=nil{
			ReturnEFormat(w, CODE_NEET_LOGIN_AGAIN, err.Error())
			return
		}else{
			hd.User=user
		}
		hd.Handle2(w,r,&hd.User)
	}else{
		hd.Handle(w,r)
	}

}

/**
设置跨域访问
 */
func SetContent(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")             //允许访问所有域
	w.Header().Set("Access-Control-Allow-Credentials","true");
	w.Header().Add("Access-Control-Allow-Headers", "x-requested-with,content-type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
}