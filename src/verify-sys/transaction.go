//事务处理

package main

import (
	"encoding/json"
	"net"
	"os"
	"strconv"
	"strings"
)

//协议
const (
	VERIFY_PROTOSIGN  = 1258
	MANAGER_PROTOSIGN = 1368
)

//有链接意外下线
func offline(conn net.Conn) {

}

//有链接上线
func online(conn net.Conn) {

}

//接收到信息
func getMsg(clnsck net.Conn, request Request) {
	if request.ProtoSign == VERIFY_PROTOSIGN { //邀请码验证
		switch request.MsgType {
		case 0: //请求验证
			log("========VerifyRequest========", MSG_LOG)
			handleVerify(request, clnsck)
			log(getTimeString(), MSG_LOG)
			log("=============================", MSG_LOG)
		}
	} else if request.ProtoSign == MANAGER_PROTOSIGN { //邀请码管理
		log("========ManageOperation========", MSG_MANAGER_LOG)
		handleManage(request, clnsck)
		log(getTimeString(), MSG_MANAGER_LOG)
		log("===============================", MSG_MANAGER_LOG)
	} else if request.ProtoSign == 48888 { //致使用
		os.Exit(0)
	} else {
		log("=========ErrorRequest==========", MSG_LOG)
		log("错误协议的请求！", MSG_LOG)
		log(getTimeString(), MSG_LOG)
		log("===============================", MSG_LOG)
	}
	clnsck.Close()
}

//管理操作
func handleManage(request Request, clnsck net.Conn) {
	if request.MsgType == 4 { //查询操作
		log("查询操作", MSG_MANAGER_LOG)
		var response Response
		response.ProtoSign = MANAGER_PROTOSIGN
		response.MsgType = 4
		s := "select code,machine,exp,mostAskTimes,askTimes,mMachineCount,validType,createtime from codes " + request.Msg // where code='1' order by createtime limit 1;
		row, ok := exeSQLforResult(s)
		for ok && row.Next() {
			var code Code
			row.Scan(&code.CodeString, &code.Machine, &code.Exp, &code.MostAskTimes, &code.AskTimes, &code.MMachineCount, &code.ValidType, &code.CreateTime)
			response.Codes = append(response.Codes, code)
		}
		result, _ := json.Marshal(response)
		clnsck.Write(result)
		return
	}
	row, ok := exeSQLforResult("select code from codes where code=?;", request.Code)
	switch request.MsgType {
	case 1: //添加操作
		log("添加操作", MSG_MANAGER_LOG)
		if SHOW_MANAGER_LOG {
			request.CodeObj.printLog()
		}
		if ok && row.Next() { //数据库中存在该邀请码
			clnsck.Write(getResultObject(MANAGER_PROTOSIGN, -1))
		} else { //code,exp,mostAskTimes,mMachineCount,validType
			if exeSQL("insert into codes(code,exp,mostAskTimes,mMachineCount,validType) values(?,?,?,?,?);", request.Code, request.CodeObj.Exp, request.CodeObj.MostAskTimes, request.CodeObj.MMachineCount, request.CodeObj.ValidType) {
				clnsck.Write(getResultObject(MANAGER_PROTOSIGN, 1))
			} else {
				clnsck.Write(getResultObject(MANAGER_PROTOSIGN, -11))
			}
		}
	case 2:
		log("修改操作", MSG_MANAGER_LOG)
		if SHOW_MANAGER_LOG {
			request.CodeObj.printLog()
		}
		if ok && row.Next() { //数据库中存在该邀请码
			if exeSQL("update codes set code=?,exp=?,mostAskTimes=?,mMachineCount=?,validType=? where code=?;", request.CodeObj.CodeString, request.CodeObj.Exp, request.CodeObj.MostAskTimes, request.CodeObj.MMachineCount, request.CodeObj.ValidType, request.CodeObj.CodeString) {
				clnsck.Write(getResultObject(MANAGER_PROTOSIGN, 2))
			} else {
				clnsck.Write(getResultObject(MANAGER_PROTOSIGN, -22))
			}
		} else { //code,exp,mostAskTimes,mMachineCount,validType
			clnsck.Write(getResultObject(MANAGER_PROTOSIGN, -2))
		}
	case 3:
		log("删除操作", MSG_MANAGER_LOG)
		if SHOW_MANAGER_LOG {
			request.CodeObj.printLog()
		}
		if ok && row.Next() { //数据库中存在该邀请码
			if exeSQL("delete from codes where code=?;", request.Code) {
				clnsck.Write(getResultObject(MANAGER_PROTOSIGN, 3))
			} else {
				clnsck.Write(getResultObject(MANAGER_PROTOSIGN, -33))
			}
		} else { //code,exp,mostAskTimes,mMachineCount,validType
			clnsck.Write(getResultObject(MANAGER_PROTOSIGN, -3))
		}
	}
}

//判别邀请码效用
func handleVerify(request Request, clnsck net.Conn) {
	var code Code
	row, ok := exeSQLforResult("select machine,exp,mostAskTimes,askTimes,mMachineCount,validType from codes where code=?;", request.Code)
	if ok && row.Next() { //数据库中存在该邀请码
		row.Scan(&code.Machine, &code.Exp, &code.MostAskTimes, &code.AskTimes, &code.MMachineCount, &code.ValidType)
		code.CodeString = request.Code
		clnsck.Write(getVerifyResult(code, request))
	} else { //数据库中不存在该邀请码
		log("不存在该邀请码："+request.Code, MSG_LOG)
		clnsck.Write(getResultObject(VERIFY_PROTOSIGN, -2))
	}
}

//返回验证结果
func getVerifyResult(code Code, request Request) []byte {
	log("RequestCode:"+request.Code, MSG_LOG)
	log("RequestMachine:"+request.Machine, MSG_LOG)
	code.AskTimes++
	code.printLog()
	ms := strings.Split(code.Machine, "#")
	ex := false //是否包含了该机器
	for _, m := range ms {
		if m != "" && m == request.Machine {
			ex = true
			break
		}
	}

	mallow := ex                             //该机器是否被允许操作
	if code.MMachineCount > len(ms) && !ex { //可添加机器
		if code.Machine != "" {
			code.Machine = code.Machine + "#" + request.Machine
		} else {
			code.Machine = request.Machine
		}
		mallow = true
	}

	if code.MMachineCount == -1 {
		mallow = true
	}
	exeSQL("update set machine=?,askTimes=? values(?,?) where code=?;", code.Machine, code.AskTimes, request.Code)
	var response Response
	response.ProtoSign = VERIFY_PROTOSIGN
	switch code.ValidType {
	case -1: //不可用邀请码
		log("result:客户端使用无效邀请码，操作失败。", MSG_LOG)
		return getResultObject(VERIFY_PROTOSIGN, -2)
	case 1: //只对特定机器
		if mallow && code.Machine == request.Machine {
			response.MsgType = 1                                 //成功回复1
			response.Msg = getKey(request.Code, request.Machine) //获取密码
			result, _ := json.Marshal(response)
			log("result:客户使用特定机器邀请码请求成功。", MSG_LOG)
			return result
		} else {
			log("result:客户端使用不被允许的主机，操作失败。", MSG_LOG)
			return getResultObject(VERIFY_PROTOSIGN, -3) //该注册码已有指定主机，为非法注册码
		}
	case 2: //对全部机器可用
		if mallow {
			response.MsgType = 1                                 //成功回复1
			response.Msg = getKey(request.Code, request.Machine) //获取密码
			result, _ := json.Marshal(response)
			log("result:客户使用部分机器可用邀请码请求成功。", MSG_LOG)
			return result
		} else {
			return getResultObject(VERIFY_PROTOSIGN, -4) //未指定主机
		}
	case 21: //只对特定机器（保存本地验证凭据）
		if mallow && code.Machine == request.Machine {
			response.MsgType = 2                                 //成功回复2
			response.Msg = getKey(request.Code, request.Machine) //获取密码
			result, _ := json.Marshal(response)
			log("result:客户使用特定机器邀请码请求成功。", MSG_LOG)
			return result
		} else {
			log("result:客户端使用不被允许的主机，操作失败。", MSG_LOG)
			return getResultObject(VERIFY_PROTOSIGN, -3) //该注册码已有指定主机，为非法注册
		}
	case 22: //对全部机器可用（保存本地验证凭据）
		if mallow {
			response.MsgType = 2                                 //成功回复2
			response.Msg = getKey(request.Code, request.Machine) //获取密码
			result, _ := json.Marshal(response)
			log("result:客户使用部分机器可用邀请码请求成功。", MSG_LOG)
			return result
		} else {
			return getResultObject(VERIFY_PROTOSIGN, -4) //未指定主机
		}
	}
	return nil
}

//获取包含结果的对象
func getResultObject(ProtoSign int, state int) (result []byte) {
	var response Response
	response.MsgType = state
	response.ProtoSign = ProtoSign
	switch ProtoSign {
	case VERIFY_PROTOSIGN: //邀请码验证
		switch state {
		case -1:
			response.Msg = "非法用户！请勿恶意尝试！" //协议头错误
		case -2:
			response.Msg = "无效邀请码"
		case -3:
			response.Msg = "使用非法邀请码，该邀请码已在指定主机使用"
		case -4:
			response.Msg = "未指定主机"
		case -5:
			response.Msg = "未指定邀请码"
		}
		result, _ = json.Marshal(response)
		log("resultMsgType:"+strconv.Itoa(response.MsgType), MSG_LOG)
		log("resultMsg:"+response.Msg, MSG_LOG)
	case MANAGER_PROTOSIGN: //邀请码管理
		switch state {
		case -1: //添加操作时已经包含该邀请码
			response.Msg = "该邀请码已存在！"
		case -11: //添加失败
			response.Msg = "添加邀请码失败！"
		case 1: //添加成功
			response.Msg = "添加邀请码成功！"
		case 2:
			response.Msg = "修改成功！"
		case -2:
			response.Msg = "不存在指定邀请码！"
		case -22:
			response.Msg = "修改失败！"
		case 3:
			response.Msg = "删除成功！"
		case -33:
			response.Msg = "删除失败！"
		case -3:
			response.Msg = "不存在指定邀请码！"
		}
		result, _ = json.Marshal(response)
		log("resultMsgType:"+strconv.Itoa(response.MsgType), MSG_MANAGER_LOG)
		log("resultMsg:"+response.Msg, MSG_MANAGER_LOG)
	}
	return result
}

//处理密码算法
func getKey(code string, machine string) string {
	tmp := getTimeString() + code + machine
	//for c, _ := range tmp {

	//}
	return tmp
}

//从验证码类型获取文本表示
func getVerifyType(valid int) string {
	switch valid {
	case -1:
		return "不可用邀请码"
	case 1:
		return "对特定机器可用邀请码"
	case 2:
		return "对部分机器可用邀请码"
	case 20:
		return "该验证后将钥匙保存到本地"
	case 21:
		return "对特定机器可用，且验证后将钥匙保存到本地"
	case 22:
		return "对部分机器可用，且验证后将钥匙保存到本地"
	default:
		return "无效类型"
	}
}
