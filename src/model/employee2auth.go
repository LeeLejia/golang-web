/**
 * 员工和权限绑定model实现
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

type Employee2auth struct {
	ID         int64 `json:"id"`
	Auth       int64 `json:"auth"`
	EmployeeID int64 `json:"employeeId"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

func Employee2authTableName() string {
	return "t_employee_2_auth"
}

func (m *Employee2auth) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(auth,employee_id,updated_at,created_at) "+
		"VALUES($1,$2,$3,$4)", Employee2authTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.UpdatedAt = time.Now()
	m.CreatedAt = time.Now()

	_, err = stmt.Exec(m.Auth, m.EmployeeID, m.UpdatedAt, m.CreatedAt)
	return
}

func FindEmployee2auths(condition, limit, order string) ([]Employee2auth, error) {
	result := []Employee2auth{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,auth,employee_id,updated_at,created_at FROM %s %s %s %s", Employee2authTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		tmp := Employee2auth{}
		err = rows.Scan(&tmp.ID, &tmp.Auth, &tmp.EmployeeID, &tmp.UpdatedAt, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func UpdateEmployee2auths(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s SET %s %s", Employee2authTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}

func CountEmployee2auths(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", Employee2authTableName(), condition)).Scan(&count)
	return
}

func DeleteEmployee2auths(condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("DELETE FROM %s %s", Employee2authTableName(), condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}
