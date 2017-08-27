package model

import (
	"fmt"
	"time"

	"../pdb"
)

type Organization2Employee struct {
	ID             int64 `json:"id"`
	OrganizationID int64 `json:"organizationId"`
	EmployeeID     int64 `json:"employeeId"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

/*
CREATE TABLE "public"."t_organization_2_employee" (
	"id" SERIAL primary key,
	"organization_id" int4,
	"employee_id" int4,
	"created_at" date,
	"updated_at" date
)
WITH (OIDS=FALSE);
*/

func Organization2EmployeeTableName() string {
	return "t_organization_2_employee"
}

func (m *Organization2Employee) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(organization_id,employee_id,updated_at,created_at) "+
		"VALUES($1,$2,$3,$4)", Organization2EmployeeTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.UpdatedAt = time.Now()
	m.CreatedAt = time.Now()

	_, err = stmt.Exec(m.OrganizationID, m.EmployeeID, m.UpdatedAt, m.CreatedAt)
	return
}

func FindOrganization2Employees(condition, limit, order string) ([]Organization2Employee, error) {
	result := []Organization2Employee{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,organization_id,employee_id,updated_at,created_at FROM %s %s %s %s", Organization2EmployeeTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		tmp := Organization2Employee{}
		err = rows.Scan(&tmp.ID, &tmp.OrganizationID, &tmp.EmployeeID, &tmp.UpdatedAt, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func UpdateOrganization2Employees(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s SET %s %s", Organization2EmployeeTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}

func CountOrganization2Employees(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", Organization2EmployeeTableName(), condition)).Scan(&count)
	return
}
