package test

import (
	"testing"
	"fmt"
	"../../common/conf"
	"time"
	m ".."
	"github.com/cjwddz/fast-model"
)

func TestT_vlog(t *testing.T) {
	conf.Init("/home/cjwddz/桌面/git-project/golang-web/src/app.toml")
	md, err := m.GetVLogModel()
	if err != nil {
		t.Fail()
	}
	// 正常插入
	vl:= m.T_vlog{
		Tag:"test",
		Code:"test14545415",
		App:"testold1234",
		Machine:"testsgeafgweafjpagewjqogjapwikpfw",
		Content:"test测试日志-长",
		CreatedAt:time.Now(),
	}
	err = md.Insert(vl)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	cond := model.DbCondition{}
	rs, err := md.Query(cond.And2("=", "tag", "test").Limit2(4, -1))
	if err != nil {
		fmt.Println("err:"+err.Error())
		t.Fail()
	}
	fmt.Println(rs)
	fmt.Println(md.CountAll())
	fmt.Println(md.Count(cond))
}