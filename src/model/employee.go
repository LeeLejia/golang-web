/**
 * 员工model实现
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

const (
	EmployeeAuthSuper  int = 1
	EmployeeAuthAdmin  int = 2
	EmployeeAuthDoctor int = 3
	EmployeeAuthNurse  int = 4
)

type Employee struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`    //姓名
	Phone        string `json:"account"` //电话
	Sex          string `json:"sex"`     //性别
	Title        string `json:"title"`   //职称
	WorkNo       string `json:"workNo"`  //工号
	ManageNum    int64  `json:"manageNum"`
	EmployeeType int64  `json:"employeeType"`
	State        int    `json:"state"`
	Pwd          string `json:"pwd"`
	OID          int64  `json:"oId"`
	WorkGroupID  int64  `json:"workGroupId"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	Auth int `json:"auth"`
}

func EmployeeTableName() string {
	return "t_employee"
}

func EmployeeViewName() string {
	return "v_employee"
}

func (m *Employee) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(name,phone,sex,title,work_no,manage_num,employee_type,pwd,state,o_id,updated_at,created_at) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", EmployeeTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.UpdatedAt = time.Now()
	m.CreatedAt = time.Now()

	_, err = stmt.Exec(m.Name, m.Phone, m.Sex, m.Title, m.WorkNo, m.ManageNum, m.EmployeeType, m.Pwd, m.State, m.OID, m.UpdatedAt, m.CreatedAt)
	return
}

func FindEmployees(condition, limit, order string) ([]Employee, error) {
	result := []Employee{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,name,phone,sex,title,work_no,manage_num,employee_type,pwd,state,o_id,work_group_id,updated_at,created_at FROM %s %s %s %s", EmployeeViewName(), condition, order, limit))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		tmp := Employee{}
		err = rows.Scan(&tmp.ID, &tmp.Name, &tmp.Phone, &tmp.Sex, &tmp.Title, &tmp.WorkNo, &tmp.ManageNum, &tmp.EmployeeType, &tmp.Pwd, &tmp.State, &tmp.OID, &tmp.WorkGroupID, &tmp.UpdatedAt, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func FindEmployee(condition string) (result Employee, err error) {
	result = Employee{}
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT id,name,phone,pwd,sex,title,work_no,manage_num,employee_type,state,o_id,work_group_id,updated_at,created_at FROM %s %s", EmployeeTableName(), condition)).
		Scan(&result.ID, &result.Name, &result.Phone, &result.Pwd, &result.Sex, &result.Title, &result.WorkNo, &result.ManageNum, &result.EmployeeType, &result.State, &result.UpdatedAt, &result.CreatedAt)
	return
}

func UpdateEmployees(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s SET %s %s", EmployeeTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}

func CountEmployees(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", EmployeeTableName(), condition)).Scan(&count)
	return
}

func DeleteEmployees(condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("DELETE FROM %s %s", EmployeeTableName(), condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}
