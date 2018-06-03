package app

import (
	"github.com/cjwddz/fast-model"
	"fmt"
	m "../model"
	"../common/conf"
)
var AppModel model.DbModel
var CodeModel model.DbModel
var UserModel model.DbModel
var VlogModel model.DbModel
var FileModel model.DbModel
var PublishModel model.DbModel
var GoodModel model.DbModel
var OrderModel model.DbModel


func Init(){
	if !model.DbHasInit{
		if conf.App.AppEvn=="dev"{
			model.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.TestDbName, conf.App.DBDriver)
		}else{
			model.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName, conf.App.DBDriver)
		}
	}
	app,err:= m.GetAppModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	code,err:=m.GetCodeModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	user,err:=m.GetUserModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	file,err:=m.GetFileModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	vlog,err:=m.GetVLogModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	publish,err:=m.GetPublishModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	good,err:=m.GetGoodModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	order,err:=m.GetOrderModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	AppModel = app
	CodeModel = code
	UserModel = user
	FileModel = file
	VlogModel = vlog
	PublishModel = publish
	GoodModel = good
	OrderModel = order
}
