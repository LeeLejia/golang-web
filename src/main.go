package main

import (
	"./common/conf"
	"fmt"
	"./pdb"
	"net/http"
	"./app"
	"./common"
	"./common/log"
	"html/template"
	"io/ioutil"
	"strings"
)

func main() {
	conf.Init("./app.toml")
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
		return
	}
	log.Init()
	BeginServer()
}

func BeginServer(){
	/**注册模板*/
	IniTemplate()
	/**注册路由*/
	routers:=[]common.BH{
		{Url:"/api/login",Check:false,Handle:app.Login},
		{Url:"/api/logout",Check:false,Handle:app.Logout},
		{Url:"/api/register",Check:false,Handle:app.Register},
		{Url:"/api/developer/addApps",Check:true,Handle2:app.AddApp},
		{Url:"/api/developer/list-apps",Check:true,Handle2:app.ListApps},
		{Url:"/api/upload-picture",Check:false,Handle2:app.UploadPicture},
		{Url:"/api/upload-file",Check:false,Handle2:app.UploadFile},
		{Url:"/api/list-file",Check:true,Handle2:app.ListFiles},
	}
	common.SetRouters(routers)
	http.Handle("/",http.FileServer(http.Dir(conf.App.StaticPath)))
	fmt.Println("开始服务！")
	err:=http.ListenAndServe(fmt.Sprintf(":%s", conf.App.ServerPort), nil)
	if err!=nil{
		fmt.Println("服务退出！"+err.Error())
	}
}

/**
初始化模板
 */
func IniTemplate(){
	/*从静态目录中加载并设置模板*/
	fmt.Println(fmt.Sprintf("初始化模板,%s",conf.App.StaticPath))
	files, err := ioutil.ReadDir(conf.App.StaticPath)
	if err != nil {
		fmt.Println("初始化模板失败！detail："+err.Error())
		return
	}
	var htmlFiles []string
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".html") {
			htmlFiles = append(htmlFiles, conf.App.StaticPath+"/"+fileName)
			fmt.Println(fmt.Sprintf("---添加模板,%s",htmlFiles))
		}
	}
	templates:=template.Must(template.ParseFiles(htmlFiles...))
	/*设置前端入口模板*/
	index:=templates.Lookup("index.html")
	handler:=func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html;charset=utf-8")
		if request.Method == "GET" {
			index.Execute(writer, nil)
			return
		}
		http.Redirect(writer,request,"/",http.StatusFound)
	}
	http.HandleFunc("/", handler)

}
