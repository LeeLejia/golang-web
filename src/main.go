package main

import (
	"./app"
	"./common/conf"
	"fmt"
	"./pdb"
	"net/http"
)

func main() {
	conf.Init("./app.toml")
	pdb.InitDB(conf.App.DBHost, conf.App.DBPort, conf.App.DBUser, conf.App.DBPassword, conf.App.DBName)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/register", app.Register)
	http.HandleFunc("/logout", app.Logout)
	http.HandleFunc("/reset_password", app.ResetPassword)

	http.HandleFunc("/users/update", app.UpdateUser)
	http.HandleFunc("/users/read", app.ListUsers)

	http.HandleFunc("/test", app.Test)

	fsh := http.FileServer(http.Dir(conf.App.Static))
	http.Handle("/static/", http.StripPrefix("/static/", fsh))

	http.ListenAndServe(fmt.Sprintf(":%s", conf.App.ServerPort), nil)
}
