package log

import (
	"testing"
	"fmt"
	"../conf"
	".."
	"time"
)

func TestD(t *testing.T) {
	common.InitDb(t)
	conf.Init("/home/cjwddz/桌面/git-project/golang-web/src/app.toml")
	fmt.Println(fmt.Sprintf("cache:%d,logInterval:%d",cache,logInterval))
	Init()
	for i:=0;i<50;i++{
		D("test","user","msg%d",i)
		fmt.Println(count)
	}
	time.Sleep(time.Second*11)
	D("test","user","msg%d",123)
	for i:=0;i<100;i++{
		D("test","user","msg%d",i)
		fmt.Println(count)
	}
	fmt.Println("ok")
}
