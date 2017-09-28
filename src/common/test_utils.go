package common

import (
	"net/url"
	"strings"
	"net/http/httptest"
	"fmt"
	"net/http"
	"errors"
	"testing"
	"./conf"
	"../pdb"
	"../model"
	"time"
)
/**
Http模拟(无验证)
 */
func MockHttp(v url.Values, route string,handle func(http.ResponseWriter,*http.Request))(res string, err error){
	req, err := http.NewRequest("POST", route, strings.NewReader(v.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		err = errors.New(fmt.Sprintf("status code: %d",status))
	}
	res = rr.Body.String()
	return
}
/**
Http模拟(过验证)
 */
func MockHttp2(v url.Values, route string,user model.T_user, handle func(http.ResponseWriter,*http.Request))(res string, err error){
	UserSessions=make(map[string]*Session)
	UserSessions["TEST_SESSION"]=&Session{
		User: user,
		Token:"TEST_TOKEN",
		Time:time.Now(),
	}
	req, err := http.NewRequest("POST", route, strings.NewReader(v.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	req.AddCookie(&http.Cookie{
		Name:"sessionId",
		Value:"TEST_SESSION",
	})
	req.AddCookie(&http.Cookie{
		Name:"token",
		Value:"TEST_TOKEN",
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handle)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		err = errors.New(fmt.Sprintf("status code: %d",status))
	}
	res = rr.Body.String()
	return
}

/**
数据库配置
 */
const conf_path="C:\\Users\\lejia\\Desktop\\git-project\\go-web\\app.toml"
func InitDb(t *testing.T){
	conf.Init(conf_path)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
}
