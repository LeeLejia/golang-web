//main
package main

import (
	"fmt"
	"../pdb"
	"./conf"
	"./tcp"
	"net"
	"../model"
	"github.com/bitly/go-simplejson"
	"time"
	"encoding/json"
)

func main() {
	conf.Init("./app.toml")
	fmt.Println("初始化数据库")
	err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	if err!=nil{
		fmt.Print("数据库配置错误。")
	}
	fmt.Println("初始化网络")
	tcp.BeginServer(conf.App.TcpAddress,OnLine,OffLine,Recv)
}

func OnLine(cln net.Conn){
	fmt.Println("["+cln.RemoteAddr().String()+"]上线！")
}

func OffLine(cln net.Conn){
	fmt.Println("["+cln.RemoteAddr().String()+"]下线！")
}

func Recv(data []byte, cln net.Conn){
	if len(data)<10{
		fmt.Println("接收到空的数据请求！data:"+string(data))
		return
	}
	js,err:=simplejson.NewJson(data)
	if err!=nil{
		log:=model.T_log{
			Type:model.LOG_TYPE_WARM,
			Tag:"邀请码验证",
			Operator:"客户机",
			Content:"请求格式错误！"+err.Error(),
			CreatedAt:time.Now(),
		}
		cln.Write([]byte("请求格式错误！"))
		log.Insert()
		return
	}
	CheckVerify(data,js,cln)
}

func CheckVerify(data []byte,js *simplejson.Json, cln net.Conn){
	log:=model.T_vlog{
		CreatedAt:time.Now(),
	}
	resp:= struct {
		ProtoSign int    `json:"protosign"`
		MsgType   int    `json:"msgType"`
		Msg       string `json:"msg"`
	}{
		ProtoSign:1258,
		MsgType:0, // 失败
	}
	proto,err:=js.Get("protosign").Int()
	code:= struct {
		machine     string
		code        string
		version     string
		application string
	}{}
	code.code,_ =js.Get("code").String()
	code.machine,_ =js.Get("machine").String()
	code.version,_ =js.Get("version").String()
	code.application,_ =js.Get("application").String()

	log.Machine=code.machine
	log.Code=code.code
	log.Content=string(data)
	log.App=code.application

	if err!=nil || proto!=1258{
		log.Tag=model.VERIFY_LOG_PROTO_NOVALID
		resp.Msg=model.VERIFY_LOG_PROTO_NOVALID
		rs,_:=json.Marshal(resp)
		cln.Write(rs)
		wlog(log)
		return
	}

	if code.code==""{
		log.Tag=model.VERIFY_LOG_EMPTY
		resp.Msg="邀请码不能为空！"
		rs,_:=json.Marshal(resp)
		cln.Write(rs)
		wlog(log)
		return
	}

	codes,err:=model.FindCodes(fmt.Sprintf("where code='%s'",code.code),"","")
	if len(codes)<=0{
		log.Tag=model.VERIFY_LOG_NOEXIST
		resp.Msg="无效邀请码！"
		rs,_:=json.Marshal(resp)
		cln.Write(rs)
		wlog(log)
		return
	}

	if(!codes[0].Valid){
		log.Tag=model.VERIFY_LOG_INVALID
		resp.Msg="该邀请码已被禁用！"
		rs,_:=json.Marshal(resp)
		cln.Write(rs)
		wlog(log)
		return
	}
	fmt.Println(codes[0].EndTime.Unix())
	if(codes[0].EnableTime && codes[0].EndTime.Unix()<time.Now().Unix()){
		log.Tag=model.VERIFY_LOG_AFTER_TIME
		resp.Msg="邀请码已经失效！"
		rs,_:=json.Marshal(resp)
		cln.Write(rs)
		wlog(log)
		return
	}else if(codes[0].EnableTime && codes[0].StartTime.Unix()>time.Now().Unix()){
		log.Tag=model.VERIFY_LOG_BEFORE_TIME
		resp.Msg="该邀请码暂未生效！"
		rs,_:=json.Marshal(resp)
		cln.Write(rs)
		wlog(log)
		return
	}
	log.Tag=model.VERIFY_LOG_SUCCESS
	resp.Msg="邀请码验证成功！"
	resp.MsgType=1
	rs,_:=json.Marshal(resp)
	cln.Write(rs)
	wlog(log)
	// {"protosign":1258,"msgType":0,"machine":"862095022571886","code":"1ab2c","version":"1.0","application":"jwechat"}
}

func wlog(l model.T_vlog){
	// todo 减少并发压力
	err:=l.Insert()
	if err!=nil{
		fmt.Println(err.Error())
	}
}