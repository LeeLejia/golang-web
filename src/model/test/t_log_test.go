package test

import (
	"testing"
	"fmt"
	"../../common/conf"
	"time"
	m ".."
	"github.com/cjwddz/fast-model"
)

func TestT_log(t *testing.T) {
	conf.Init("/home/cjwddz/桌面/git-project/golang-web/src/app.toml")
	md, err := m.GetLogModel()
	if err != nil {
		t.Fail()
	}
	// 正常插入
	item:=m.T_log{
		Type:      "debug",
		Tag:       "test 测试用例",
		Operator:  "tester",
		Content:   "test 忽略我吧，我是测试用例！",
		CreatedAt: time.Now(),
	}
	err = md.Insert(item)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	cond := model.DbCondition{}
	rs, err := md.Query(cond.And2("=", "tag", "test 测试用例").Limit2(4, -1))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	sq:=model.DbSetCondition{}
	err = md.Update(sq.And2("=", "tag", "test 测试用例").Set2("tag","test update"))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(rs)
	fmt.Println(md.CountAll())
	fmt.Println(md.Count(cond))
}