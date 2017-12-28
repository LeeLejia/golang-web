package model

import (
	"testing"
	"../common/conf"
	"fmt"
	"../pdb"
)
func TestCountFiles(t *testing.T) {
	conf.Init(ConfigPath)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	c,e:=CountFiles("")
	if e!=nil{
		t.Fail()
	}
	fmt.Println(fmt.Sprintf("无筛选条件：%d",c))
	c,e=CountFiles("WHERE owner=1")
	if e!=nil{
		t.Fail()
	}
	fmt.Println(fmt.Sprintf("筛选：%d",c))
}

func TestT_File_Insert(t *testing.T) {
	conf.Init(ConfigPath)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	files:=[]T_File{
		{
			FileKey:"xxxxxxxxxxxxxxx",
			Owner:1,
		},{
		},{
			FileKey:"xxxxxxxxxxxxxxx",
		},{
			Owner:1,
		},
	}
	describe:=[]string{
		"完整数据",
		"空数据",
		"缺陷数据",
		"缺陷数据",
	}
	for i,f:=range(files){
		err:=f.Insert()
		if err!=nil{
			t.Errorf("%s err:%s",describe[i],err.Error())
		}else{
			fmt.Println("pass:",describe[i])
		}
	}
}
