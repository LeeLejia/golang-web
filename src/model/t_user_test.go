package model

import (
	"testing"
	"../common/conf"
	"github.com/bitly/go-simplejson"
	"../pdb"
	"fmt"
)


func TestT_user_Insert(t *testing.T) {
	conf.Init(conf_path)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	// 正常插入
	user:=&T_user{
		Role:USER_ROLE_ADMIN,
		Nick:"白菜",
		Pwd:"imjia123",
		Status:2,
		Avatar:"https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
		Phone:"13480332034",
		Email:"cjwddz@qq.com",
		QQ:"1436983000",
		Expend: simplejson.New(),
	}
	fmt.Println(user)
	err=user.Insert()
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	// 插入空对象
	user=&T_user{
	}
	fmt.Println(user)
	err=user.Insert()
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
}

func TestFindUser(t *testing.T) {
	conf.Init(conf_path)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	// 空条件
	users,err:= FindUsers("","","")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("no condition len:%d\n",len(users))
	// 有条件
	users,err = FindUsers("where id>5","limit 4","order by updated_at")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("with condition len:%d\n",len(users))
}

func TestCountUsers2(t *testing.T) {
	conf.Init(conf_path)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	count,err:=CountUsers("")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("countUser with no condition count:%d\n",count)

	count,err=CountUsers("where role=0")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("countUser with condition count:%d\n",count)
}

func TestUpdateUsers(t *testing.T) {
	conf.Init(conf_path)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	count,err:=CountUsers("where role=1")
	if err!=nil{
		t.Fail()
	}
	fmt.Printf("count before:%d\n",count)
	err=UpdateUsers("set role=0","where role =1")
	if err!=nil{
		t.Fail()
	}
	count,err=CountUsers("where role=1")
	if err!=nil{
		t.Fail()
	}
	fmt.Printf("count after:%d\n",count)
}