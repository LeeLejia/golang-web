package app

import (
	"fmt"
	"../common"
	"../conf"
	"../model"
	"net/http"
	"os"
	"strings"
	"time"
)

func ListSystemLogs(w http.ResponseWriter, r *http.Request) {
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	conditions := []string{}
	order := r.FormValue("order")
	start := r.FormValue("start")
	end := r.FormValue("end")
	types := r.FormValue("types")
	if types != "" {
		conditions = append(conditions, fmt.Sprintf("type='%s'", types))
	}
	if start != "" && end != "" {
		conditions = append(conditions, fmt.Sprintf("created_at>=timestamp '%s' AND created_at<=timestamp '%s'", start, end))
	}
	limit := common.PageParams(r)
	sort := ""
	if order != "" && strings.Contains(order, "DESC") {
		sort = " ORDER BY " + order
	}
	condition := ""
	if len(conditions)>0 {
		condition = " WHERE " + strings.Join(conditions, " AND ")
	}
	fmt.Println(condition)
	systemLogs, err := model.FindSystemLogs(condition, limit, sort)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	total, err := model.CountSystemLogs(condition)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	CreateOperateLogs("info", "查看", account, "查看系统日志列表", r)
	common.ReturnFormat(w, 200, map[string]interface{}{"systemLogs": systemLogs, "total": total}, "SUCCESS")
}

func ExportSystemLogs(w http.ResponseWriter, r *http.Request) {
	account, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	conditions := []string{}
	order := r.FormValue("order")
	start := r.FormValue("start")
	end := r.FormValue("end")
	types := r.FormValue("types")
	if types != "" {
		conditions = append(conditions, fmt.Sprintf("type='%s'", types))
	}
	if start != "" && end != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= timestamp '%s' AND created_at <= timestamp '%s'", start, end))
	}
	limit := common.PageParams(r)
	sort := ""
	if order != "" && strings.Contains(order, "DESC") {
		sort = " ORDER BY " + order
	}
	condition := ""
	if len(conditions)>0 {
		condition = " WHERE " + strings.Join(conditions, " AND ")
	}
	systemLogs, err := model.FindSystemLogs(condition, limit, sort)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}

	filePath := fmt.Sprintf(conf.App.FileUploadPath+"systemLogs_%d.txt", time.Now().Unix())

	f, err := os.OpenFile("."+filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	} else {
		for _, systemLog := range systemLogs {
			content := fmt.Sprintf("time:%s type:%s model:%s content:%s\r\n", systemLog.CreatedAt.Format("2006-01-02 15:04:05 -0700"), systemLog.Type, systemLog.Part, systemLog.Content)
			f.Write([]byte(content))
		}
	}
	f.Close()
	CreateOperateLogs("info", "查看", account, "导出系统日志列表", r)
	common.ReturnFormat(w, 200, conf.App.ServerHost+filePath, "SUCCESS")
}
