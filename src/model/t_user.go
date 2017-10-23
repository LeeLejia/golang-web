package model

import (
	"time"
	"github.com/bitly/go-simplejson"
	"fmt"
	"../pdb"
	"math/rand"
)

/*
DROP TABLE IF EXISTS t_user;
CREATE TABLE t_user (
"id" serial NOT NULL,
"role" varchar(16) COLLATE "default",
"nick" varchar(128) COLLATE "default",
"pwd" varchar(128) COLLATE "default",
"avatar" varchar(128) COLLATE "default",
"phone" varchar(11) COLLATE "default",
"email" varchar(128) COLLATE "default",
"qq" varchar(13) COLLATE "default",
"status" int4 DEFAULT 1,
"expend" jsonb NOT NULL,
"updated_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */

const (
	USER_ROLE_COMMON = "common"
	USER_ROLE_DEVELOPER = "developer"
	USER_ROLE_ADMIN = "admin"
	USER_ROLE_SUPER = "super"

	USER_STATUS_INVALID=0
	USER_STATUS_VALID=1

)
// 用戶表
type T_user struct {
	Id        int                `json:"id"`
	Role      string             `json:"role"`
	Nick      string             `json:"nick"`
	Pwd       string             `json:"pwd"`
	Avatar    string             `json:"avatar"`
	Phone     string             `json:"phone"`
	Email     string             `json:"email"`
	QQ        string             `json:"qq"`
	Status    int             	 `json:"status"`
	Expend    * simplejson.Json  `json:"expend"`
	UpdatedAt time.Time          `json:"updated_at"`
	CreatedAt time.Time          `json:"created_at"`
}

func UserTableName() string{
	return "t_user"
}

func (u* T_user) Insert() (err error){
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(role,nick,pwd,avatar,phone,email,qq,status,expend,created_at,updated_at) "+
			  "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)", UserTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	u.UpdatedAt = time.Now()
	u.CreatedAt = time.Now()
	if u.Expend==nil{
		u.Expend=&simplejson.Json{}
	}
	if u.Role==USER_ROLE_DEVELOPER{
		u.Expend.Set("devId",getRandomString(8))
		u.Expend.Set("devKey",getRandomString(8))
	}
	d,err:=u.Expend.MarshalJSON()
	if err!=nil{
		return err
	}
	_, err = stmt.Exec(u.Role,u.Nick,u.Pwd,u.Avatar,u.Phone,u.Email,u.QQ,u.Status,string(d),u.CreatedAt,u.UpdatedAt)
	return
}

func FindUsers(condition,limit,order string)(user []T_user,err error){
	rows,err:=pdb.Session.Query(fmt.Sprintf("select id,role,nick,pwd,avatar,phone,email,qq,status,expend,updated_at,created_at from %s %s %s %s",UserTableName(),condition,order,limit))
	if err!=nil{
		return
	}
	for rows.Next(){
		tmp:=T_user{}
		bs:=new([]byte)
		err=rows.Scan(&tmp.Id,&tmp.Role,&tmp.Nick,&tmp.Pwd,&tmp.Avatar,&tmp.Phone,&tmp.Email,&tmp.QQ,&tmp.Status,bs,&tmp.UpdatedAt,&tmp.CreatedAt)
		if err!=nil{
			fmt.Println(err.Error())
			return
		}
		tmp.Expend,err=simplejson.NewJson(*bs)
		if err!=nil{
			fmt.Println(err.Error())
			return
		}
		user=append(user,tmp)
	}
	return
}

func UpdateUsers(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s %s %s", UserTableName(), update, condition))
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


/**
获取随机字符串
 */
func  getRandomString(l int) string {
	str := "QWERTYUIOPASDFGHJKLZXCVBNM"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
