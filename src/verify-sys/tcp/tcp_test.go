package tcp

import (
	"testing"
	"net"
	"fmt"
	"os"
	"bytes"
	"time"
)

func TestBeginServer(t *testing.T) {
	//go BeginServer(":8881",OnLine,OffLine,recv)
	fmt.Println("openOK")
	time.Sleep(time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:10005")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	conn.Write([]byte(`{"protosign":1258,"msgType":0,"machine":"862095022571886","code":"112235","version":"1.0","application":"jwechat"}`))
	bs:=bytes.Buffer{}
	conn.Read(bs.Bytes())
	buf := make([]byte, 4096)
	conn.Read(buf)
	fmt.Println(string(buf))
	os.Exit(0)
}

func OnLine(cln net.Conn){
	fmt.Println("["+cln.RemoteAddr().String()+"]上线！")
}

func OffLine(cln net.Conn){
	fmt.Println("["+cln.RemoteAddr().String()+"]下线！")
}

func recv(data []byte, cln net.Conn){
	fmt.Println("接收到["+cln.RemoteAddr().String()+"]数据：")
	fmt.Println(string(data))
	if string(data)!=`{"protosign":1258,"msgType":0,"machine":"862095022571886","code":"1ab2c","version":"1.0","application":"jwechat"}`{
		fmt.Println("test failed!!")
		return
	}
	cln.Write([]byte("i get message!"))
	fmt.Println("test success!!")
}
