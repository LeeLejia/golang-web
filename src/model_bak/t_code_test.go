package model

import (
	"testing"
	"fmt"
	"../common/conf"
	"../pdb"
	"time"
	"github.com/bitly/go-simplejson"
)

func TestT_code_Insert(t *testing.T) {
	conf.Init(ConfigPath)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	// 正常插入
	// simplejson.NewJson([]byte("{'name':'stupy,'home':'www.baidu.com''}"))
	code:=&T_code{
		Code:"ssgo45",
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
	fmt.Println(code)
	err=code.Insert()
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	// 插入空对象
	code=&T_code{
	}
	fmt.Println(code)
	err=code.Insert()
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
}

func TestFindCodes(t *testing.T) {
	conf.Init(ConfigPath)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	// 空条件
	codes,err:= FindCodes("","","")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	if len(codes)>0{
		fmt.Println(codes[0])
	}
	fmt.Printf("findCodes with no condition len:%d\n",len(codes))
	// 有条件
	codes,err = FindCodes("where id>5","limit 4","order by created_at")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("findCodes with with condition len:%d\n",len(codes))
}

func TestCountCodes(t *testing.T) {
	conf.Init(ConfigPath)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	count,err:=CountCodes("")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("countCodes with no condition count:%d\n",count)

	count,err=CountCodes("where id>1")
	if err!=nil{
		fmt.Println(err.Error())
		t.Fail()
	}
	fmt.Printf("countCodes with condition count:%d\n",count)
}

func TestUpdateCodes(t *testing.T) {
	conf.Init(ConfigPath)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	count,err:=CountCodes("where developer=1")
	if err!=nil{
		t.Fail()
	}
	fmt.Printf("count before:%d\n",count)
	err=UpdateCodes("set developer=1","where developer !=1")
	if err!=nil{
		t.Fail()
	}
	count,err=CountCodes("where developer=1")
	if err!=nil{
		t.Fail()
	}
	fmt.Printf("count after:%d\n",count)
}