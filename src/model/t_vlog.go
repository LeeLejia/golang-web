package model

import (
	"time"
	m "github.com/cjwddz/fast-model"
)

/*
DROP TABLE IF EXISTS t_vlog;
CREATE TABLE t_vlog (
"id" serial NOT NULL,
"tag" varchar(255) COLLATE "default",
"code" varchar(16) COLLATE "default",
"app_id" varchar(16) COLLATE "default",
"machine" varchar(128) COLLATE "default",
"content" varchar(255) COLLATE "default",
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */
const (
	VERIFY_LOG_SUCCESS         = "验证成功"

	VERIFY_LOG_NOEXIST         = "邀请码不存在"
	VERIFY_LOG_INVALID         = "邀请码已经禁用"
	VERIFY_LOG_NOFIND_APP      = "找不到指定app"
	VERIFY_LOG_APP_NOVALID     = "该APP已失效"
	VERIFY_LOG_EMPTY           = "邀请码为空"
	VERIFY_LOG_MACHINE_NOVALID = "机器码不可用"
	VERIFY_LOG_TOOMUCH_MACHINE = "设备数超限制"
	VERIFY_LOG_BEFORE_TIME     = "尚未生效"
	VERIFY_LOG_AFTER_TIME      = "已经过期"

	VERIFY_LOG_PROTO_NOVALID = "协议码不可用"
)

type T_vlog struct {
	ID        int64     `json:"id"`
	Tag       string    `json:"tag"`
	Code      string    `json:"code"`
	App 	  string 	`json:"app_id"`
	Machine   string    `json:"machine"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func GetVLogModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_vlog",
		InsertColumns:  []string{"tag","code","app_id","machine","content","created_at"},
		QueryColumns:   []string{"id","tag","code","app_id","machine","content","created_at"},
		InSertFields:   insertVLogFields,
		QueryField2Obj: queryVLogField2Obj,
	}
	return m.GetModel(sc)
}

func insertVLogFields(obj interface{}) []interface{} {
	vlog:=obj.(T_vlog)
	return []interface{}{
		vlog.Tag,vlog.Code,vlog.App,vlog.Machine,vlog.Content,vlog.CreatedAt,
	}
}
func queryVLogField2Obj(fields []interface{}) interface{} {
	vl:=T_vlog{
		ID:m.GetInt64(fields[0],0),
		Tag:m.GetString(fields[1]),
		Code:m.GetString(fields[2]),
		App:m.GetString(fields[3]),
		Machine:m.GetString(fields[4]),
		Content:m.GetString(fields[5]),
		CreatedAt:m.GetTime(fields[6],time.Now()),
	}
	return vl
}