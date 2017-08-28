package app

import (
	"fmt"
	"../src/common"
	"../src/common/logger"
	"../src/common/conf"
	"../src/model"
	"net/http"
	"os"
	"strings"
	"time"
)

func CreateOperateLogs(types, operateType, operator, content string, r *http.Request) {
	ip := r.Header.Get("Remote_addr")
	if ip == "" {
		ip = r.RemoteAddr
	}
	operateLog := model.OperateLog{}
	operateLog.IP = ip
	operateLog.Type = types
	operateLog.OperateType = operateType
	operateLog.Operator = operator
	operateLog.Content = content
	err := operateLog.Insert()
	if err != nil {
		log.E("operateLogs.go",err.Error())
	}
}

func ListOperateLogs(w http.ResponseWriter, r *http.Request) {
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
	operateLogs, err := model.FindOperateLogs(condition, limit, sort)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	total, err := model.CountOperateLogs(condition)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	CreateOperateLogs("info", "查看", account, "查看操作日志列表", r)
	common.ReturnFormat(w, 200, map[string]interface{}{"operateLogs": operateLogs, "total": total}, "SUCCESS")
}

func ExportOperateLogs(w http.ResponseWriter, r *http.Request) {
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
	operateLogs, err := model.FindOperateLogs(condition, limit, sort)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}

	filePath := fmt.Sprintf(conf.App.FileUploadPath+"operateLogs_%d.txt", time.Now().Unix())

	f, err := os.OpenFile("."+filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	} else {
		for _, operateLog := range operateLogs {
			content := fmt.Sprintf("time:%s type:%s operateType:%s operator:%s IP:%s content:%s\r\n",
				operateLog.CreatedAt.Format("2006-01-02 15:04:05 -0700"), operateLog.Type, operateLog.OperateType,
				operateLog.Operator, operateLog.IP, operateLog.Content)
			f.Write([]byte(content))
		}
	}
	f.Close()
	CreateOperateLogs("info", "查看", account, "导出操作日志列表", r)
	common.ReturnFormat(w, 200, conf.App.ServerHost+filePath, "SUCCESS")
}
