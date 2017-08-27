package app

import (
	"fmt"
	"../common"
	"../model"
	"net/http"
	"strconv"
	"strings"
)

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	_, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	oIDStr := r.PostFormValue("oId")
	name := r.PostFormValue("name")
	phone := r.PostFormValue("phone")
	sex := r.PostFormValue("sex")
	title := r.PostFormValue("title")
	workNo := r.PostFormValue("workNo")
	manageNumStr := r.PostFormValue("manageNum")
	employeeTypeStr := r.PostFormValue("employeeType")
	stateStr := r.PostFormValue("state")
	workGroupIDStr := r.PostFormValue("workGroupId")
	//auths := r.PostFormValue("auths")
	if oIDStr == "" {
		common.ReturnEFormat(w, 500, "机构id不能为空")
		return
	}
	oID, err := strconv.ParseInt(oIDStr, 10, 64)
	if err != nil {
		common.ReturnEFormat(w, 500, "机构id格式错误")
		return
	}

	if name == "" {
		common.ReturnEFormat(w, 500, "姓名不能为空")
		return
	}
	if phone == "" {
		common.ReturnEFormat(w, 500, "手机号不能为空")
		return
	}
	if workNo == "" {
		common.ReturnEFormat(w, 500, "工号不能为空")
		return
	}
	where := fmt.Sprintf("WHERE phone='%s'", phone)

	count, err := model.CountEmployees(where)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	if count > 0 {
		common.ReturnEFormat(w, 500, "手机号已存在")
		return
	}

	employee := model.Employee{}
	employee.Name = name
	employee.Phone = phone
	employee.Sex = sex
	employee.Title = title
	employee.WorkNo = workNo
	employee.Pwd = common.MD5Password("123456")
	employee.State = 0
	employee.OID = oID

	manageNum, err := strconv.ParseInt(manageNumStr, 10, 64)
	if err == nil {
		employee.ManageNum = manageNum
	}
	employeeType, err := strconv.ParseInt(employeeTypeStr, 10, 64)
	if err == nil {
		employee.EmployeeType = employeeType
	}
	state, err := strconv.Atoi(stateStr)
	if err == nil {
		employee.State = state
	}
	if workGroupIDStr != "" {
		workGroupID, err := strconv.ParseInt(workGroupIDStr, 10, 64)
		if err == nil {
			employee.WorkGroupID = workGroupID
		}
	}

	err = employee.Insert()
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}

	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	_, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	accountStr := r.PostFormValue("accounts")
	pwd := r.PostFormValue("pwd")
	state := r.PostFormValue("state")
	update := []string{}
	password := common.MD5Password(pwd)
	accounts := strings.Split(accountStr, ",")
	for _, account := range accounts {
		if account == "" {
			common.ReturnEFormat(w, 500, "手机号不能为空")
			return
		}
		where := fmt.Sprintf("WHERE phone='%s'", account)
		users, _ := model.FindEmployees(where, "", "")
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
		//if userType != "" {
		//	typeInt, err := strconv.Atoi(userType)
		//	if err == nil {
		//		if user.Type != typeInt {
		//			update = append(update, fmt.Sprintf("type=%d", typeInt))
		//		}
		//	}
		//}
		if len(update) > 0 {
			updater := strings.Join(update, ",")
			fmt.Println(updater)
			err := model.UpdateEmployees(updater, where)
			if err != nil {
				common.ReturnEFormat(w, 500, err.Error())
				return
			}
		}
	}
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func ListEmployees(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	q := r.FormValue("q")
	state := r.FormValue("state")
	order := r.FormValue("order")
	oID := r.FormValue("oId")
	workGroupID := r.FormValue("workGroupId")
	limit := common.PageParams(r)
	stateInt, err := strconv.Atoi(state)
	conditions := []string{}
	condition := ""
	if oID == "" {
		common.ReturnEFormat(w, 500, "机构id不能为空")
		return
	}
	if err == nil {
		conditions = append(conditions, fmt.Sprintf("state=%d", stateInt))
	}
	if workGroupID != "" {
		conditions = append(conditions, fmt.Sprintf("work_group_id=%s", workGroupID))
	}
	if q != "" {
		conditions = append(conditions, "(name LIKE '"+q+"' OR work_no LIKE '%"+q+"%' OR phone LIKE '%"+q+"%')")
	}
	conditions = append(conditions, "o_id="+oID)
	if len(conditions) > 0 {
		condition = " WHERE " + strings.Join(conditions, " AND ")
	}
	if order != "" {
		order = " ORDER BY " + order
	}
	fmt.Println(condition)
	fmt.Println(limit)
	employees, err := model.FindEmployees(condition, limit, order)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	total, err := model.CountEmployees(condition)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	CreateOperateLogs("info", "查看", account, "查看员工列表", r)
	common.ReturnFormat(w, 200, map[string]interface{}{"employees": employees, "total": total}, "SUCCESS")
}

func DeleteEmployees(w http.ResponseWriter, r *http.Request) {
	employeeIDs := r.PostFormValue("employeeIds")
	if employeeIDs == "" {
		common.ReturnEFormat(w, 500, "员工id不能为空")
		return
	}

	condition := fmt.Sprintf("WHERE id IN (%s)", employeeIDs)
	err := model.DeleteEmployees(condition)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}
