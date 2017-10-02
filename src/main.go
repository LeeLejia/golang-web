package main

import (
	"./common/conf"
	"fmt"
	"./pdb"
	"net/http"
	"./app"
	"./common"
)

func main() {
	conf.Init("./app.toml")
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
	}
	BeginServer()
}

func BeginServer(){
	/**注册路由*/
	routers:=[]common.BH{
		{Url:"/login",Check:false,Handle:app.Login},
		{Url:"/logout",Check:false,Handle:app.Logout},
		{Url:"/register",Check:false,Handle:app.Register},
		{Url:"/developer/addApps",Check:true,Handle2:app.AddApp},
		{Url:"/developer/list-apps",Check:true,Handle2:app.ListApps},
		{Url:"/upload-picture",Check:false,Handle2:app.UploadPicture},
		{Url:"/upload-file",Check:false,Handle2:app.UploadFile},
		{Url:"/list-file",Check:true,Handle2:app.ListFiles},
	}
	common.SetRouters(routers)
	fsh := http.FileServer(http.Dir(conf.App.StaticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fsh))
	http.ListenAndServe(fmt.Sprintf(":%s", conf.App.ServerPort), nil)
}