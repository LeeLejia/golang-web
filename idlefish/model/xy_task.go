package model

import (
	m "github.com/cjwddz/fast-model"
	"time"
)
/*
DROP TABLE IF EXISTS xy_task;
CREATE TABLE xy_task (
"machine" varchar(255) COLLATE "default",
"type" varchar(255) COLLATE "default",
"title" varchar(255) COLLATE "default",
"addr" varchar(255) COLLATE "default",
"raw_money" varchar(255) COLLATE "default",
"public_type" varchar(255) COLLATE "default",
"describe" varchar(255) COLLATE "default",
"category" varchar(255) COLLATE "default",
"money" varchar(255) COLLATE "default",
"tran_money" varchar(256) COLLATE "default",
"pond" varchar(255) COLLATE "default",
"pics" varchar(1024) COLLATE "default",
"created" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */
type Task struct {
	Machine    string `json:"machine"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Addr       string `json:"addr"`
	RawMoney   string `json:"rawMoney"`
	PublicType string `json:"publicType"`
	Desc       string `json:"describe"`
	Category   string `json:"category"`
	Money      string `json:"money"`
	TranMoney  string `json:"tranMoney"`
	Pond       string `json:"pond"`
	Pics	   string `json:"pics"`
	CreatedTime time.Time `json:"created"`
}

func GetTaskModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "xy_task",
		InsertColumns:  []string{"machine","type","title","addr","raw_money","public_type","describe","category","money","tran_money","pond","pics","created"},
		QueryColumns:   []string{"machine","type","title","addr","raw_money","public_type","describe","category","money","tran_money","pond","pics","created"},
		InSertFields:   insertTaskFields,
		QueryField2Obj: queryTaskField2Obj,
	}
	return m.GetModel(sc)
}

func insertTaskFields(obj interface{}) []interface{} {
	task :=obj.(Task)
	return []interface{}{
		task.Machine,task.Type, task.Title, task.Addr, task.RawMoney, task.PublicType,task.Desc,task.Category,task.Money,task.TranMoney,task.Pond,task.Pics,task.CreatedTime,
	}
}
func queryTaskField2Obj(fields []interface{}) interface{} {
	task :=Task{
		Machine:     m.GetString(fields[0]),
		Type:        m.GetString(fields[1]),
		Title:       m.GetString(fields[2]),
		Addr:        m.GetString(fields[3]),
		RawMoney:    m.GetString(fields[4]),
		PublicType:  m.GetString(fields[5]),
		Desc:        m.GetString(fields[6]),
		Category:    m.GetString(fields[7]),
		Money:       m.GetString(fields[8]),
		TranMoney:   m.GetString(fields[9]),
		Pond:        m.GetString(fields[10]),
		Pics: 		 m.GetString(fields[11]),
		CreatedTime: m.GetTime(fields[12],time.Now()),
	}
	return task
}