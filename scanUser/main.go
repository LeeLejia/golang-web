package main

import "github.com/frank2019/mahonia"
//import "github.com/changvvb/alinotify"

var enc = mahonia.NewDecoder("gbk")
func main(){

}

// https://bbs.125.la/home.php?mod=space&uid=%s
func scan(uid string){

}

//1.安装go，配置好GOPATH和GOBIN环境变量(如果没有配置好，自行google)
//2.下载项目
//$ go get github.com/changvvb/alinotify
//3.编译
//$ go install github.com/changvvb/alinotify
//4.运行
//$ alinotify
//5.打开浏览其 https://my.alipay.com 登陆后，获得请求头中的cookies(许多cookie!)复制
//6.打开浏览器 http://127.0.0.1:2048/setcookie ,将复制到的cookies粘贴提交
//7.交易后到 http://127.0.0.1:2048/exam?tel=<电话>&email=< Email> 查看