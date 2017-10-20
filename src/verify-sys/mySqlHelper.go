//mysqlHelper
package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var rootDbPwd = "imjia123"
var dB = "verifyNew"

//初始化数据库
func setupDB() {
	var err error
	connStr := "root:" + rootDbPwd + "@/mysql?charset=utf8&loc=Local&parseTime=true"
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		log(err.Error(), MSG_ERR)
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log(err.Error(), MSG_ERR)
		panic(err.Error())
	}

	cr_db := "CREATE DATABASE IF NOT EXISTS " + dB + " DEFAULT CHARSET utf8 COLLATE utf8_general_ci;"
	exeSQL(cr_db)
	grantSQL := "grant all on " + dB + ".* to root identified by '" + rootDbPwd + "';"
	exeSQL(grantSQL)
	grantSQL = "grant all on " + dB + ".* to root@'' identified by '" + rootDbPwd + "';"
	exeSQL(grantSQL)
	grantSQL = "grant all on " + dB + ".* to root@'localhost' identified by '" + rootDbPwd + "';"
	exeSQL(grantSQL)
	grantSQL = "grant all on " + dB + ".* to root@'127.0.0.1' identified by '" + rootDbPwd + "';"
	exeSQL(grantSQL)

	db.Close()

	connStr = "root:" + rootDbPwd + "@/" + dB + "?charset=utf8&loc=Local&parseTime=true"
	db, err = sql.Open("mysql", connStr)
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	//初始化表
	//valid -1是不可用,1是可用对特定机器,2是可用对全部机器,3运行将验证结果保存到机器
	exeSQL("create table if not exists codes(code char(10) not null primary key,exp varchar(800) not null default '',mostAskTimes int not null DEFAULT -1, askTimes int not null DEFAULT 0,mMachineCount int not null DEFAULT -1,machine varchar(1000) not null default '',validType int not null,createtime TIMESTAMP DEFAULT current_timestamp() not null);")
	//黑名单
	exeSQL("create table if not exists blackList(exp varchar(800),machine varchar(100),ip char(20),createtime TIMESTAMP DEFAULT current_timestamp() not null);")
	//白名单
	exeSQL("create table if not exists whiteList(exp varchar(800),machine varchar(100),ip char(20),createtime TIMESTAMP DEFAULT current_timestamp() not null);")
	//日志
	exeSQL("create table if not exists logs(protosign int not null,msgType int not null,code char(10) not null,exp varchar(800),machine varchar(1000),ip char(20) not null,createtime TIMESTAMP DEFAULT current_timestamp() not null);")
}

//执行sql,不获取返回值
func exeSQL(command string, args ...interface{}) bool {
	lock.Lock()
	log("sql:"+command, MSG_DEBUG)
	stmt, err := db.Prepare(command)
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
		lock.Unlock()
	}()
	if err != nil {
		log(err.Error(), MSG_ERR)
		//panic(err.Error())
		return false
	}
	if args != nil {
		_, err = stmt.Exec(args...)
	} else {
		_, err = stmt.Exec()
	}

	if err != nil {
		log(err.Error(), MSG_ERR)
		return false
	}
	return true
}

//获取数据，返回数据列表
func exeSQLforResult(command string, args ...interface{}) (*sql.Rows, bool) {
	lock.Lock()
	log("sql:"+command, MSG_DEBUG)
	ok := true
	stmt, err := db.Prepare(command)
	if err != nil {
		log(err.Error(), MSG_ERR)
		//panic(err.Error())
		ok = false
	}
	if err != nil {
		log(err.Error(), MSG_ERR)
		//panic(err.Error())
		ok = false
	}
	var rows *sql.Rows
	if args != nil {
		rows, err = stmt.Query(args...)
	} else {
		rows, err = stmt.Query()
	}

	if err != nil {
		log(err.Error(), MSG_ERR)
		//panic(err.Error())
		ok = false
	}
	defer func() {
		if stmt != nil {
			stmt.Close()
		}
		lock.Unlock()
	}()
	return rows, ok
}
