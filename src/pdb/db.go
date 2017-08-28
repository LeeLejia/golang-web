/**
 * 公共数据库连接配置
 */
package pdb

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

var Session *sql.DB

func InitDB(host, port, user, pwd, dbName string) error {
	dateSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, dbName)
	fmt.Println(dateSource)
	db, _ := sql.Open("postgres", dateSource)
	Session = db
	err := Session.Ping()
	if err != nil {
		go reInit(dateSource, 1)
	}
	return nil
}

func reInit(dateSource string, seconds int) {
	for {
		db, _ := sql.Open("postgres", dateSource)
		if err := db.Ping(); err == nil {
			Session = db
			break
		} else {
			fmt.Println("数据库连接失败，2分钟后重试")
			time.Sleep(time.Minute * 2)
			reInit(dateSource, seconds*2)
		}
	}
}
