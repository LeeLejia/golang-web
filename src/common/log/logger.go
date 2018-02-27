/**
	日志类
 */
package log

import (
	"fmt"
	"../../model"
	fm "github.com/cjwddz/fast-model"
	"time"
	"../conf"
	"runtime"
	"strings"
)

const (
	color_red = uint8(iota + 91)
	color_green
	color_yellow
	color_blue
	color_magenta //洋红
)

/**是否插入数据库*/
var CONF_WRITE_TO_DB = true
var FMT_OUT = true

var logs []model.T_log
var cache int
var count int
var logInterval int
var lastTime int64
var LogModel fm.DbModel
func Init() {
	cache = conf.App.LogCache
	logInterval = conf.App.LogInterval
	if cache<=0{
		cache=100
	}
	if logInterval<=0{
		logInterval=10
	}
	lastTime = time.Now().Unix()
	count=0
	fmt.Println(fmt.Sprintf("log system->cache:%d,logInterval:%d,lastTime:%d",cache,logInterval,lastTime))
	logs=make([]model.T_log,cache)
	l,err:=model.GetLogModel()
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
	LogModel = l
}
/**
	debug
 */
func D(tag string,operator string,msg ...interface{}) {
	var content=(msg[0]).(string)
	if len(msg)>1{
		content=fmt.Sprintf(msg[0].(string),msg[1:])
	}
	if FMT_OUT{
		fmt.Println(fmt.Sprintf("D[tag:%s,operator:%s,time:%s]>>%s",tag,operator,time.Now().Format("2006/01/02 15:04:05"),content))
	}
	insertSystemLogs("debug",tag,operator,content)
}
/**
	warn
 */
func W(tag string,operator string,msg ...interface{}) {
	var content=(msg[0]).(string)
	if len(msg)>1{
		content=fmt.Sprintf(msg[0].(string),msg[1:])
	}
	if FMT_OUT{
		fmt.Println(fmt.Sprintf("W[tag:%s,operator:%s,time:%s]>>%s",tag,operator,time.Now().Format("2006/01/02 15:04:05"),content))
	}
	insertSystemLogs("warn",tag,operator,content)
}
/**
	info
 */
func I(tag string,operator string,msg ...interface{}) {
	var content=(msg[0]).(string)
	if len(msg)>1{
		content=fmt.Sprintf(msg[0].(string),msg[1:]...)
	}
	if FMT_OUT {
		fmt.Println(Blue(fmt.Sprintf("I[tag:%s,operator:%s,time:%s]>>%s", tag, operator, time.Now().Format("2006/01/02 15:04:05"), content)))
	}
	insertSystemLogs("info",tag,operator,content)
}
/**
	error
 */
func E(tag string,operator string,msg ...interface{}) {
	var content=(msg[0]).(string)
	if len(msg)>1{
		content=fmt.Sprintf(msg[0].(string),msg[1:])
	}
	if FMT_OUT {
		fmt.Println(Red(fmt.Sprintf("E[tag:%s operator:%s %s] %s", tag, operator, time.Now().Format("2006/01/02 15:04:05"), content)))
	}
	insertSystemLogs("error",tag,operator,content)
}
/**
	normal
 */
func N(tag string,operator string,msg ...interface{}) {
	var content=(msg[0]).(string)
	if len(msg)>1{
		content=fmt.Sprintf(msg[0].(string),msg[1:])
	}
	if FMT_OUT {
		fmt.Println(Green(fmt.Sprintf("N[tag:%s,operator:%s,time:%s]>>%s", tag, operator, time.Now().Format("2006/01/02 15:04:05"), content)))
	}
	insertSystemLogs("normal",tag,operator,content)
}
/**
提示用户
 */
func Notify(tag string,operator string,msg ...interface{}){
	var content=(msg[0]).(string)
	if len(msg)>1{
		content=fmt.Sprintf(msg[0].(string),msg[1:])
	}
	if FMT_OUT {
		fmt.Println(Magenta(fmt.Sprintf("Notify[tag:%s,operator:%s,time:%s]>>%s", tag, operator, time.Now().Format("2006/01/02 15:04:05"), content)))
	}
	// todo notify action
	insertSystemLogs("notify",tag,operator,content)
}

/**
	写入数据库
 */
func insertSystemLogs(logType,tag, operator, content string) {
	if !CONF_WRITE_TO_DB{
		return
	}
	_,f,l,_:=runtime.Caller(2)
	cl:=fmt.Sprintf("%s line:%d",strings.TrimPrefix(f,conf.AppRootPath),l)
	tl := model.T_log{
		Tag:tag,
		Type:logType,
		Caller:cl,
		Operator:operator,
		Content:content,
	}
	logs[count]=tl
	count++
	if count>=cache || int(time.Now().Unix()-lastTime)>logInterval{
		Flush()
	}
	lastTime=time.Now().Unix()
}

/**
	强制写出
	使用事务批量写出日志和单条插入日志分别占用时间单位(100,200)：20：13,51:34
 */
func Flush(){
	for i:=0;i<count;i++{
		err:=LogModel.Insert(logs[i])
		if err != nil {
			fmt.Println(Red(err.Error()))
		}
	}
	count=0
}
/**
	获取控制台红色格式文本
 */
func Red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}
/**
	获取控制台绿色格式文本
 */
func Green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}
/**
	获取控制台黄色格式文本
 */
func Yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_yellow, s)
}
/**
	获取控制台蓝色格式文本
 */
func Blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}
/**
	获取控制台洋红色格式文本
 */
func Magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}
