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
	"os"
	"./plugin"
)

func main() {
	// 初始化配置文件
	conf.Init("./app.toml")
	// 初始化模型
	app.Init()
	// 初始化日志系统
	log.Init()
	// 初始化插件
	go plugin.Init()
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
		{Url:"/api/setUserAvatar",Check:true,Handle:app.SetUserAvatar},
		// 发布任务
		{Url:"/api/publish",Check:true,Handle:app.Publish},
		{Url:"/api/getTask",Check:false,Handle:app.GetTask},
		// 文件校验/上传图片/上传文件
		{Url:"/api/checkSha256",Check:true,Handle:app.CheckSha256},
		{Url:"/api/uploadFile",Check:true,Handle:app.UploadFile},
		{Url:"/api/listFiles",Check:true,Handle:app.ListFiles},
		{Url:"/api/deleteFile",Check:true,Handle:app.DeleteFile},
		// App添加/删除/列表获取
		{Url:"/api/developer/addApp",Check:true,Handle:app.AddApp},
		{Url:"/api/developer/listApps",Check:true,Handle:app.ListApps},
		// 商品和交易
		{Url:"/api/pay",Check:true,Handle:app.Pay},
		{Url:"/api/getGoods",Check:true,Handle:app.GetGoods},
		{Url:"/api/addGood",Check:true,Handle:app.AddGood},
		{Url:"/api/getOrders",Check:true,Handle:app.GetOrders},
		{Url:"/api/pay/notify",Check:true,Handle:func(sess *common.Session,w http.ResponseWriter, r *http.Request){
			r.ParseForm()
			fmt.Printf("notify_callback!!")
			fmt.Printf(fmt.Sprintf("%v",r.Form))
			log.N("NOTIFY","",fmt.Sprintf("%v",r.Form))
		}},
		//{Url:"/api/developer/add-code",Check:true,Handle2:app.AddCode},
		//{Url:"/api/developer/list-codes",Check:true,Handle2:app.ListCodes},
	}
	common.SetRouters(routers)
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
		if request.Method == "GET"{
			err:=index.Execute(writer, nil)
			if err!=nil{
				http.Redirect(writer,request,"/404.html",http.StatusFound)
			}
			return
		}
		http.Redirect(writer,request,"/404.html",http.StatusFound)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		f,err := os.Stat(conf.App.StaticPath + r.URL.String())
		if err != nil || f.IsDir() {
			// 不是文件
			handler(w,r)
			return
		}
		http.FileServer(http.Dir(conf.App.StaticPath)).ServeHTTP(w,r)
	})
	// 更新模板文件
	http.HandleFunc("/resetTemplate", func(writer http.ResponseWriter, request *http.Request) {
		index,err=getTemplate()
		if err!=nil{
			rs:=fmt.Sprintf("设置模板出错!原因：%s",err.Error())
			writer.Write([]byte(rs))
			return
		}
		rs:=fmt.Sprintf("更新成功，入口模板为%s",index.Name())
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
		return nil,fmt.Errorf("在%s中找不到模板文件",conf.App.StaticPath)
	}
	templates:=template.Must(template.ParseFiles(htmlFiles...))
	/*设置前端入口模板*/
	index:=templates.Lookup("index.html")
	return index,nil
}