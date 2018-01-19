package test

import (
	"testing"
	"fmt"
	"../../common/conf"
	"github.com/bitly/go-simplejson"
	"time"
	m ".."
	"github.com/cjwddz/fast-model"
)

func TestT_code(t *testing.T) {
	conf.Init("/home/cjwddz/桌面/git-project/golang-web/src/app.toml")
	md, err := m.GetCodeModel()
	if err != nil {
		t.Fail()
	}
	// 正常插入
	code:=m.T_code{
		Code:"testCode....",
		AppId:"sgwj2000jojo",
		Developer:125,
		Consumer:simplejson.New(),
		Describe:"测试",
		Valid:true,
		MachineCount:5,
		EnableTime:false,
		StartTime:time.Now(),
		EndTime:time.Now(),
		CreatedAt:time.Now(),
	}
	err = md.Insert(code)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	cond := model.DbCondition{}
	rs, err := md.Query(cond.And2("=", "valid", true).Limit2(4, -1))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	sq:=model.DbSetCondition{}
	err = md.Update(sq.And2("=", "app_id", "sgwj2000jojo").Set2("app_id","resetAppId"))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(rs)
	fmt.Println(md.CountAll())
	fmt.Println(md.Count(cond))
}