package model

import (
	"time"
	m "github.com/cjwddz/fast-model"
)
/*
DROP TABLE IF EXISTS t_log;
CREATE TABLE t_log (
"id" serial NOT NULL,
"type" varchar(16) COLLATE "default",
"tag" varchar(256) COLLATE "default",
"operator" varchar(128) COLLATE "default",
"content" varchar(512) COLLATE "default",
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */

const(
	LOG_TYPE_DEBUG="debug"
	LOG_TYPE_INFO="info"
	LOG_TYPE_WARM="warn"
	LOG_TYPE_ERROR="error"
	LOG_TYPE_NORMAL="normal"
)

type T_log struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Tag       string    `json:"tag"`
	Operator  string    `json:"operator"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func GetLogModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_log",
		InsertColumns:  []string{"type","tag","operator","content","created_at"},
		QueryColumns:   []string{"id","type","tag","operator","content","created_at"},
		InSertFields:   insertLogFields,
		QueryField2Obj: queryLogField2Obj,
	}
	return m.GetModel(sc)
}
func insertLogFields(obj interface{}) []interface{} {
	log :=obj.(T_log)
	return []interface{}{
		log.Type, log.Tag, log.Operator, log.Content, log.CreatedAt,
	}
}
func queryLogField2Obj(fields []interface{}) interface{} {
	log:=T_log{
		ID:m.GetInt64(fields[0],0),
		Type:m.GetString(fields[1]),
		Tag:m.GetString(fields[2]),
		Operator:m.GetString(fields[3]),
		Content:m.GetString(fields[4]),
		CreatedAt:m.GetTime(fields[5],time.Now()),
	}
	return log
}