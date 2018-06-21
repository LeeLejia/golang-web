package model

import (
	"time"
	m "github.com/cjwddz/fast-model"
)

/*
DROP TABLE IF EXISTS t_file;
CREATE TABLE t_vk (
"key" varchar(16) UNIQUE COLLATE "default",
"value" text,
"update_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */
type T_Vk struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Owner     string    `json:"owner,omitempty"`
	UpdateAt  time.Time `json:"update_at"`
	CreatedAt time.Time `json:"created_at"`
}

func GetVkModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_vk",
		InsertColumns:  []string{"key","value"},
		QueryColumns:   []string{"key","value","update_at","created_at"},
		InSertFields:   insertVkFields,
		QueryField2Obj: queryVkField2Obj,
	}
	return m.GetModel(sc)
}

func insertVkFields(obj interface{}) []interface{} {
	vk:=obj.(T_Vk)
	return []interface{}{
		vk.Key, vk.Value,
	}
}
func queryVkField2Obj(fields []interface{}) interface{} {
	file := T_Vk{
		Key:       m.GetString(fields[1]),
		Value:     m.GetString(fields[2]),
		UpdateAt:  m.GetTime(fields[9], time.Now()),
		CreatedAt: m.GetTime(fields[9], time.Now()),
	}
	return file
}