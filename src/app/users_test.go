package app

import (
	"net/url"
	"testing"
	"../common"
	"fmt"
)

func TestRegister(t *testing.T) {
	common.InitDb(t)
	for true {
		v := url.Values{}
		v.Add("role", "common")
		v.Add("pwd", "xxxxxxx")
		v.Add("phone", "13480387032")
		v.Add("email", "4849@qq.com")
		res,err:=common.MockHttp(v,"/register",Register)
		if err!=nil{
			t.Errorf("error:%s",err.Error())
		}
		fmt.Println(res)
		expected := `{"code":200,"msg":"success"}`
		if res != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				res, expected)
			break
		}
	}
}

func TestLogin(t *testing.T) {
	common.InitDb(t)
	// 用手机号登录
	v := url.Values{}
	v.Add("role", "common")
	v.Add("pwd", "xxxxxxx")
	v.Add("account", "13480335035")
	res,err:=common.MockHttp(v,"/login",Login)
	fmt.Println(res)
	fmt.Println(err)
	// 用邮箱登录
	v = url.Values{}
	v.Add("role", "common")
	v.Add("pwd", "xxxxxxx")
	v.Add("account", "123456@qq.com")
	res,err=common.MockHttp(v,"/login",Login)
	fmt.Println(res)
	fmt.Println(err)
}


