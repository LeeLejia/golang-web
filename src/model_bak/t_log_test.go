package model

import (
	"testing"
	"../common/conf"
	"../pdb"
	"fmt"
	"time"
)

const ConfigPath ="/home/cjwddz/桌面/git-project/golang-web/src/app.toml"

func TestFindLogs(t *testing.T) {
	conf.Init(ConfigPath)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	// 空条件
	logs,err:= FindLogs("","","")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("no condition len:%d\n",len(logs))
	// 有条件
	logs,err = FindLogs("where id<5","limit 2","order by created_at")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("with condition len:%d\n",len(logs))
}

func TestLog_Insert(t *testing.T) {
	conf.Init(ConfigPath)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	item:=&T_log{
		Type:      "debug",
		Tag:       "测试用例",
		Operator:  "tester",
		Content:   "忽略我吧，我是测试用例！",
		CreatedAt: time.Now(),
	}
	err=item.Insert()
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
}

func TestCountLogs(t *testing.T) {
	conf.Init(ConfigPath)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	count,err:=CountLogs("")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("CountLogs with no condition count:%d\n",count)

	count,err=CountLogs("where type='debug'")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("CountLogs with condition count:%d\n",count)
}