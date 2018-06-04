package app

import (
	"../common"
	"../common/log"
	"../model"
	fm "github.com/cjwddz/fast-model"
	"net/http"
	"strings"
	"time"
)

/**
登录
 */
func Login(_ *common.Session, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	account := r.PostFormValue("account")
	pwd := r.PostFormValue("pwd")
	osType := r.PostFormValue("osType")
	if account == "" {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "账号不能为空")
		return
	}
	if pwd == "" {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "密码不能为空")
		return
	}
	if osType == "" {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "终端类型缺失")
		return
	}
	log.I("login请求","","account:%s,pwd:%s,osType:%s",account,pwd,osType)
	user, err := UserModel.Query(fm.DbCondition{}.Or("=","email",account).Or("=","phone",account).Limit(1,-1))
	if err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "服务器出错！请稍后重试.")
		log.E("login出错","",err.Error())
		return
	}
	if len(user) == 0 {
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "用户不存在")
		log.W("login失败","","account:%s,pwd:%s,osType:%s,reason:%s",account,pwd,osType,"用户不存在")
		return
	}
	usr:=user[0].(model.T_user)
	if usr.Pwd!= pwd{
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "密码错误")
		log.W("login失败","","account:%s,pwd:%s,osType:%s,reason:%s",account,pwd,osType,"密码错误")
		return
	}
	if usr.Status == model.USER_STATUS_INVALID {
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "账号已被停用,请咨询平台管理员")
		log.W("login失败","","account:%s,pwd:%s,osType:%s,reason:%s",account,pwd,osType,"账号已被停用,请咨询平台管理员")
		return
	}
	if usr.Status == model.USER_STATUS_WAITING_VERIFY {
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "该帐号正在等待审核,请咨询平台管理员")
		log.W("login失败","","account:%s,pwd:%s,osType:%s,reason:%s",account,pwd,osType,"该帐号正在等待审核,请咨询平台管理员")
		return
	}
	sessionId,session := common.SaveSession(usr,osType)
	http.SetCookie(w,&http.Cookie{Name:"sessionId",Value:sessionId,Path:"/"})
	http.SetCookie(w,&http.Cookie{Name:"token",Value:session.Token,Path:"/"})
	expend,_:=usr.Expend.Encode()
	tmp:=map[string]interface{}{
		"role": usr.Role,
		"nick": usr.Nick,
		"avatar": usr.Avatar,
		"qq":usr.QQ,
		"status":usr.Status,
		"email": usr.Email,
		"phone":usr.Phone,
		"expend":string(expend),
		"update_at":usr.UpdatedAt,
		"created_at":usr.CreatedAt,
	}
	result := map[string]interface{}{"user": tmp, "sessionId":sessionId, "token": session.Token,"msg":"登录成功！"}
	common.ReturnFormat(w, common.CODE_OK, result)
	log.N("login成功","","account:%s,pwd:%s,osType:%s",account,pwd,osType)
}

/**
用户注销
 */
func Logout(sess *common.Session, w http.ResponseWriter, r *http.Request) {
	log.N("logout",sess.User.Email,"注销成功！")
	common.RemoveSession(r)
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"msg":"注销成功！"})
}

/**
注册
 */
func Register(_ *common.Session, w http.ResponseWriter, r *http.Request) {
	phone := r.PostFormValue("phone")
	email := r.PostFormValue("email")
	pwd := r.PostFormValue("password")
	role := r.PostFormValue("roles")
	nick := r.PostFormValue("nick")
	if phone == "" {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "手机号不能为空")
		return
	}
	if !common.IsPhone(phone) {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "手机号格式错误")
		return
	}
	if email == ""{
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "邮箱不能为空")
		return
	}
	if !common.IsEmail(email){
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "邮箱格式错误")
		return
	}
	if pwd == "" {
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID, "密码为空")
		return
	}
	if len(pwd) < 6 {
		common.ReturnEFormat(w,common.CODE_VERIFY_FAIL, "请输入至少6位的密码")
		return
	}
	roles:=strings.Split(role,",")
	for _,r:= range roles{
		if r != model.USER_ROLE_EMPLOYER && r!=model.USER_ROLE_ADMIN && r!=model.USER_ROLE_DEVELOPER {
			common.ReturnEFormat(w,common.CODE_VERIFY_FAIL, "不能设置非法角色")
			return
		}
		if r== model.USER_ROLE_ADMIN && len(roles)>1{
			common.ReturnEFormat(w,common.CODE_VERIFY_FAIL, "管理员角色不能和其它角色一起设置")
			return
		}
	}
	if nick == "" {
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID, "昵称不能为空!")
		return
	}
	log.I("register请求","","phone:%s,email:%s,role:%s,nick:%s",phone,email,role,nick)
	users, err := UserModel.Query(fm.DbCondition{}.Or("=","phone",phone).Or("=","email",email).Limit(1,-1))
	if err==nil && len(users) > 0 {
		user:=users[0].(model.T_user)
		if user.Phone==phone{
			common.ReturnEFormat(w, common.REGISTER_ACCOUNT_EXIST, "手机号已经被注册，不可重复注册.")
			log.I("register失败","","phone:%s,email:%s,role:%s,reason:%s",phone,email,role,"手机号已经被注册，不可重复注册.")
		}else if user.Email==email{
			common.ReturnEFormat(w, common.REGISTER_ACCOUNT_EXIST, "所使用邮箱已经被注册，不可重复注册.")
			log.I("register失败","","phone:%s,email:%s,role:%s,reason:%s",phone,email,role,"所使用邮箱已经被注册，不可重复注册.")
		}
		return
	}
	status:= model.USER_STATUS_VALID
	if common.IsRole(role,model.USER_ROLE_ADMIN){
		status = model.USER_STATUS_WAITING_VERIFY
	}
	user := model.T_user{
		Role:role,
		Nick:nick,
		Pwd:pwd,
		Status:status,
		Avatar:"https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
		Phone:phone,
		Email:email,
		QQ:"",
		UpdatedAt:time.Now(),
		CreatedAt:time.Now(),
	}
	err = UserModel.Insert(user)
	if err != nil {
		log.E("register失败","","phone:%s,email:%s,role:%s,reason:%s",phone,email,role,err.Error())
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "恭喜你找到个bug..创建新用户出错,劳烦转告下管理员!")
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"msg":"注册成功，请重新登录！"})
	log.N("register成功","","注册成功=phone:%s,email:%s,role:%s",phone,email,role)
}

/**
修改
todo 记得同步session中的用户
 */
func SetUserAvatar(sess *common.Session, w http.ResponseWriter, r *http.Request) {
	avatar := r.PostFormValue("avatar")
	if len(avatar) != 64 {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "无效的hash值")
		return
	}
	err:=UserModel.Update(fm.DbSetCondition{}.And("=","email",sess.User.Email).Set("avatar",avatar))
	if err != nil {
		log.E("setUserAvatar失败",sess.User.Email,"reason:%s",err.Error())
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "修改头像失败!")
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"msg":"修改头像成功!"})
	log.N("setUserAvatar成功",sess.User.Email,"修改头像")
}