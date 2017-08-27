/**
 * 机构model实现
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
	OrganizeAuditStateWait   int = 0
	OrganizeAuditStatePass   int = 1
	OrganizeAuditStateRefuse int = 2
)

type Organization struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`            // 名称
	Code            string `json:"code"`            // 组织机构代码
	Contacts        string `json:"contacts"`        // 联系人
	ContactsPhone   string `json:"contactsPhone"`   // 联系电话
	Email           string `json:"email"`           // 邮箱
	BusinessLicense string `json:"businessLicense"` // 营业执照
	CodePic         string `json:"codePic"`         // 组织机构代码证
	AuditState      int    `json:"-"`               // 审核状态
	UserAccount     string `json:"userAccount"`     // 用户
	RefuseReason    string `json:"refuseReason"`

	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`

	AccountState int `json:"accountState"` //账号状态
}

func OrganizationTableName() string {
	return "t_organization"
}

func OrganizationViewName() string {
	return "v_organization"
}

func (m *Organization) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(name,code,contacts,contacts_account,email,business_license,code_pic,audit_state,user_account,updated_at,created_at) "+
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)", OrganizationTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	m.UpdatedAt = time.Now()
	m.CreatedAt = time.Now()

	_, err = stmt.Exec(m.Name, m.Code, m.Contacts, m.ContactsPhone, m.Email, m.BusinessLicense, m.CodePic,
		m.AuditState, m.UserAccount, m.UpdatedAt, m.CreatedAt)
	return
}

func FindOrganizations(condition, limit, order string) ([]Organization, error) {
	result := []Organization{}
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,name,code,contacts,contacts_account,email,business_license,code_pic,state,user_account,updated_at,created_at FROM %s %s %s %s", OrganizationViewName(), condition, order, limit))
	if err != nil {
		return result, err
	}

	for rows.Next() {
		tmp := Organization{}
		err = rows.Scan(&tmp.ID, &tmp.Name, &tmp.Code, &tmp.Contacts, &tmp.ContactsPhone, &tmp.Email,
			&tmp.BusinessLicense, &tmp.CodePic, &tmp.AccountState, &tmp.UserAccount, &tmp.UpdatedAt, &tmp.CreatedAt)
		result = append(result, tmp)
	}
	return result, err
}

func FindOrganization(condition string) (result Organization, err error) {
	result = Organization{}
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT id,name,code,contacts,contacts_account,email,business_license,code_pic,state,user_account,updated_at,created_at FROM %s %s", OrganizationViewName(), condition)).
		Scan(&result.ID, &result.Name, &result.Code, &result.Contacts, &result.ContactsPhone, &result.Email,
			&result.BusinessLicense, &result.CodePic, &result.AccountState, &result.UserAccount, &result.UpdatedAt, &result.CreatedAt)
	return
}

func UpdateOrganizations(update, condition string) (err error) {
	fmt.Println(fmt.Sprintf("UPDATE %s SET %s %s", OrganizationTableName(), update, condition))
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s SET %s %s", OrganizationTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = stmt.Exec()
	return
}

func CountOrganizations(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", OrganizationViewName(), condition)).Scan(&count)
	return
}
