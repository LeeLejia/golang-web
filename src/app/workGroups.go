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

// CreateWorkGroup 创建工作组
func CreateWorkGroup(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	ownerID := r.PostFormValue("ownerID")
	ownerName := r.PostFormValue("ownerName")

	if name == "" {
		common.ReturnEFormat(w, 500, "工作组名称不能为空")
		return
	}

	if ownerID == "" {
		common.ReturnEFormat(w, 500, "负责人id不能为空")
		return
	}

	total, err := model.CountWorkGroups(fmt.Sprintf("WHERE name='%s'", name))
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	if total > 0 {
		common.ReturnEFormat(w, 500, "名称不能重复")
		return
	}

	workGroup := model.WorkGroup{}
	workGroup.Name = name
	workGroup.OwnerName = ownerName
	owner, err := strconv.ParseInt(ownerID, 10, 64)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	workGroup.OwnerID = owner
	err = workGroup.Insert()
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	workGroupNew, err := model.FindWorkGroup(fmt.Sprintf("WHERE name='%s'", name))
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}

	common.ReturnFormat(w, 200, map[string]interface{}{"workGroupNew": workGroupNew}, "SUCCESS")
}

func ListWorkGroups(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	q := r.FormValue("q")
	order := r.FormValue("order")
	limit := common.PageParams(r)
	conditions := []string{}
	condition := ""
	if q != "" {
		conditions = append(conditions, "name LIKE '%"+q+"%'")
	}
	if len(conditions) > 0 {
		condition = " WHERE " + strings.Join(conditions, " AND ")
	}
	if order != "" {
		order = " ORDER BY " + order
	}
	workGroups, err := model.FindWorkGroups(condition, limit, order)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	type workGroupJSON struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		OwnerID   int64  `json:"ownerId"`
		OwnerName string `json:"ownerName"`

		UpdatedAt time.Time        `json:"updatedAt"`
		CreatedAt time.Time        `json:"createdAt"`
		Employees []model.Employee `json:"employees"`
	}
	result := []workGroupJSON{}
	for _, workGroup := range workGroups {
		tmp := workGroupJSON{}
		tmp.ID = workGroup.ID
		tmp.Name = workGroup.Name
		tmp.OwnerID = workGroup.OwnerID
		tmp.OwnerName = workGroup.OwnerName
		tmp.CreatedAt = workGroup.CreatedAt
		tmp.UpdatedAt = workGroup.UpdatedAt
		employees, err := model.FindEmployees(fmt.Sprintf("WHERE work_group_id=%d", workGroup.ID), "", "")
		if err == nil {
			tmp.Employees = employees
		}
		result = append(result, tmp)
	}
	total, err := model.CountWorkGroups(condition)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	CreateOperateLogs("info", "查看", account, "查看工作组列表", r)
	common.ReturnFormat(w, 200, map[string]interface{}{"workGroups": result, "total": total}, "SUCCESS")
}

// UpdateWorkGroup 创建工作组
func UpdateWorkGroup(w http.ResponseWriter, r *http.Request) {
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	id := r.PostFormValue("id")
	name := r.PostFormValue("name")
	ownerID := r.PostFormValue("ownerID")
	ownerName := r.PostFormValue("ownerName")

	if name == "" {
		common.ReturnEFormat(w, 500, "工作组名称不能为空")
		return
	}

	if ownerID == "" {
		common.ReturnEFormat(w, 500, "负责人id不能为空")
		return
	}

	where := fmt.Sprintf("WHERE id=%s", id)
	workGroups, err := model.FindWorkGroups(where, "", "")
	if err != nil || len(workGroups) == 0 {
		common.ReturnEFormat(w, 500, "工作组不存在")
		return
	}
	workGroup := workGroups[0]
	update := []string{}
	if workGroup.Name != name {
		update = append(update, fmt.Sprintf("name='%s'", name))
	}
	if workGroup.OwnerName != ownerName {
		update = append(update, fmt.Sprintf("owner_name='%s'", ownerName))
	}

	owner, err := strconv.ParseInt(ownerID, 10, 64)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	if workGroup.OwnerID != owner {
		update = append(update, fmt.Sprintf("owner_id=%d", owner))
	}
	if len(update) > 0 {
		updater := strings.Join(update, ",")
		err := model.UpdateWorkGroups(updater, where)
		if err != nil {
			common.ReturnEFormat(w, 500, err.Error())
			return
		}
	}
	CreateOperateLogs("info", "修改", account, "修改工作组，修改结果：修改成功", r)
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func DeleteWorkGroups(w http.ResponseWriter, r *http.Request) {
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	workGroupIds := r.PostFormValue("workGroupIds")
	if workGroupIds == "" {
		common.ReturnEFormat(w, 500, "员工id不能为空")
		return
	}

	condition := fmt.Sprintf("WHERE id IN (%s)", workGroupIds)
	err = model.DeleteWorkGroups(condition)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	CreateOperateLogs("info", "删除", account, "删除工作组，删除结果：删除成功", r)
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func UpdateWorkGroupEmployees(w http.ResponseWriter, r *http.Request) {
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	id := r.PostFormValue("id")
	employeeIDStr := r.PostFormValue("employeeIds")

	workGroup, err := model.FindEmployee("WHERE id=" + id)
	if err != nil {
		common.ReturnEFormat(w, 500, "工作组不存在")
		return
	}

	employeeIDs := strings.Split(employeeIDStr, ",")
	employees := []model.Employee{}
	for _, employeeID := range employeeIDs {
		employee, err := model.FindEmployee("WHERE id=" + employeeID)
		if err == nil {
			if employee.WorkGroupID != 0 {
				common.ReturnEFormat(w, 500, fmt.Sprintf("员工%s已存在其他工作组中", employee.Name))
				return
			}
			employees = append(employees, employee)
		}
	}

	for _, employee := range employees {
		err = model.UpdateEmployees(fmt.Sprintf("work_group_id=%d", workGroup.ID), fmt.Sprintf("WHERE id=%d", employee.ID))
		if err != nil {
			logger.Error("workGroup.go", err.Error())
		}
	}
	CreateOperateLogs("info", "修稿", account, "修改工作组成员，修改结果：修改成功", r)
	common.ReturnFormat(w, 200, "", "SUCCESS")
}
