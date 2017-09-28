package app

import (
	"testing"
	"../common"
	"net/url"
	"fmt"
	"../model"
	"github.com/bitly/go-simplejson"
)

/**
添加App
 */
func TestAddApp(t *testing.T) {
	common.InitDb(t)
	bh := &common.BH{
		Url:"xxxx",
		Check:true,
		Handle2:AddApp,
	}
	v := url.Values{}
	v.Set("version","2.0")
	v.Set("describe","这是一个很棒的程序！")
	v.Set("name","棒棒")
	v.Set("file","xxxx")
	v.Set("src","xxxxx")

	user:=model.T_user{
		Id:1,
		Role:model.USER_ROLE_SUPER,
		Nick:"白菜",
		Pwd:"xxxx",
		Status:2,
		Avatar:"https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
		Phone:"13480332034",
		Email:"cjwddz@qq.com",
		QQ:"1436983000",
		Expend: simplejson.New(),
	}
	res,err:=common.MockHttp2(v,"xxxx",user,bh.BsHandle)
	if err!=nil{
		t.Log(err.Error())
		t.Fail()
	}
	fmt.Println(res)
}

/**
获取App列表
 */
func TestListApps(t *testing.T) {
	common.InitDb(t)
	users:=[]model.T_user{
		{
			Id:1111,
			Role:model.USER_ROLE_SUPER,
			Nick:"白菜",
			Pwd:"xxxx",
			Status:2,
			Avatar:"https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
			Phone:"13480332034",
			Email:"cjwddz@qq.com",
			QQ:"1436983000",
			Expend: simplejson.New(),
		},{
			Id:1,
			Role:model.USER_ROLE_DEVELOPER,
			Nick:"白菜",
			Pwd:"xxxx",
			Status:2,
			Avatar:"https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
			Phone:"13480332034",
			Email:"cjwddz@qq.com",
			QQ:"1436983000",
			Expend: simplejson.New(),
		},{
			Id:2,
			Role:model.USER_ROLE_DEVELOPER,
			Nick:"白菜",
			Pwd:"xxxx",
			Status:2,
			Avatar:"https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
			Phone:"13480332034",
			Email:"cjwddz@qq.com",
			QQ:"1436983000",
			Expend: simplejson.New(),
		},
	}
	bh := &common.BH{
		Url:"xxxx",
		Check:true,
		Handle2:ListApps,
	}
	v := url.Values{}
	for i,user:=range(users){
		res,err:=common.MockHttp2(v,"xxxx",user,bh.BsHandle)
		if err!=nil{
			t.Log(err.Error())
			t.Fail()
		}
		fmt.Println(fmt.Sprintf("第%d个样例：",i+1)+res)
	}
}


