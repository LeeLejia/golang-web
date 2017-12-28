package model

import (
	"time"
	"fmt"
	"../pdb"
)
/*
DROP TABLE IF EXISTS t_log;
CREATE TABLE t_log (
"id" serial NOT NULL,
"type" varchar(16) COLLATE "default",
"tag" varchar(256) COLLATE "default",
"operator" varchar(128),
"content" varchar(512) COLLATE "default",
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */

const(
	LOG_TYPE_DEBUG="debug"
	LOG_TYPE_INFO="info"
	LOG_TYPE_WARM="warn"
	LOG_TYPE_ERROR="error"
	LOG_TYPE_NORMAL="normal"
r
)

type T_log struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Tag       string    `json:"tag"`
	Operator  string    `json:"operator"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func LogTableName() string {
	return "t_log"
}

func InsertAllLog(ls []T_log)(err error){
	tx,err:=pdb.Session.Begin()
	if err != nil {
		fmt.Println("Beginx error:", err)
		return
	}
	sql:=fmt.Sprintf("INSERT INTO %s(type,tag,operator,content,created_at) "+
		"VALUES($1,$2,$3,$4,$5)", LogTableName())
	for _,m:=range ls{
		_,err:=tx.Exec(sql,m.Type, m.Tag, m.Operator, m.Content, m.CreatedAt)
		if err!=nil{
			tx.Rollback()
		}
	}
	tx.Commit()
	return
}

func (m *T_log) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(type,tag,operator,content,created_at) "+
			  "VALUES($1,$2,$3,$4,$5)", LogTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()
	m.CreatedAt = time.Now()
	_, err = stmt.Exec(m.Type, m.Tag, m.Operator, m.Content, m.CreatedAt)
	return
}

func FindLogs(condition, limit, order string) ([]T_log, error) {
	result := []T_log{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,type,tag,operator,content,created_at FROM %s %s %s %s", LogTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}
	for rows.Next() {
		tmp := T_log{}
		err = rows.Scan(&tmp.ID, &tmp.Type,&tmp.Tag,&tmp.Operator, &tmp.Content, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func CountLogs(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", LogTableName(), condition)).Scan(&count)
	return
}
