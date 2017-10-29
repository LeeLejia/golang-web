package main

import (
	"./common/conf"
	"fmt"
	"./pdb"
	"net/http"
	"./app"
	"./common"
	"html/template"
	"io/ioutil"
	"strings"
	"time"
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
	/**注册模板*/
	IniTemplate()
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

	http.Handle("/static",http.FileServer(http.Dir(conf.App.StaticPath)))
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
			htmlFiles = append(htmlFiles, conf.App.StaticPath+fileName)
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
	http.HandleFunc("/developer", handler)

}


func xysb(){
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		t1:=time.Now().Unix()
		fmt.Println(fmt.Sprintf("===come time: %s",time.Now().String()))

		bs:=make([]byte,512)
		for {
			c,_:=r.Body.Read(bs)
			if c<=0{
				break
			}
			fmt.Println(string(bs))
		}
		w.Write([]byte("{\"code\":\"ok\";}"))
		fmt.Println(fmt.Sprintf("ok time: %s",time.Now().String()))
		fmt.Println(fmt.Sprintf("时间差 %d",time.Now().Unix()-t1))
	})
}