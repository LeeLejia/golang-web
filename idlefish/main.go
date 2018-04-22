package main

import (
	"net/http"
	"fmt"
	"os"
	"math/rand"
	"time"
	fm "github.com/cjwddz/fast-model"
	"io"
	"./model"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"encoding/json"
	"strings"
)
/**
DROP TABLE IF EXISTS xy_pic;
CREATE TABLE xy_pic (
"file" varchar(32) COLLATE "default"
)
WITH (OIDS=FALSE);
 */
var(
	picDir = "./pic"
)
// 在线设备
var onlineMachines []string
var Ponds []FishPond
/**
检查文件目录
 */
func checkDirs(){
	if _,err:=os.Stat(picDir);err!=nil{
		err:=os.MkdirAll(picDir,0755)
		if err!=nil{
			fmt.Println("创建文件夹失败！")
		}
	}
}

/**
获取随机字符串
 */
func  GetRandomString(l int) string {
	str := "0123456789ASDFGHJKLPOIUYTREWQZXCVBNM"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
type FishPond struct{
	Machine string   `json:"machine"`
	Ponds   []string `json:"ponds"`
}
func main(){
	defer func() {
		fps,err:=json.Marshal(Ponds)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err=ioutil.WriteFile("./fishpond.json",fps,0777)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}()
	initPonds()
	checkDirs()
	fm.InitDB("www.cjwddz.cn","5432","xy","xy123","xy","postgres")
	tm,err:=model.GetTaskModel()
	if err!=nil{
		fmt.Println(err)
		panic("找不到数据表对象.")
		return
	}
	pic,err:=model.GetPicModel()
	if err!=nil{
		fmt.Println(err)
		panic("找不到数据表对象.")
		return
	}
	onlineMachines=make([]string,0)
	http.Handle("/",http.FileServer(http.Dir("./")))
	// 获取机器及鱼塘配置
	http.HandleFunc("/api/machines", func(writer http.ResponseWriter, request *http.Request) {
		fps,err:=json.Marshal(Ponds)
		if err != nil {
			fmt.Fprint(writer,fmt.Sprintf(`{"res":false,"msg":'%s'}`,err.Error()))
			return
		}
		fmt.Fprint(writer,fmt.Sprintf(`{"res":true,"data":%s}`,string(fps)))
	})
	// 更新鱼塘配置
	http.HandleFunc("/api/setPond", func(writer http.ResponseWriter, request *http.Request) {
		name:=request.FormValue("name")
		ponds:=request.FormValue("ponds")
		if len(name)<4{
			fps,_:=json.Marshal(Ponds)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Fprint(writer,fmt.Sprintf(`{"res":false,"msg":"参数出错！","data":%s}`,string(fps)))
			return
		}
		for i := range Ponds {
			if Ponds[i].Machine == name {
				Ponds = append(Ponds[:i], Ponds[i+1:]...)
				fps,_:=json.Marshal(Ponds)
				if err != nil {
					fmt.Println(err.Error())
				}
				if len(ponds)==0{ // 删除的情况
					fmt.Fprint(writer,fmt.Sprintf(`{"res":true,"msg":"更新完成，删除设备记录成功！","data":%s}`,string(fps)))
					return
				}
				break
			}
		}
		if len(ponds)==0 { // 删除的情况
			fps,_:=json.Marshal(Ponds)
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Fprint(writer,fmt.Sprintf(`{"res":true,"msg":"更新完成！","data":%s}`,string(fps)))
			return
		}

		fp:=FishPond{
			Machine:name,
			Ponds:[]string{},
		}
		for _,v:=range strings.Split(ponds,"\r\n"){
			if len(v)!=0{
				fp.Ponds = append(fp.Ponds,v)
			}
		}
		// 添加的情况
		Ponds=append(Ponds,fp)
		fps,_:=json.Marshal(Ponds)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Fprint(writer,fmt.Sprintf(`{"res":true,"msg":"添加成功！","data":%s}`,string(fps)))
	})
	// 填写任务
	http.HandleFunc("/addTask", func(writer http.ResponseWriter, request *http.Request) {
		f, err := os.Open("./index.html")
		if err != nil {
			fmt.Fprint(writer,fmt.Sprintf(`{"res":false,"msg":'%s'}`,err.Error()))
			return
		}
		io.Copy(writer,f)
		return
	})
	// 获取任务
	http.HandleFunc("/getTask", func(writer http.ResponseWriter, request *http.Request) {
		t:=request.FormValue("time")
		if t==""{
			fmt.Fprint(writer,fmt.Sprintf(`{"res":false,"msg":'%s'}`,"time参数不能为空！"))
			return
		}
		mc:=request.FormValue("machine")
		rs,err:=tm.Query(fm.DbCondition{}.And2(">","created",t).And2("like","machine",mc))
		if err!=nil{
			fmt.Println(err)
			fmt.Fprint(writer,fmt.Sprintf(`{"res":false,"msg":'%s'}`,err.Error()))
			return
		}
		bs,err:=json.Marshal(rs)
		if err!=nil{
			fmt.Println(err)
			fmt.Fprint(writer,fmt.Sprintf(`{"res":false,"msg":'%s'}`,err.Error()))
			return
		}
		fmt.Fprint(writer,string(bs))
		return
	})
	// 处理表单
	http.HandleFunc("/api/upload", func(writer http.ResponseWriter, request *http.Request) {
		task:=model.Task{
			Machine:	 request.FormValue("machines"),
			Type:        request.FormValue("type"),
			Title:       request.FormValue("title"),
			Addr:        request.FormValue("addr"),
			RawMoney:    request.FormValue("rawMoney"),
			Desc:        request.FormValue("desc"),
			PublicType:  request.FormValue("publicType"),
			Category:    request.FormValue("category"),
			Money:       request.FormValue("money"),
			TranMoney:   request.FormValue("tranMoney"),
			Pond:        request.FormValue("pond"),
			Pics:        request.FormValue("pics"),
			CreatedTime: time.Now(),
		}
		fmt.Println(task)
		err:=tm.Insert(task)
		if err!=nil{
			fmt.Println(err)
			fmt.Fprint(writer,"fail")
			return
		}
		fmt.Fprint(writer,"success")
		return
	})
	// 处理文件
	http.HandleFunc("/api/uploadFile", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseMultipartForm(1<<32)
		fmt.Println(request.PostForm)
		f,h,_:=request.FormFile("file")
		fmt.Println(h.Filename)
		key:=GetRandomString(32)
		filePath := picDir +"/"+ key
		t, err := os.Create(filePath)
		if err != nil {
			fmt.Fprint(writer,`{"res":false,"msg":'创建文件失败！'}`)
			return
		}
		if _,err:=io.Copy(t,f);err!=nil{
			fmt.Fprint(writer,`{"res":false,"msg":'写入文件失败'}`)
			return
		}
		err=pic.Insert(model.Pic{Name:key})
		if err!=nil{
			fmt.Println(err)
			fmt.Fprint(writer,"fail")
			return
		}
		fmt.Fprint(writer,key)
	})
	// 上下线登记
	http.HandleFunc("/api/reg", func(writer http.ResponseWriter, request *http.Request) {
		online:=request.FormValue("online")
		id:=request.FormValue("id")
		for i := range onlineMachines {
			if onlineMachines[i] == id {
				if online=="on"{
					// 已经存在
					fmt.Fprint(writer,"ID已经存在！")
					return
				}else{
					onlineMachines = append(onlineMachines[:i], onlineMachines[i+1:]...)
					fmt.Fprint(writer,"登出成功！")
					return
				}
			}
		}
		if online=="on"{
			onlineMachines = append(onlineMachines,id)
			fmt.Fprint(writer,"登入成功！")
		}
	})
	// 退出
	http.HandleFunc("/exit", func(writer http.ResponseWriter, request *http.Request) {
		os.Exit(0)
	})
	err=http.ListenAndServe(fmt.Sprintf(":%d", 6789), nil)
	if err!=nil{
		fmt.Println("服务退出！"+err.Error())
	}
}

func initPonds(){
	Ponds=make([]FishPond,0)
	bs,err:=ioutil.ReadFile("./fishpond.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	js,err:=simplejson.NewJson(bs)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	is,err:=js.Array()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, i := range is {
		fp, _ := i.(map[string]interface{})
		machine := fp["machine"]
		ponds := fp["ponds"]
		if machine==nil || ponds==nil{
			continue
		}
		ps:=make([]string,len(ponds.([]interface{})))
		for j,p:=range ponds.([]interface{}){
			ps[j]=p.(string)
		}
		t:=FishPond{Machine:machine.(string), Ponds:ps}
		Ponds = append(Ponds,t)
	}
}