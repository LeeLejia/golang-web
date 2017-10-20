//common method

package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

//log类型
const (
	MSG_NORMAL      = 0
	MSG_ERR         = 1
	MSG_DEBUG       = 2
	MSG_LOG         = 3
	MSG_MANAGER_LOG = 4

	DEBUG            = false
	SHOW_LOG         = true
	SHOW_ERR         = true
	SHOW_MANAGER_LOG = false

	LOG_TO_TXT = false
)

var lock sync.Mutex

//打印日志
func log(msg string, msgType int) {
	if LOG_TO_TXT {
		msg = msg + "\r\n"
	}
	switch msgType {
	case MSG_NORMAL:
		fmt.Println(msg)
	case MSG_ERR:
		if SHOW_ERR {
			fmt.Errorf(msg)
		}
	case MSG_DEBUG:
		if DEBUG {
			fmt.Println(msg)
		}
	case MSG_LOG:
		if SHOW_LOG {
			fmt.Println(msg)
		}
	case MSG_MANAGER_LOG:
		if SHOW_MANAGER_LOG {
			fmt.Println(msg)
		}
	}
}

//<<<<<<<<<<<<<<<<<<<<<<<反馈信息结构体
type Response struct {
	ProtoSign int    `json:"protosign"`
	MsgType   int    `json:"msgType"`
	Msg       string `json:"msg"`
	Codes     []Code `json:"codes,omitempty"`
}

//<<<<<<<<<<<<<<<<<<<<<<<请求结构体
type Request struct {
	ProtoSign int    `json:"protosign"`
	MsgType   int    `json:"msgtype"`
	Code      string `json:"code"`
	Machine   string `json:"machine"`
	CodeObj   Code   `json:"codeObj,omitempty"`
	SendTime  int64  `json:"sendtime,omitempty"`
	Msg       string `json:"msg,omitempty"`
}

//验证码
type Code struct {
	CodeString    string `json:"codeString"`
	Machine       string `json:"machine"`
	Exp           string `json:"exp"`
	MostAskTimes  int    `json:"mostAskTimes"`
	AskTimes      int    `json:"askTimes"`
	MMachineCount int    `json:"MMachineCount"`
	ValidType     int    `json:"validType"`
	CreateTime    int    `json:"createTime,omitempty"`
}

func (code *Code) printLog() {
	log("Code:"+code.CodeString, MSG_LOG)
	log("Machines:"+code.Machine, MSG_LOG)
	log("Expression:"+code.Exp, MSG_LOG)
	log("VarifyType:"+getVerifyType(code.ValidType), MSG_LOG)
	if code.MostAskTimes == -1 {
		log("MostAskTime:无上限", MSG_LOG)
	} else {
		log("MostAskTime:"+strconv.Itoa(code.MostAskTimes), MSG_LOG)
	}
	log("AskTime:"+strconv.Itoa(code.AskTimes), MSG_LOG)
	ms := strings.Split(code.Machine, "#")
	if code.MMachineCount == -1 {
		log("MostMachineCount:无上限", MSG_LOG)
	} else {
		log("MostMachineCount:"+strconv.Itoa(code.MMachineCount), MSG_LOG)
	}
	log("MachineCount:"+strconv.Itoa(len(ms)), MSG_LOG)
}
func getTimeString() string {
	t := time.Now().Unix()
	return time.Unix(t, 0).Format("2006-01-02 03:04:05 PM")
}

//{"protosign":1258,"msgtype":0,"machine":"ahgbvikkgenjl15h","code":"12354"}
//{"protosign":1368,"msgtype":2,"code":"12ddds","codeObj":{"codeString":"12ddds","machine":"","exp":"ss你好啊sss","mostAskTimes":20,"askTimes":0,"MMachineCount":-1,"validType":0},"machine":""}
