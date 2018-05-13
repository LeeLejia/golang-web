package app

import (
	"net/http"
	"../common"
	"../common/log"
	"../model"
	"strconv"
	"github.com/bitly/go-simplejson"
	fm "github.com/cjwddz/fast-model"
	"time"
	"fmt"
)

type Task struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Describe    string    `json:"describe"`
	MoneyLow    int       `json:"money_lower"`
	MoneyUp     int       `json:"money_upper"`
	//OutSourcing bool      `json:"outsourcing"`
	Labels      string    `json:"labels"`
	//Commission  string    `json:"commission"`
	NeedCode    bool      `json:"need_code"`
	Annex       string    `json:"annex"`
	FromTime    int64 `json:"from_time"`
	ToTime      int64 `json:"to_time"`
	UpdateTime  int64 `json:"update_at"`
	CreatedTime int64 `json:"created_at"`
}

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
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"msg":"发布任务成功!"})
	log.N("publish成功",publish.Owner,"projectName:%s,money_lower:%d,money_upper:%d",name,money_lower,money_upper)
}

func GetTask(_ *common.Session, w http.ResponseWriter, r *http.Request) {
	cond:=fm.DbCondition{}
	total,err:=PublishModel.Count(cond)
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("GetTask","",err.Error())
		return
	}
	result,err:=PublishModel.Query(cond.Limit2(r,"start","count").Order("order by created_at desc"))
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("GetTask","",err.Error())
		return
	}
	tasks:= make([]Task,len(result))
	for i,item:=range result{
		task:=item.(model.T_publish)
		annex,_:=task.Annex.Encode()
		tasks[i] = Task{
			ID:          task.ID,
			Name:        task.Name,
			Describe:    task.Describe,
			MoneyLow:    task.MoneyLow,
			MoneyUp:     task.MoneyUp,
			Labels:      task.Labels,
			NeedCode:    task.NeedCode,
			Annex:       string(annex),
			FromTime:    task.FromTime.Unix(),
			ToTime:      task.ToTime.Unix(),
			UpdateTime:  task.UpdateTime.Unix(),
			CreatedTime: task.CreatedTime.Unix(),
		}
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"tasks": tasks,"total":total})
	log.N("GetTask","",fmt.Sprintf("listCount=%d,total=%d",len(result),total))
	return
}