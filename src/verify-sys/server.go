//server
package main

import (
	"encoding/json"
	"net"
)

const (
	port = ":10005"
)

var clnOnLineChannel = make(chan net.Conn)
var clnOffLineChannel = make(chan net.Conn)

func beginServer() {
	srvSck, err := net.Listen("tcp", port)
	if err != nil {
		log("tcp serve failed!\n"+err.Error(), MSG_ERR)
		return
	}
	log("tcp start,port"+port, MSG_NORMAL)

	defer func() {
		log("tcp serve failed!\n"+err.Error(), MSG_ERR)
		srvSck.Close()
	}()

	go linkEvent()

	//监听到连接通知channel，开启接收线程
	for {
		clnSck, err := srvSck.Accept()
		if err != nil {
			log(err.Error(), MSG_ERR)
			return
		}
		clnOnLineChannel <- clnSck

		go recv(clnSck)
	}

}

//连接接入和断开事件
func linkEvent() {
	for {
		select {
		case clnSck := <-clnOnLineChannel:
			log("online:"+clnSck.RemoteAddr().String(), MSG_DEBUG)
			online(clnSck)
		case clnSck := <-clnOffLineChannel:
			log("offline:"+clnSck.RemoteAddr().String(), MSG_DEBUG)
			offline(clnSck)
			clnSck.Close()
		}
	}
}

//监收信息
func recv(clnSck net.Conn) {
	buf := make([]byte, 1024)
	//for {
	dataLen, err := clnSck.Read(buf)
	if err != nil {
		log(err.Error(), MSG_ERR)
		clnOffLineChannel <- clnSck
		return
	}

	if dataLen > 10 {
		var request Request
		json.Unmarshal(buf[:dataLen], &request)
		//处理消息事务
		getMsg(clnSck, request)
	}
}
