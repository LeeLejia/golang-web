package app

import (
	"fmt"
	"../src/common"
	"../src/common/logger"
	"../src/model"
	"net/http"
	"strconv"
	"strings"
)

/**
     * 获取员工权限配置信息
	 * @param auth 权限 1-超级管理员 2-管理员 3-医生 4-护士
     * @return Object 如果有员工权限配置信息那么返回相应的列表，否则返回空
     */
func ListEmployee2Auth(w http.ResponseWriter, r *http.Request) {
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	auth := r.FormValue("auth")
	if auth == "" {
		common.ReturnEFormat(w, 500, "请选择权限")
		return
	}
	limit := common.PageParams(r)
	where := fmt.Sprintf("WHERE auth=%s", auth)
	auths, err := model.FindEmployee2auths(where, "", "")
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	employeeIDs := []string{}
	for _, x := range auths {
		employeeIDs = append(employeeIDs, fmt.Sprintf("%d", x.EmployeeID))
	}
	fmt.Println(fmt.Sprintf("WHERE id IN (%s)", strings.Join(employeeIDs, ",")))
	result := map[string]interface{}{}
	if len(employeeIDs) > 0 {
		employees, err := model.FindEmployees(fmt.Sprintf("WHERE id IN (%s)", strings.Join(employeeIDs, ",")), limit, "")
		if err != nil {
			common.ReturnEFormat(w, 500, err.Error())
			return
		}
		CreateOperateLogs("info", "查看", account, "查看员工权限列表", r)
		result["employees"] = employees
	}

	common.ReturnFormat(w, 200, result, "SUCCESS")

}

func ListEmployeeWithoutAuth(w http.ResponseWriter, r *http.Request) {
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	auth := r.FormValue("auth")
	oID := r.FormValue("oId")
	name := r.FormValue("name")
	if auth == "" {
		common.ReturnEFormat(w, 500, "请选择权限")
		return
	}
	if oID == "" {
		common.ReturnEFormat(w, 500, "机构id为空")
		return
	}
	where := fmt.Sprintf("WHERE auth=%s", auth)
	limit := common.PageParams(r)
	conditions := []string{}
	conditions = append(conditions, fmt.Sprintf("o_id=%s", oID))
	if name != "" {
		conditions = append(conditions, "(name LIKE '%"+name+"%' OR work_no LIKE '%"+name+"%' OR phone LIKE '%"+name+"%')")
	}
	auths, err := model.FindEmployee2auths(where, "", "")
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	employeeIDs := []string{}
	for _, x := range auths {
		employeeIDs = append(employeeIDs, fmt.Sprintf("%d", x.EmployeeID))
	}
	result := map[string]interface{}{}
	if len(employeeIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("id NOT IN (%s)", strings.Join(employeeIDs, ",")))
		condition := strings.Join(conditions, " AND ")
		fmt.Println(condition)
		employees, err := model.FindEmployees(" WHERE "+condition, limit, "")
		if err != nil {
			common.ReturnEFormat(w, 500, err.Error())
			return
		}
		result["employees"] = employees
	} else {
		condition := strings.Join(conditions, " AND ")
		employees, err := model.FindEmployees(" WHERE "+condition, limit, "")
		if err != nil {
			common.ReturnEFormat(w, 500, err.Error())
			return
		}
		result["employees"] = employees
	}
	CreateOperateLogs("info", "查看", account, "查看员工权限列表", r)
	common.ReturnFormat(w, 200, result, "SUCCESS")

}

func AddEmployee2Auth(w http.ResponseWriter, r *http.Request) {
	authStr := r.PostFormValue("auth")
	employeeIDs := r.PostFormValue("employeeIds")
	method := r.PostFormValue("method")
	if authStr == "" {
		common.ReturnEFormat(w, 500, "请选择权限")
		return
	}
	if employeeIDs == "" {
		common.ReturnEFormat(w, 500, "请选择员工")
		return
	}
	if method == "" {
		common.ReturnEFormat(w, 500, "操作为空")
		return
	}
	auth, err := strconv.ParseInt(authStr, 10, 64)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	auths, err := model.FindEmployee2auths(fmt.Sprintf("WHERE auth=%d", auth), "", "")
	if err != nil {
		common.ReturnEFormat(w, 500, "查找权限失败")
		return
	}
	if method == "add" {
		employeeID := strings.Split(employeeIDs, ",")
		for _, x := range employeeID {
			id, err := strconv.ParseInt(x, 10, 64)
			if err == nil {
				check := 0
				for _, auth := range auths {
					if auth.EmployeeID == id {
						check = 1
						break
					}
				}
				if check == 0 {
					employee2auth := model.Employee2auth{}
					employee2auth.EmployeeID = id
					employee2auth.Auth = auth
					err = employee2auth.Insert()
					if err != nil {
						Log.E("employee2auths.go", err.Error())
					}
				}
			}
		}
	} else if method == "delete" {
		employeeID := strings.Split(employeeIDs, ",")
		if len(employeeID) > 0 {
			err = model.DeleteEmployee2auths(fmt.Sprintf("WHERE auth=%d AND employee_id IN (%s)", auth, employeeIDs))
			if err != nil {
				common.ReturnEFormat(w, 500, err.Error())
				return
			}
		}
	}
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}
