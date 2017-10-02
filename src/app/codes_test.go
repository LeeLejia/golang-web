package app

import (
	"testing"
	"../common"
	"github.com/bitly/go-simplejson"
	"fmt"
	"net/url"
	"../model"
)

func TestAddCode(t *testing.T) {
	common.InitDb(t)
	bh := []common.BH{
		{
			Url:"xxxx",
			Check:true,
			User:model.T_user{
				Id:1,
				Role:model.USER_ROLE_COMMON,
				Nick:"白菜",
				Pwd:"xxxx",
				Status:2,
				Avatar:"https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
				Phone:"13480332034",
				Email:"cjwddz@qq.com",
				QQ:"1436983000",
				Expend: simplejson.New(),
			},
			Handle2:AddCode,
		},{
			Url:"xxxx",
			Check:true,
			User:model.T_user{
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
			},
			Handle2:AddCode,
		},{
			Url:"xxxx",
			Check:true,
			User:model.T_user{
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
			},
			Handle2:AddCode,
		},
	}
	v := url.Values{}
	for _,b:=range(bh){
		res,err:=common.MockHttp2(v,"xxxx",b.User,b.BsHandle)
		if err!=nil{
			t.Log(err.Error())
			t.Fail()
		}
		fmt.Println(res)
	}

}