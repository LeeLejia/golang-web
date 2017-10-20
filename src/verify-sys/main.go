//main
package main

import (
	"fmt"
	//"../pdb"
	"./conf"
	"./tcp"
	"net"
)

func main() {
	conf.Init("./app.toml")
	log("初始化数据库", MSG_NORMAL)
	//err:=pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	//if err!=nil{
	//	fmt.Print("数据库配置错误。")
	//}
	log("初始化网络", MSG_NORMAL)
	tcp.BeginServer(OnLine,OffLine,Recv)
}

func OnLine(cln net.Conn){
	fmt.Println("["+cln.RemoteAddr().String()+"]上线！")
}

func OffLine(cln net.Conn){
	fmt.Println("["+cln.RemoteAddr().String()+"]下线！")
}

func Recv(data []byte, cln net.Conn){
	fmt.Println("接收到["+cln.RemoteAddr().String()+"]数据：")
	fmt.Println(string(data))
}
