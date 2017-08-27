/**
 * 日志实现
 *
 * @author zhangqichao
 * Created on 2017-08-18
 */
package logger

import (
	"fmt"
	"../../../back-end-api/model"
	"mypack/log"
)


func InitLogger() {

}

func Debug(msg, path string) {
	log.D(msg+" [path]%s",path)
	InsertSystemLogs("debug", msg, path)
}

func Error(msg, path string) {
	log.E(msg+" [path]%s",path)
	InsertSystemLogs("error", msg, path)
}

func Info(msg, path string) {
	log.I(msg+" [path]%s",path)
	InsertSystemLogs("info", msg, path)
}

func InsertSystemLogs(types, msg, path string) {
	systemLog := model.SystemLog{}
	systemLog.Type = types
	systemLog.Content = msg
	systemLog.Part = path
	err := systemLog.Insert()
	if err != nil {
		fmt.Println(err.Error())
	}
}
