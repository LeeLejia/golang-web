package main

import (
	"./common/conf"
	"fmt"
	"net/http"
	"./app"
	"./common"
	"./common/log"
	"html/template"
	"io/ioutil"
	"strings"
)

func main() {
	// 初始化配置文件
	conf.Init("./app.toml")
	// 初始化模型
	app.Init()
	// 初始化日志系统
	log.Init()
	defer func() {
		log.Flush()
	}()
	BeginServer()
}

func BeginServer(){
	InitTemplate()
	/**注册路由*/
	routers:=[]common.BH{
		// 用户登录/注销/注册
		{Url:"/api/login",Check:false,Handle:app.Login},
		{Url:"/api/logout",Check:true,Handle:app.Logout},
		{Url:"/api/register",Check:false,Handle:app.Register},
		// 文件校验/上传图片/上传文件
		{Url:"/api/checkSha256",Check:false,Handle:app.CheckSha256},
		{Url:"/api/uploadPicture",Check:true,Handle:app.UploadPicture},
		{Url:"/api/uploadFile",Check:true,Handle:app.UploadFile},

		//{Url:"/api/developer/add-app",Check:true,Handle2:app.AddApp},
		//{Url:"/api/developer/list-apps",Check:true,Handle2:app.ListApps},
		//{Url:"/api/developer/add-code",Check:true,Handle2:app.AddCode},
		//{Url:"/api/developer/list-codes",Check:true,Handle2:app.ListCodes},
	}
	common.SetRouters(routers)
	http.Handle("/",http.FileServer(http.Dir(conf.App.StaticPath)))
	fmt.Println("开始服务！")
	err:=http.ListenAndServe(fmt.Sprintf(":%s", conf.App.ServerPort), nil)
	if err!=nil{
		fmt.Println("服务退出！"+err.Error())
	}
}
/** 初始化模板*/
func InitTemplate(){
	/**获取模板,设置页面*/
	index,err:=getTemplate()
	if err!=nil{
		fmt.Println(fmt.Sprintf("设置模板出错!原因：%s",err.Error()))
	}
	handler:=func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/html;charset=utf-8")
		if request.Method == "GET" && err==nil{
			index.Execute(writer, nil)
			return
		}
		http.Redirect(writer,request,"/",http.StatusFound)
	}
	http.HandleFunc("/fangtan", handler)
	// 更新模板文件
	http.HandleFunc("/resetTemplate", func(writer http.ResponseWriter, request *http.Request) {
		index,err=getTemplate()
		if err!=nil{
			rs:=fmt.Sprintf("设置模板出错!原因：%s",err.Error())
			fmt.Println(rs)
			writer.Write([]byte(rs))
			return
		}
		rs:=fmt.Sprintf("更新成功，入口模板为%s",index.Name())
		fmt.Println(rs)
		writer.Write([]byte(rs))
		return
	})
}
/**
	获取模板入口
 */
func getTemplate() (*template.Template,error){
	/*从静态目录中加载并设置模板*/
	fmt.Println(fmt.Sprintf("从[%s]中加载模板..",conf.App.StaticPath))
	files, err := ioutil.ReadDir(conf.App.StaticPath)
	if err != nil {
		return nil,err
	}
	var htmlFiles []string
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".html") {
			htmlFiles = append(htmlFiles, conf.App.StaticPath+"/"+fileName)
			fmt.Println(fmt.Sprintf("添加模板:%s",fileName))
		}
	}
	if len(htmlFiles)==0{
		return nil,fmt.Errorf("在%s中找不到模板文件.",conf.App.StaticPath)
	}
	templates:=template.Must(template.ParseFiles(htmlFiles...))
	/*设置前端入口模板*/
	index:=templates.Lookup("index.html")
	return index,nil
}
