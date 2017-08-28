/**
 * 工作组model实现
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

type WorkGroup struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	OwnerID   int64  `json:"ownerId"`
	OwnerName string `json:"ownerName"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

/*
CREATE TABLE "public"."t_work_group" (
	"id" SERIAL primary key,
	"name" varchar(500) COLLATE "default",
	"owner_id" int4,
	"owner_name" varchar(500) COLLATE "default",
	"created_at" date,
	"updated_at" date
)
WITH (OIDS=FALSE);
*/

func WorkGroupTableName() string {
	return "t_work_group"
}

func (m *WorkGroup) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(name,owner_id,owner_name,updated_at,created_at) "+
		"VALUES($1,$2,$3,$4,$5)", WorkGroupTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.UpdatedAt = time.Now()
	m.CreatedAt = time.Now()

	_, err = stmt.Exec(m.Name, m.OwnerID, m.OwnerName, m.UpdatedAt, m.CreatedAt)
	return
}

func FindWorkGroups(condition, limit, order string) ([]WorkGroup, error) {
	result := []WorkGroup{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,name,owner_id,owner_name,updated_at,created_at FROM %s %s %s %s", WorkGroupTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		tmp := WorkGroup{}
		err = rows.Scan(&tmp.ID, &tmp.Name, &tmp.OwnerID, &tmp.OwnerName, &tmp.UpdatedAt, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func FindWorkGroup(condition string) (result WorkGroup, err error) {
	result = WorkGroup{}
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT id,name,owner_id,owner_name,updated_at,created_at FROM %s %s", WorkGroupTableName(), condition)).
		Scan(&result.ID, &result.Name, &result.OwnerID, &result.OwnerName, &result.UpdatedAt, &result.CreatedAt)
	return
}

func UpdateWorkGroups(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s SET %s %s", WorkGroupTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}

func CountWorkGroups(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", WorkGroupTableName(), condition)).Scan(&count)
	return
}

func DeleteWorkGroups(condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("DELETE FROM %s %s", WorkGroupTableName(), condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}
