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

func TestT_user(t *testing.T) {
	conf.Init("/home/cjwddz/桌面/git-project/golang-web/src/app.toml")
	md, err := m.GetUserModel()
	if err != nil {
		t.Fail()
	}
	// 正常插入
	v:=simplejson.New()
	v.Set("test","avalue")
	user:=m.T_user{
		Role:m.USER_ROLE_ADMIN,
		Nick:"test nick",
		Account:"test account",
		Pwd:"test xxxxx",
		Status:2,
		Avatar:"test https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
		Phone:"13480332034",
		Email:"cjwddz@qq.com",
		QQ:"1436983000",
		Expend: v,
		UpdatedAt:time.Now(),
		CreatedAt:time.Now(),
	}
	err = md.Insert(user)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	cond := model.DbCondition{}
	rs, err := md.Query(cond.And2("=", "nick", "test nick").Limit2(4, -1))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	sq:=model.DbSetCondition{}
	err = md.Update(sq.And2("=", "nick", "test nick").Set2("nick","test update nickkk"))
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Println(rs)
	fmt.Println(md.CountAll())
	fmt.Println(md.Count(cond))
}