/**
 * 会员及管理员model实现
 *
 * @author zhangqichao
 * Created on 2017-08-18
 */
package model

import (
	"fmt"
	"../pdb"
	"time"
)

const (
	UserStateOn  int = 1
	UserStateOff int = 2
)

const (
	UserTypeSuper  int = 0
	UserTypeAdmin  int = 1
	UserTypeDoctor int = 2
	UserTypeNurse  int = 3
)

type User struct {
	Uid       int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	UserNo    string `json:"userNo"`
	UserLevel string `json:"userLevel"`
	Account   string `json:"account"`
	Pwd       string `json:"pwd"`
	State     int    `json:"state"`
	UserType  string `json:"userType"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

func UserTableName() string {
	return "t_user_info"
}

func (m *User) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(username,name,user_type,user_no,user_level,account,pwd,state,updated_at,created_at) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)", UserTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.UpdatedAt = time.Now()
	m.CreatedAt = time.Now()

	_, err = stmt.Exec(m.Username, m.Name, m.UserType, m.UserNo, m.UserLevel, m.Account, m.Pwd, m.State, m.UpdatedAt, m.CreatedAt)
	return
}

func FindUsers(condition, limit, order string) ([]User, error) {
	result := []User{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT uid,username,name,user_type,user_no,user_level,account,pwd,state,updated_at,created_at FROM %s %s %s %s", UserTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		tmp := User{}
		err = rows.Scan(&tmp.Uid, &tmp.Username, &tmp.Name, &tmp.UserType, &tmp.UserNo, &tmp.UserLevel, &tmp.Account, &tmp.Pwd, &tmp.State, &tmp.UpdatedAt, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func UpdateUsers(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s SET %s %s", UserTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}

func CountUsers(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", UserTableName(), condition)).Scan(&count)
	return
}
