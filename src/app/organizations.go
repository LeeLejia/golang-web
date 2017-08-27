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

func ListOrganizations(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	q := r.FormValue("q")
	order := r.FormValue("order")
	auditState := r.FormValue("auditState")
	accountState := r.FormValue("accountState")
	limit := common.PageParams(r)

	conditions := []string{}
	condition := ""

	auditStateInt, err := strconv.Atoi(auditState)
	if err == nil {
		conditions = append(conditions, fmt.Sprintf("audit_state=%d", auditStateInt))
	}
	if q != "" {
		conditions = append(conditions, "(name LIKE '%"+q+"%' OR code LIKE '%"+q+"%')")
	}
	if accountState != "" {
		conditions = append(conditions, fmt.Sprintf("state=%s", accountState))
	}
	if len(conditions) > 0 {
		condition = " WHERE " + strings.Join(conditions, " AND ")
	}
	if order != "" {
		order = " ORDER BY " + order
	}
	fmt.Println(condition)
	fmt.Println(limit)
	organizations, err := model.FindOrganizations(condition, limit, order)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	total, err := model.CountOrganizations(condition)
	if err != nil {
		common.ReturnEFormat(w, 503, err.Error())
		return
	}
	CreateOperateLogs("info", "查看", "不需要登录的接口", "查看机构列表", r)
	common.ReturnFormat(w, 200, map[string]interface{}{"organizations": organizations, "total": total}, "SUCCESS")
}

func CreateOrganization(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)

	name := r.PostFormValue("name")
	code := r.PostFormValue("code")
	contacts := r.PostFormValue("contacts")
	contactsPhone := r.PostFormValue("contactsPhone")
	email := r.PostFormValue("email")
	businessLicense := r.PostFormValue("businessLicense")
	codePic := r.PostFormValue("codePic")
	if name == "" {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因：名称不能为空", r)
		common.ReturnEFormat(w, 500, "名称不能为空")
		return
	}
	if code == "" {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因：组织机构代码不能为空", r)
		common.ReturnEFormat(w, 500, "组织机构代码不能为空")
		return
	}
	if contacts == "" {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因：联系人不能为空", r)
		common.ReturnEFormat(w, 500, "联系人不能为空")
		return
	}
	if contactsPhone == "" {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因：联系人电话不能为空", r)
		common.ReturnEFormat(w, 500, "联系人电话不能为空")
		return
	}
	if email == "" {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因：邮箱不能为空", r)
		common.ReturnEFormat(w, 500, "邮箱不能为空")
		return
	}
	if businessLicense == "" {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因：请上传营业执照", r)
		common.ReturnEFormat(w, 500, "请上传营业执照")
		return
	}
	if codePic == "" {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因：请上传组织机构代码证", r)
		common.ReturnEFormat(w, 500, "请上传组织机构代码证")
		return
	}
	where := fmt.Sprintf("WHERE code='%s'", code)
	total, err := model.CountOrganizations(where)
	if err != nil {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因："+err.Error(), r)
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	if total > 0 {
		common.ReturnEFormat(w, 500, "机构已存在")
		return
	}

	organization := model.Organization{}
	organization.Name = name
	organization.Code = code
	organization.Contacts = contacts
	organization.ContactsPhone = contactsPhone
	organization.Email = email
	organization.BusinessLicense = businessLicense
	organization.CodePic = codePic
	user, err := generateUser(organization.Code, "123456", organization.Name, "admin")
	if err != nil {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因："+err.Error(), r)
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	organization.UserAccount = user.Account
	err = organization.Insert()
	if err != nil {
		CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：失败，失败原因："+err.Error(), r)
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	CreateOperateLogs("error", "创建", "不需要登录的接口", "创建机构，操作结果：成功"+err.Error(), r)
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	account, err := common.CheckSession(r)
	if err != nil {
		CreateOperateLogs("error", "修改", "机构管理员", "修改机构信息，操作结果：失败，失败原因："+err.Error(), r)
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	id := r.PostFormValue("id")
	name := r.PostFormValue("name")
	code := r.PostFormValue("code")
	contacts := r.PostFormValue("contacts")
	contactsPhone := r.PostFormValue("contactsPhone")
	email := r.PostFormValue("email")
	if name == "" {
		CreateOperateLogs("error", "修改", account, "修改机构信息，操作结果：失败，失败原因：名称不能为空", r)
		common.ReturnEFormat(w, 500, "名称不能为空")
		return
	}
	if code == "" {
		CreateOperateLogs("error", "修改", account, "修改机构信息，操作结果：失败，失败原因：组织机构代码不能为空", r)
		common.ReturnEFormat(w, 500, "组织机构代码不能为空")
		return
	}
	if contacts == "" {
		CreateOperateLogs("error", "修改", account, "修改机构信息，操作结果：失败，失败原因：联系人不能为空", r)
		common.ReturnEFormat(w, 500, "联系人不能为空")
		return
	}
	if contactsPhone == "" {
		CreateOperateLogs("error", "修改", account, "修改机构信息，操作结果：失败，失败原因：联系人电话不能为空", r)
		common.ReturnEFormat(w, 500, "联系人电话不能为空")
		return
	}
	if email == "" {
		CreateOperateLogs("error", "修改", account, "修改机构信息，操作结果：失败，失败原因：邮箱不能为空", r)
		common.ReturnEFormat(w, 500, "邮箱不能为空")
		return
	}
	if !common.IsEmail(email) {
		CreateOperateLogs("error", "修改", account, "修改机构信息，操作结果：失败，失败原因：邮箱格式错误", r)
		common.ReturnEFormat(w, 500, "邮箱格式错误")
		return
	}
	where := fmt.Sprintf("WHERE code='%s'", code)
	total, err := model.CountOrganizations(where)
	if err != nil {
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	if total == 0 {
		CreateOperateLogs("error", "修改", account, "修改机构信息，操作结果：失败，失败原因：机构不存在", r)
		common.ReturnEFormat(w, 500, "机构不存在")
		return
	}

	condition := fmt.Sprintf("WHERE id=%s", id)
	organizations, err := model.FindOrganizations(condition, "", "")
	if err == nil && len(organizations) > 0 {
		organization := organizations[0]
		update := []string{}
		if organization.Name != name {
			update = append(update, fmt.Sprintf("name='%s'", name))
		}
		if organization.Contacts != contacts {
			update = append(update, fmt.Sprintf("contacts='%s'", contacts))
		}
		if organization.ContactsPhone != contactsPhone {
			update = append(update, fmt.Sprintf("contacts_account='%s'", contactsPhone))
		}
		if organization.Email != email {
			update = append(update, fmt.Sprintf("email='%s'", email))
		}
		if len(update) > 0 {
			updater := strings.Join(update, ",")
			err := model.UpdateOrganizations(updater, condition)
			if err != nil {
				CreateOperateLogs("error", "修改", account, "修改机构信息，操作结果：失败，失败原因："+err.Error(), r)
				common.ReturnEFormat(w, 500, err.Error())
				return
			}
			CreateOperateLogs("error", "修改", account, "修改机构信息，操作结果：成功", r)
		}
	}
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func AuditOrganization(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
	account, err := common.CheckSession(r)
	if err != nil {
		CreateOperateLogs("error", "修改", account, "审核机构信息，操作结果：失败，失败原因："+err.Error(), r)
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	id := r.PostFormValue("id")
	pass := r.PostFormValue("pass")
	reason := r.PostFormValue("reason")
	if id == "" {
		CreateOperateLogs("error", "修改", account, "审核机构信息，操作结果：失败，失败原因：id不能为空", r)
		common.ReturnEFormat(w, 500, "id不能为空")
		return
	}
	update := []string{}
	where := fmt.Sprintf("WHERE id=%s", id)
	organization, err := model.FindOrganization(where)
	if err != nil {
		CreateOperateLogs("error", "修改", account, "审核机构信息，操作结果：失败，失败原因："+err.Error(), r)
		common.ReturnEFormat(w, 500, "获取机构信息失败")
		return
	}
	if pass == "true" {
		update = append(update, fmt.Sprintf("audit_state=%d", model.OrganizeAuditStatePass))
		SendEmail(organization.Contacts, organization.Email, organization.Code, true)
	} else if pass == "false" {
		update = append(update, fmt.Sprintf("audit_state=%d", model.OrganizeAuditStateRefuse))
		update = append(update, fmt.Sprintf("refuse_reason='%s'", reason))
		SendEmail(organization.Contacts, organization.Email, organization.Code, false)
	} else {
		CreateOperateLogs("error", "修改", account, "审核机构信息，操作结果：失败，失败原因：审核状态错误", r)
		common.ReturnEFormat(w, 500, "审核状态错误")
		return
	}
	updater := strings.Join(update, ",")
	err = model.UpdateOrganizations(updater, where)
	if err != nil {
		CreateOperateLogs("error", "修改", account, "审核机构信息，操作结果：失败，失败原因："+err.Error(), r)
		common.ReturnEFormat(w, 500, err.Error())
		return
	}
	common.ReturnFormat(w, 200, nil, "SUCCESS")
}

func SendEmail(name, email, code string, state bool) {
	user := "hm-service@hm.rn830.com"
	password := "GDCChm12345"
	host := "smtpdm.aliyun.com:25"
	to := email
	subject := "健康管理系统"
	body := ``
	fmt.Println("send email")
	if state {
		body = `
        <html>
        <body>
        <h3>
        <div>尊敬的` + name + `先生：</div>
        	<div>您好！恭喜您，贵机构申请入驻健康管理平台已通过审核，请尽快登录添加并管理机构员工。</div>
       		<div>管理员账号：` + code + `</div>
       		<div>管理员密码：123456</div>
       		<div>如需修改密码，请联系平台管理员</div>
        </h3>
        </body>
        </html>
        `
	} else {
		body = `
        <html>
        <body>
        <h3>
        <div>尊敬的` + name + `先生：</div>
        	<div>您好！很抱歉，由于您提供的重要信息缺失或有误，贵机构申请入驻健康管理平台未通过审核，请核对信息重新申请入驻本平台。</div>
       		<div>申请入驻访问，或直接联系我们</div>
        </h3>
        </body>
        </html>
        `
	}
	err := common.SendToMail(user, password, host, to, subject, body, "html")
	if err != nil {
		logger.Error(fmt.Sprintf("Send mail error! %s", err.Error()), "organizations.go")
	} else {
		logger.Info("Send mail success!", "organizations.go")
	}
}

func CheckAudit(w http.ResponseWriter, r *http.Request) {
	_, err := common.CheckSession(r)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	now := time.Now()
	past := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()-30, 0, now.Location())
	where := fmt.Sprintf(" WHERE created_at >= timestamp '%s' AND  created_at <= timestamp '%s' ", past.Local().Format("2006-01-02 15:04:05"), now.Local().Format("2006-01-02 15:04:05"))
	count, err := model.CountOrganizations(where)
	if err != nil {
		common.ReturnEFormat(w, 403, err.Error())
		return
	}
	if count > 0 {
		common.ReturnFormat(w, 200, true, "SUCCESS")
	} else {
		common.ReturnFormat(w, 200, false, "SUCCESS")
	}
}
