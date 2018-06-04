package model

import (
	"time"
	"github.com/bitly/go-simplejson"
	m "github.com/cjwddz/fast-model"
)

/*
DROP TABLE IF EXISTS t_user;
CREATE TABLE t_user (
"id" serial NOT NULL,
"account" varchar(128) COLLATE "default",
"pwd" varchar(128) COLLATE "default",
"role" varchar(16) COLLATE "default",
"nick" varchar(128) COLLATE "default",
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
	USER_ROLE_EMPLOYER = "employer"

	USER_STATUS_INVALID=0
	USER_STATUS_WAITING_VERIFY = 2
	USER_STATUS_VALID=1

)
// 用戶表
type T_user struct {
	Id        int64                `json:"id"`
	Pwd       string             `json:"pwd"`
	Role      string             `json:"role"`
	Nick      string             `json:"nick"`
	Avatar    string             `json:"avatar"`
	Phone     string             `json:"phone"`
	Email     string             `json:"email"`
	QQ        string             `json:"qq"`
	Status    int             	 `json:"status"`
	Expend    * simplejson.Json  `json:"expend"`
	UpdatedAt time.Time          `json:"updated_at"`
	CreatedAt time.Time          `json:"created_at"`
}
func GetUserModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_user",
		InsertColumns:  []string{"pwd","role","nick","avatar","phone","email","qq","status","expend","updated_at","created_at"},
		QueryColumns:   []string{"id","pwd","role","nick","avatar","phone","email","qq","status","expend","updated_at","created_at"},
		InSertFields:   insertUserFields,
		QueryField2Obj: queryUserField2Obj,
	}
	return m.GetModel(sc)
}

func insertUserFields(obj interface{}) []interface{} {
	user :=obj.(T_user)
	expend := []byte("{}")
	if user.Expend != nil {
		bs, err := user.Expend.MarshalJSON()
		if err==nil{
			expend = bs
		}
	}
	return []interface{}{
		user.Pwd, user.Role, user.Nick, user.Avatar, user.Phone, user.Email, user.QQ, user.Status, expend, user.UpdatedAt, user.CreatedAt,
	}
}
func queryUserField2Obj(fields []interface{}) interface{} {
	expend,_:=simplejson.NewJson(m.GetByteArr(fields[9]))
	user :=T_user{
		Id:m.GetInt64(fields[0],0),
		Pwd:m.GetString(fields[1]),
		Role:m.GetString(fields[2]),
		Nick:m.GetString(fields[3]),
		Avatar:m.GetString(fields[4]),
		Phone:m.GetString(fields[5]),
		Email:m.GetString(fields[6]),
		QQ:m.GetString(fields[7]),
		Status:m.GetInt(fields[8],USER_STATUS_INVALID),
		Expend:expend,
		UpdatedAt:m.GetTime(fields[10],time.Now()),
		CreatedAt:m.GetTime(fields[11],time.Now()),
	}
	return user
}