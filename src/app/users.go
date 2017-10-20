package app

import (
	"fmt"
	"../common"
	"../common/log"
	"../model"
	"net/http"
	"github.com/bitly/go-simplejson"
)

/**
登录
 */
func Login(w http.ResponseWriter, r *http.Request) {
	account := r.PostFormValue("account")
	pwd := r.PostFormValue("pwd")
	osType := r.PostFormValue("osType")
	if account == "" {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "账号不能为空")
		return
	}
	if pwd == "" {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "请输入密码")
		return
	}

	where := fmt.Sprintf("WHERE  pwd='%s' AND (nick='%s' OR phone='%s' OR email='%s')", pwd, account,account,account)
	fmt.Println(log.Green(where))
	user, err := model.FindUsers(where, "", "")
	if err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, err.Error())
		return
	}
	if len(user) == 0 {
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "用户名或密码错误")
		return
	}

	if user[0].Status == model.USER_STATUS_INVALID {
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "账号已被停用")
		return
	}
	sessionId,session := common.SaveSession(user[0],osType)
	user_tmp:=simplejson.New()
	user_tmp.Set("role", user[0].Role)
	user_tmp.Set("nick",user[0].Nick)
	user_tmp.Set("avatar", user[0].Avatar)
	user_tmp.Set("phone", user[0].Phone)
	user_tmp.Set("email", user[0].Email)
	user_tmp.Set("qq", user[0].QQ)
	result := map[string]interface{}{"user": user_tmp, "sessionId":sessionId, "token": session.Token}
	http.SetCookie(w,&http.Cookie{Name:"sessionId",Value:sessionId,Path:"/"})
	http.SetCookie(w,&http.Cookie{Name:"token",Value:session.Token,Path:"/"})
	common.ReturnFormat(w, 200, result, "success")
}
/**
用户注销
 */
func Logout(w http.ResponseWriter, r *http.Request) {
	common.RemoveSession(r)
	common.ReturnFormat(w, 200, nil, "success")
}
/**
注册
 */
func Register(w http.ResponseWriter, r *http.Request) {
	phone := r.PostFormValue("phone")
	email := r.PostFormValue("email")
	pwd := r.PostFormValue("pwd")
	role := r.PostFormValue("role")
	if phone == "" {
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "手机号不能为空")
		return
	}
	if !common.IsPhone(phone) {
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "手机号格式错误")
		return
	}
	if email == ""{
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "邮箱不能为空")
		return
	}
	if !common.IsEmail(email){
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "邮箱格式错误")
		return
	}
	if pwd == "" {
		common.ReturnEFormat(w,common.CODE_VERIFY_FAIL, "请输入密码")
		return
	}
	if len(pwd) < 6 {
		common.ReturnEFormat(w,common.CODE_VERIFY_FAIL, "请输入至少6位的密码")
		return
	}
	if role!= model.USER_ROLE_COMMON && role!=model.USER_ROLE_DEVELOPER{
		common.ReturnEFormat(w,common.CODE_VERIFY_FAIL, "请选择正确的账号类型")
		return
	}
	where := fmt.Sprintf("WHERE phone='%s' or email='%s'", phone, email)
	users, err := model.FindUsers(where, "", "")
	if len(users) > 0 {
		common.ReturnEFormat(w, 500, "账号已存在，请直接登录")
		return
	}
	user := &model.T_user{
		Role:role,
		Nick:"新手",
		Pwd:pwd,
		Status:model.USER_STATUS_VALID,
		Avatar:"https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
		Phone:phone,
		Email:email,
		QQ:"",
		Expend: simplejson.New(),
	}
	err = user.Insert()
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	common.ReturnFormat(w, 200, nil, "注册成功，请重新登录！")
}

/**
修改
todo 记得同步session中的用户
 */
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	//uaccount, err := common.CheckSession(r)
	//if err != nil {
	//	common.ReturnEFormat(w, 403, err.Error())
	//	return
	//}
	//accountStr := r.PostFormValue("accounts")
	//pwd := r.PostFormValue("pwd")
	////username := r.PostFormValue("username")
	//state := r.PostFormValue("state")
	////userType := r.PostFormValue("type")
	//password := common.MD5Password(pwd)
	//accounts := strings.Split(accountStr, ",")
	//for _, account := range accounts {
	//	update := []string{}
	//	if account == "" {
	//		common.ReturnEFormat(w, 500, "手机号不能为空")
	//		return
	//	}
	//	where := fmt.Sprintf("WHERE account='%s'", account)
	//	users, _ := model.FindUsers(where, "", "")
	//	if len(users) <= 0 {
	//		common.ReturnEFormat(w, 404, "用户不存在")
	//		return
	//	}
	//	user := users[0]
	//	if user.Pwd != password && pwd != "" {
	//		update = append(update, fmt.Sprintf("pwd='%s'", password))
	//	}
	//	if state != "" {
	//		stateInt, err := strconv.Atoi(state)
	//		if err == nil {
	//			if user.State != stateInt {
	//				update = append(update, fmt.Sprintf("state=%d", stateInt))
	//			}
	//		}
	//	}
	//	if len(update) > 0 {
	//		updater := strings.Join(update, ",")
	//		fmt.Println(updater)
	//		fmt.Println(where)
	//		if state != "" && user.UserType == "admin" {
	//			//organization, err := model.FindOrganization("WHERE user_account='" + account + "'")
	//			//if err == nil {
	//			//	err = model.UpdateEmployees(fmt.Sprintf("state=%s", state), fmt.Sprintf("WHERE o_id=%d", organization.ID))
	//			//	if err != nil {
	//			//		CreateOperateLogs("info", "修改", uaccount, fmt.Sprintf("修改会员资料，用户账号：%s，修改结果：修改失败，失败原因：%s", account, err.Error()), r)
	//			//		logger.Error(fmt.Sprintf("更新机构员工状态失败，%s", err.Error()), "users.go")
	//			//	}
	//			//}
	//		}
	//		err := model.UpdateUsers(updater, where)
	//		if err != nil {
	//			common.ReturnEFormat(w, 500, err.Error())
	//			return
	//		}
	//	}
	//}
	//common.ReturnFormat(w, 200, nil, "SUCCESS")
}
/**
列出用户
 */
func ListUsers(w http.ResponseWriter, r *http.Request) {
	//account, err := common.CheckSession(r)
	//if err != nil {
	//	common.ReturnEFormat(w, 403, err.Error())
	//	return
	//}
	//q := r.FormValue("q")
	//state := r.FormValue("state")
	//order := r.FormValue("order")
	//limit := common.PageParams(r)
	//stateInt, err := strconv.Atoi(state)
	//conditions := []string{}
	//condition := "WHERE user_type='user' "
	//if err == nil {
	//	conditions = append(conditions, fmt.Sprintf("state=%d", stateInt))
	//}
	//if q != "" {
	//	conditions = append(conditions, "(username LIKE '%"+q+"%' OR account LIKE '%"+q+"%')")
	//}
	//if len(conditions) > 0 {
	//	condition = condition + " AND " + strings.Join(conditions, " AND ")
	//}
	//if order != "" {
	//	if strings.Contains(order, "DESC") || strings.Contains(order, "ASC") {
	//		order = " ORDER BY " + order
	//	}
	//}
	//fmt.Println(condition)
	//fmt.Println(limit)
	//users, err := model.FindUsers(condition, limit, order)
	//if err != nil {
	//	common.ReturnEFormat(w, 503, err.Error())
	//	return
	//}
	//total, err := model.CountUsers(condition)
	//if err != nil {
	//	common.ReturnEFormat(w, 503, err.Error())
	//	return
	//}
	//
	//type userJSON struct {
	//	Uid       int64  `json:"id"`
	//	Username  string `json:"username"`
	//	Name      string `json:"name"`
	//	UserNo    string `json:"userNo"`
	//	UserLevel string `json:"userLevel"`
	//	Account   string `json:"account"`
	//	State     int    `json:"state"`
	//	UserType  string `json:"userType"`
	//
	//	UpdatedAt time.Time `json:"updatedAt"`
	//	CreatedAt time.Time `json:"createdAt"`
	//}
	//
	//result := []userJSON{}
	//for _, x := range users {
	//	tmp := userJSON{}
	//	tmp.Uid = x.
	//	tmp.Username = x.Username
	//	tmp.Name = x.Name
	//	tmp.UserNo = x.UserNo
	//	tmp.UserLevel = x.UserLevel
	//	tmp.Account = x.Account
	//	tmp.State = x.State
	//	tmp.UserType = x.UserType
	//	tmp.UpdatedAt = x.UpdatedAt
	//	tmp.CreatedAt = x.CreatedAt
	//	result = append(result, tmp)
	//}
	////CreateOperateLogs("info", "查看", account, "查看会员列表", r)
	//common.ReturnFormat(w, 200, map[string]interface{}{"users": result, "total": total}, "SUCCESS")
}

func generateUser(account, pwd, username, userType string) (user *model.T_user, err error) {
	user = &model.T_user{
		Role:model.USER_ROLE_ADMIN,
		Nick:"白菜",
		Pwd:"imjia123",
		Status:2,
		Avatar:"https://avatars2.githubusercontent.com/u/24471738?v=4&s=40",
		Phone:"13480332034",
		Email:"cjwddz@qq.com",
		QQ:"1436983000",
		Expend: simplejson.New(),
	}
	err = user.Insert()
	return
}

func CheckMobileCode(accountStr, code string) bool {
	//TODO 短信验证代码
	return true
}