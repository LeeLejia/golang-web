package model

import (
	"testing"
	"fmt"
	"../common/conf"
	"../pdb"
	"time"
)

func TestT_vlog_Insert(t *testing.T) {
	conf.Init(conf_path)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	its:=[]T_vlog{
		{Tag:"验证成功",Code:"14gaseg211521515154545415",App:"old1234", Machine:"sgeafgweafjpagewjqogjapwikpfw",Content:"测试日志-长",CreatedAt:time.Now()},
		{Tag:"验证成功",Code:"1",App:"old1234", Machine:"sgeafgweafjpagewjqogjapwikpfw",Content:"测试日志-短",CreatedAt:time.Now()},
		{Tag:"验证成功",Code:"3sgawefs",App:"old1234", Machine:"sgeafgweafjpagewjqogjapwikpfw",Content:`{"protosign":1258,"msgType":0,"machine":"862095022571886","code":"1ab2c","version":"1.0","application":"jwechat"}`,CreatedAt:time.Now()},
	}
	for _,it:=range its{
		err=it.Insert()
		if err==nil{
			t.Fail()
		}
	}
}
func TestCountVLogs(t *testing.T) {
	conf.Init(conf_path)
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		t.Fail()
	}
	count,_:=CountVLogs("")
	fmt.Println(fmt.Sprintf("无条件：%d",count))
	count,_=CountVLogs("where Code=1")
	fmt.Println(fmt.Sprintf("有条件：%d",count))
}
// todo add