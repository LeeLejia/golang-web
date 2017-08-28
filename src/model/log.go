package model

import (
	"time"
	"fmt"
	"../pdb"
)

type Log struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Tag       string `json:"tag"`
	User      string `json:"user"`
	Content   string `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
func LogTableName() string {
	return "t_log"
}
/**
	插入一条日志
 */
func (m *Log) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(type,tag,user,content,created_at) "+
			  "VALUES($1,$2,$3,$4,$5)", LogTableName()))
	if err != nil {
		fmt.Println(fmt.Sprintf("\x1b[%dm插入日志失败！detail:%s\x1b[0m",91,err.Error()))
		return
	}
	m.CreatedAt = time.Now()
	_, err = stmt.Exec(m.Type, m.Tag, m.User, m.Content, m.CreatedAt)
	return
}
/**
	查找日志
 */
func ListLogs(condition, limit, order string) ([]Log, error) {
	result := []Log{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,type,tag,user,created_at FROM %s %s %s %s", LogTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}
	for rows.Next() {
		tmp := Log{}
		err = rows.Scan(&tmp.ID, &tmp.Type,&tmp.Tag,&tmp.User, &tmp.Content, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}
/**
	获取日志数量
 */
func GetCount(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", LogTableName(), condition)).Scan(&count)
	return
}
