package app

import (
	"net/http"
	"../common"
	"../common/log"
	"../model"
	//
	//"fmt"
	"strconv"
	"github.com/bitly/go-simplejson"
	"time"
)



func Publish(sess *common.Session, w http.ResponseWriter, r *http.Request){

	if !common.IsRole(sess.User.Role,model.USER_ROLE_EMPLOYER){
		log.W(common.ACTION_VIOLENCE,sess.User.Email,"该用户roles=%s,尝试发布任务.",sess.User.Role)
		common.ReturnEFormat(w,common.CODE_ROLE_INVADE, "抱歉,你不并是雇主,请先到个人信息页面修改角色.")
		return
	}
	r.ParseForm()
	name := r.PostFormValue("name")
	describe := r.PostFormValue("describe")
	commission := r.PostFormValue("commission")
	money_lower,_:=strconv.Atoi(r.PostFormValue("money_lower"))
	money_upper,_:=strconv.Atoi(r.PostFormValue("money_upper"))
	outsourcing := false
	if r.PostFormValue("outsourcing")=="true"{
		outsourcing = true
	}
	types := r.PostFormValue("type")
	code := false
	if r.PostFormValue("code")=="true"{
		code = true
	}
	annex,_ := simplejson.NewJson([]byte(r.PostFormValue("annex")))
	from_time,_ := strconv.ParseInt(r.PostFormValue("from_time"),10,64)
	to_time,_ := strconv.ParseInt(r.PostFormValue("to_time"),10,64)
	if len(name)<3{
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "项目名不能少于三个字符")
		return
	}
	if len(describe)<10{
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "项目描述不能少于10个字符")
		return
	}
	if len(name)<3{
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "项目名不能少于三个字符")
		return
	}
	if len(name)<3{
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "项目名不能少于三个字符")
		return
	}
	if len(name)<3{
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "项目名不能少于三个字符")
		return
	}
	if from_time>to_time{
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "项目终止时间不能大于项目开始时间")
		return
	}
	if money_upper<money_lower{
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "请正确设置预算上下限")
		return
	}
	if money_lower<1{
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "一块钱都不给,你在逗我??")
		return
	}
	publish := model.T_publish{
		Owner:       sess.User.Email,
		Name:        name,
		Describe:    describe,
		MoneyLow:    money_lower,
		MoneyUp:     money_upper,
		OutSourcing: outsourcing,
		Labels:        types,
		Commission:  commission,
		NeedCode:    code,
		Annex:       annex,
		FromTime:    time.Unix(from_time, 0),
		ToTime:      time.Unix(to_time, 0),
		UpdateTime:  time.Now(),
		CreatedTime: time.Now(),
	}
	err:=PublishModel.Insert(publish)
	if err != nil {
		log.E("publish失败",publish.Owner,"projectName:%s,describe:%s,money_lower:%d,money_upper:%d,reason:%s",name,describe,money_lower,money_upper,err.Error())
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "恭喜你找到个bug..创建新用户出错,劳烦转告下管理员!")
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"msg":"注册成功，请重新登录！"})
	log.N("publish成功",publish.Owner,"projectName:%s,money_lower:%d,money_upper:%d",name,money_lower,money_upper)
}