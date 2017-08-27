package app

import (
	"fmt"
	"../common"
	"../common/logger"
	"../model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	account := r.PostFormValue("account")
	pwd := r.PostFormValue("pwd")
	roleType := r.PostFormValue("roleType")
	oID := r.PostFormValue("oId")
	osType := r.FormValue("osType")
	if account == "" {
		common.ReturnEFormat(w, 500, "账号不能为空")
		return
	}
	if pwd == "" {
		common.ReturnEFormat(w, 500, "请输入密码")
		return
	}
	if roleType == "user" || roleType == "admin" || roleType == "super" {
		where := fmt.Sprintf("WHERE account='%s' AND pwd='%s'", account, common.MD5Password(pwd))
		fmt.Println(where)
		user, err := model.FindUsers(where, "", "")
		if err != nil {
			common.ReturnEFormat(w, 503, err.Error())
			return
		}
		if len(user) == 0 {
			common.ReturnEFormat(w, 404, "用户名或密码错误")
			return
		}
		if user[0].State == model.UserStateOff {
			common.ReturnEFormat(w, 500, "账号已被停用")
			return
		}

		token := common.SaveSession(user[0].Uid, account, osType, user[0].UserType)
		result := map[string]interface{}{"user": user, "token": token}
		if user[0].UserType == "admin" {
			organization, err := model.FindOrganization("WHERE code='" + user[0].Account + "'")
			if err == nil {
				result["oId"] = organization.ID
			}
		}
		CreateOperateLogs("info", "登录", account, "本地登录，登录结果：登录成功", r)
		common.ReturnFormat(w, 200, result, "SUCCESS")
	} else if roleType == "employee" {
		where := fmt.Sprintf("WHERE work_no='%s' AND pwd='%s'", account, common.MD5Password(pwd))
		fmt.Println(where)
		user, err := model.FindEmployees(where, "", "")
		if err != nil {
			common.ReturnEFormat(w, 503, err.Error())
			return
		}
		if len(user) == 0 {
			common.ReturnEFormat(w, 404, "用户名或密码错误")
			return
		}
		if user[0].State == model.UserStateOff {
			common.ReturnEFormat(w, 500, "账号已被停用")
			return
		}
		if fmt.Sprintf("%d", user[0].OID) != oID {
			common.ReturnEFormat(w, 500, "机构选择错误")
			return
		}
		token := common.SaveSession(user[0].ID, account, osType, roleType)
		CreateOperateLogs("info", "登录", account, "本地登录，登录结果：登录成功", r)
		common.ReturnFormat(w, 200, map[string]interface{}{"user": user, "token": token, "oId": oID}, "SUCCESS")
	} else {
		common.ReturnEFormat(w, 500, "未识别的用户角色")
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	account := r.PostFormValue("account")
	pwd := r.PostFormValue("pwd")
	//username := r.PostFormValue("username")
	name := r.PostFormValue("name")
	if account == "" {
		common.ReturnEFormat(w, 500, "手机号不能为空")
		return
	}
	if !common.IsPhone(account) {
		common.ReturnEFormat(w, 500, "手机号格式错误")
		return
	}
	if pwd == "" {
		common.ReturnEFormat(w, 500, "请输入密码")
		return
	}
	where := fmt.Sprintf("WHERE account='%s'", account)
	users, err := model.FindUsers(where, "", "")
	if len(users) > 0 {
		common.ReturnEFormat(w, 500, "账号已存在，请直接登录")
		return
	}
	user := model.User{}
	user.Account = account
	user.Pwd = common.MD5Password(pwd)
	user.State = model.UserStateOn
	user.Username = account
	user.Name = name
	user.UserNo = ""
	user.UserLevel = ""
	user.UserType = "user"
	err = user.Insert()
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	uaccount, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	accountStr := r.PostFormValue("accounts")
	pwd := r.PostFormValue("pwd")
	//username := r.PostFormValue("username")
	state := r.PostFormValue("state")
	//userType := r.PostFormValue("type")
	password := common.MD5Password(pwd)
	accounts := strings.Split(accountStr, ",")
	for _, account := range accounts {
		update := []string{}
		if account == "" {
			common.ReturnEFormat(w, 500, "手机号不能为空")
			return
		}
		where := fmt.Sprintf("WHERE account='%s'", account)
		users, _ := model.FindUsers(where, "", "")
		if len(users) <= 0 {
			common.ReturnEFormat(w, 404, "用户不存在")
			return
		}
		user := users[0]
		if user.Pwd != password && pwd != "" {
			update = append(update, fmt.Sprintf("pwd='%s'", password))
		}
		if state != "" {
			stateInt, err := strconv.Atoi(state)
			if err == nil {
				if user.State != stateInt {
					update = append(update, fmt.Sprintf("state=%d", stateInt))
				}
			}
		}
		if len(update) > 0 {
			updater := strings.Join(update, ",")
			fmt.Println(updater)
			fmt.Println(where)
			if state != "" && user.UserType == "admin" {
				organization, err := model.FindOrganization("WHERE user_account='" + account + "'")
				if err == nil {
					err = model.UpdateEmployees(fmt.Sprintf("state=%s", state), fmt.Sprintf("WHERE o_id=%d", organization.ID))
					if err != nil {
						CreateOperateLogs("info", "修改", uaccount, fmt.Sprintf("修改会员资料，用户账号：%s，修改结果：修改失败，失败原因：%s", account, err.Error()), r)
						logger.Error(fmt.Sprintf("更新机构员工状态失败，%s", err.Error()), "users.go")
					}
				}
			}
			err := model.UpdateUsers(updater, where)
			if err != nil {
				common.ReturnEFormat(w, 500, err.Error())
				return
			}
			CreateOperateLogs("info", "修改", uaccount, fmt.Sprintf("修改会员资料，用户账号：%s，修改结果：修改成功", account), r)
		}
	}
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	q := r.FormValue("q")
	state := r.FormValue("state")
	order := r.FormValue("order")
	limit := common.PageParams(r)
	stateInt, err := strconv.Atoi(state)
	conditions := []string{}
	condition := "WHERE user_type='user' "
	if err == nil {
		conditions = append(conditions, fmt.Sprintf("state=%d", stateInt))
	}
	if q != "" {
		conditions = append(conditions, "(username LIKE '%"+q+"%' OR account LIKE '%"+q+"%')")
	}
	if len(conditions) > 0 {
		condition = condition + " AND " + strings.Join(conditions, " AND ")
	}
	if order != "" {
		if strings.Contains(order, "DESC") || strings.Contains(order, "ASC") {
			order = " ORDER BY " + order
		}
	}
	fmt.Println(condition)
	fmt.Println(limit)
	users, err := model.FindUsers(condition, limit, order)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	total, err := model.CountUsers(condition)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}

	type userJSON struct {
		Uid       int64  `json:"id"`
		Username  string `json:"username"`
		Name      string `json:"name"`
		UserNo    string `json:"userNo"`
		UserLevel string `json:"userLevel"`
		Account   string `json:"account"`
		State     int    `json:"state"`
		UserType  string `json:"userType"`

		UpdatedAt time.Time `json:"updatedAt"`
		CreatedAt time.Time `json:"createdAt"`
	}

	result := []userJSON{}
	for _, x := range users {
		tmp := userJSON{}
		tmp.Uid = x.Uid
		tmp.Username = x.Username
		tmp.Name = x.Name
		tmp.UserNo = x.UserNo
		tmp.UserLevel = x.UserLevel
		tmp.Account = x.Account
		tmp.State = x.State
		tmp.UserType = x.UserType
		tmp.UpdatedAt = x.UpdatedAt
		tmp.CreatedAt = x.CreatedAt
		result = append(result, tmp)
	}
	CreateOperateLogs("info", "查看", account, "查看会员列表", r)
	common.ReturnFormat(w, 200, map[string]interface{}{"users": result, "total": total}, "SUCCESS")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	_, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	common.RemoveSession(r)
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}
func generateUser(account, pwd, username, userType string) (user *model.User, err error) {
	user = new(model.User)
	user.Account = account
	user.Pwd = common.MD5Password(pwd)
	user.State = model.UserStateOn
	user.Username = username
	user.Name = username
	user.UserType = userType
	err = user.Insert()
	return
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	accountStr := r.PostFormValue("account")
	pwd := r.PostFormValue("pwd")
	code := r.PostFormValue("code")
	if !CheckMobileCode(accountStr, code) {
		common.ReturnEFormat(w, 500, "验证码错误")
		return
	}
	updater := fmt.Sprintf("pwd='%s'", common.MD5Password(pwd))
	where := fmt.Sprintf("WHERE account='%s'", accountStr)
	users, _ := model.FindUsers(where, "", "")
	if len(users) <= 0 {
		common.ReturnEFormat(w, 404, "用户不存在")
		return
	}
	err := model.UpdateUsers(updater, where)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func CheckMobileCode(accountStr, code string) bool {
	//TODO 短信验证代码
	return true
}

func Test(w http.ResponseWriter, r *http.Request){
	phone := r.PostFormValue("phone")
	common.SendSMSCode(phone)
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}