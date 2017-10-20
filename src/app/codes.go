package app

import (
	"net/http"
	"../model"
	"../common"
	"fmt"
	"time"
	"strconv"
	"github.com/bitly/go-simplejson"
)

/**
添加邀请码
 */
func AddCode(w http.ResponseWriter, r *http.Request, user *model.T_user){
	if user.Role!=model.USER_ROLE_ADMIN && user.Role!=model.USER_ROLE_DEVELOPER &&user.Role!=model.USER_ROLE_SUPER{
		common.ReturnEFormat(w,common.CODE_ROLE_INVADE,"你没有添加邀请码的权限！")
		return
	}
	appId:=r.PostFormValue("appId")
	if appId==""{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"未指定App！")
		return
	}
	consumer_:=r.PostFormValue("consumer")
	consumer,err:=simplejson.NewJson([]byte(consumer_))
	if err!=nil{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"客户信息格式不合法！")
		return
	}
	describe:=r.PostFormValue("describe")
	valid_:=r.PostFormValue("valid")
	valid:=true
	if valid_ =="" || valid_ =="false"{
		valid=false
	}
	mostCount_:=r.PostFormValue("mostCount")
	mostCount,_:=strconv.Atoi(mostCount_)
	EnableTime_:=r.PostFormValue("enableTime")
	EnableTime:=true
	if EnableTime_ == "" || EnableTime_ == "false" {
		EnableTime=false
	}
	startTime_:=r.PostFormValue("startTime")
	endTime_ :=r.PostFormValue("endTime")
	startTime,_:=strconv.ParseInt(startTime_,10,64)
	endTime,_:=strconv.ParseInt(endTime_,10,64)
	if endTime!=0 && startTime>=endTime{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"邀请码生效时间不能超过终止时间！")
		return
	}
	code:=common.GetRandomString(6)
	developer:=user.Id

	t_code:=model.T_code{
		Code:code,
		AppId:appId,
		Developer:developer,
		 Consumer:consumer,
		Describe:describe,
		Valid:valid,
		MachineCount:0,
		MostCount:mostCount,
		EnableTime:EnableTime,
		 StartTime:time.Unix(startTime,0),
		 EndTime:time.Unix(endTime,0),
		CreatedAt:time.Now(),
	}
	err=t_code.Insert()
	if err==nil{
		common.ReturnEFormat(w,common.CODE_DB_RW_ERR,"系统写入出错，请稍后重试！")
		return
	}
	common.ReturnFormat(w,200,nil,"添加成功！")
}
/**
获取邀请码
 */
func ListCodes(w http.ResponseWriter, r *http.Request, user *model.T_user){
	if user.Role!=model.USER_ROLE_DEVELOPER && user.Role!=model.USER_ROLE_SUPER && user.Role!=model.USER_ROLE_ADMIN {
		common.ReturnEFormat(w, common.CODE_ROLE_INVADE,"你没有获取邀请码列表的权限！")
		return
	}
	cond:=""
	if user.Role==model.USER_ROLE_DEVELOPER{
		cond=fmt.Sprintf("WHERE DEVELOPER=%d",user.Id)
	}
	codes,err:= model.FindCodes(cond,"","")
	if err!=nil{
		fmt.Println(err.Error())
	}
	common.ReturnFormat(w, 200, codes, "success")
}
/**
删除邀请码
 */
func DeleteCodes(w http.ResponseWriter, r *http.Request, user *model.T_user){

}
/**
更新邀请码
 */
func UpdateCodes(w http.ResponseWriter, r *http.Request, user *model.T_user){

}