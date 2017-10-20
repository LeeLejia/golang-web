package tcp

import (
	"net"
	"../conf"
	"fmt"
)

type Handle struct {
	Online  func(net.Conn)
	Offline func(net.Conn)
	Recv    func(data []byte,clnSck net.Conn)
}

var EvnHandle Handle
func BeginServer(online func(net.Conn),offline func(net.Conn),recv func([]byte,net.Conn)) {
	srvSck, err := net.Listen("tcp", conf.App.TcpAddress)
	if err != nil {
		fmt.Println("打开tcp端口失败："+err.Error())
		return
	}
	fmt.Println("tcp start,address:"+conf.App.TcpAddress)
	defer srvSck.Close()
	go linkEvent()
	//监听到连接通知channel，开启接收线程
	EvnHandle.Offline=offline
	EvnHandle.Online = online
	EvnHandle.Recv = recv
	for {
		clnSck, err := srvSck.Accept()
		if err != nil {
			fmt.Println("客户端连入错误！"+err.Error())
			return
		}
		clnOnLineChannel <- clnSck
		go Recv(clnSck)
	}
}