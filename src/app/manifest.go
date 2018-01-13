package app

import (
	"../model"
	"fmt"
)
var AppModel model.DbModel
var CodeModel model.DbModel
func Init(){
	app,err:= model.GetAppModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	code,err:=model.GetCodeModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	AppModel = app
	CodeModel = code
}
