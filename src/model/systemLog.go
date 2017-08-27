/**
 * 系统日志model实现
 *
 * @author zhangqichao
 * Created on 2017-08-18
 */
package model

import (
	"fmt"
	"time"
	"../pdb"
)

type SystemLog struct {
	ID      int64  `json:"id"`
	Type    string `json:"type"`
	Part    string `json:"part"`
	Content string `json:"content"`

	CreatedAt time.Time `json:"createdAt"`
}

func SystemLogTableName() string {
	return "t_system_log"
}

func (m *SystemLog) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(type,part,content,created_at) "+
		"VALUES($1,$2,$3,$4)", SystemLogTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.CreatedAt = time.Now()

	_, err = stmt.Exec(m.Type, m.Part, m.Content, m.CreatedAt)
	return
}

func FindSystemLogs(condition, limit, order string) ([]SystemLog, error) {
	result := []SystemLog{}
	fmt.Println(fmt.Sprintf("SELECT id,type,part,content,created_at FROM %s %s %s %s", SystemLogTableName(), condition, order, limit))
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,type,part,content,created_at FROM %s %s %s %s", SystemLogTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		tmp := SystemLog{}
		err = rows.Scan(&tmp.ID, &tmp.Type, &tmp.Part, &tmp.Content, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func UpdateSystemLogs(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s SET %s %s", SystemLogTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}

func CountSystemLogs(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", SystemLogTableName(), condition)).Scan(&count)
	return
}

func DeleteSystemLogs(condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("DELETE FROM %s %s", SystemLogTableName(), condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}
