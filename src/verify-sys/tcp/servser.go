package tcp

import (
	"net"
)

var clnOnLineChannel = make(chan net.Conn)
var clnOffLineChannel = make(chan net.Conn)

func linkEvent() {
	for {
		select {
		case clnSck := <-clnOnLineChannel:
			EvnHandle.Online(clnSck)
		case clnSck := <-clnOffLineChannel:
			EvnHandle.Offline(clnSck)
		}
	}
}

//消息到达
func Recv(clnSck net.Conn) {
	// todo
	buf := make([]byte, 4096)
	for {
		dataLen, err := clnSck.Read(buf)
		if err != nil {
			clnOffLineChannel <- clnSck
			return
		}
		//过滤心跳等
		if dataLen > 10 {
			EvnHandle.Recv(buf,clnSck)
		}
	}
}
