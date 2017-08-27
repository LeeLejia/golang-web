package main

import (
	"../back-end-api/app"
	"../back-end-api/conf"
	"fmt"
	"../back-end-api/common/logger"
	"../back-end-api/pdb"
	"net/http"
)

func main() {
	conf.Init("./app.toml")
	pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	logger.InitLogger()
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/register", app.Register)
	http.HandleFunc("/logout", app.Logout)
	http.HandleFunc("/reset_password", app.ResetPassword)

	http.HandleFunc("/users/update", app.UpdateUser)
	http.HandleFunc("/users/read", app.ListUsers)

	http.HandleFunc("/upload_image", app.UploadPicture)

	http.HandleFunc("/organizations/read", app.ListOrganizations)
	http.HandleFunc("/organizations/create", app.CreateOrganization)
	http.HandleFunc("/organizations/update", app.UpdateOrganization)
	http.HandleFunc("/organizations/audit", app.AuditOrganization)
	http.HandleFunc("/organizations/check", app.CheckAudit)

	http.HandleFunc("/work_groups/create", app.CreateWorkGroup)
	http.HandleFunc("/work_groups/read", app.ListWorkGroups)
	http.HandleFunc("/work_groups/update", app.UpdateWorkGroup)
	http.HandleFunc("/work_groups/delete", app.DeleteWorkGroups)
	http.HandleFunc("/work_groups_employees/update", app.UpdateWorkGroupEmployees)

	http.HandleFunc("/employees/create", app.CreateEmployee)
	http.HandleFunc("/employees/read", app.ListEmployees)
	http.HandleFunc("/employees/update", app.UpdateEmployee)
	http.HandleFunc("/employees/delete", app.DeleteEmployees)

	http.HandleFunc("/employees/auth/read", app.ListEmployee2Auth)
	http.HandleFunc("/employees/auth/no_auth", app.ListEmployeeWithoutAuth)
	http.HandleFunc("/employees/auth/update", app.AddEmployee2Auth)

	http.HandleFunc("/system_logs/read", app.ListSystemLogs)
	http.HandleFunc("/system_logs/export", app.ExportSystemLogs)

	http.HandleFunc("/operate_logs/read", app.ListOperateLogs)
	http.HandleFunc("/operate_logs/export", app.ExportOperateLogs)

	http.HandleFunc("/test", app.Test)

	fsh := http.FileServer(http.Dir(conf.App.Static))
	http.Handle("/static/", http.StripPrefix("/static/", fsh))

	http.ListenAndServe(fmt.Sprintf(":%s", conf.App.ServerPort), nil)
}
