/**
	日志类
 */
package log

import (
	"fmt"
	"../../model"
	"time"
)

const (
	color_red = uint8(iota + 91)
	color_green
	color_yellow
	color_blue
	color_magenta //洋红
)

/**
	是否插入数据库
 */
var CONF_WRITE_TO_DB = true

/**
	debug
 */
func D(tag string,user string,msg ...interface{}) {
	var content=(msg[0]).(string)
	if len(msg)>1{
		content=fmt.Sprintf(msg[0].(string),msg[1:])
	}
	fmt.Println(fmt.Sprintf("D[tag:%s user:%s %s] %s",tag,user,time.Now().Format("2006/01/02 15:04:05"),content))
	insertSystemLogs("debug",tag,user,content)
}
/**
	info
 */
func I(tag string,user string,msg ...interface{}) {
	var content=(msg[0]).(string)
	if len(msg)>1{
		content=fmt.Sprintf(msg[0].(string),msg[1:])
	}
	fmt.Println(Blue(fmt.Sprintf("I[tag:%s user:%s %s] %s",tag,user,time.Now().Format("2006/01/02 15:04:05"),content)))
	insertSystemLogs("info",tag,user,content)
}
/**
	error
 */
func E(tag string,user string,msg ...interface{}) {
	var content=(msg[0]).(string)
	if len(msg)>1{
		content=fmt.Sprintf(msg[0].(string),msg[1:])
	}
	fmt.Println(Red(fmt.Sprintf("E[tag:%s user:%s %s] %s",tag,user,time.Now().Format("2006/01/02 15:04:05"),content)))
	insertSystemLogs("error",tag,user,content)
}

/**
	写入数据库
 */
func insertSystemLogs(logType,tag, user, content string) {
	if !CONF_WRITE_TO_DB{
		return
	}
	Log := model.T_log{}
	Log.Tag=tag
	Log.Type=logType
	Log.User=user
	Log.Content= content
	err := Log.Insert()
	if err != nil {
		fmt.Println(Red(err.Error()))
	}
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
