package model

import (
	m "github.com/cjwddz/fast-model"
)

type Task struct{
	Id      string
	Uid     string
	Content string
	Money   string
	Label   string
	Time    string
}

func GetTaskModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_task",
		InsertColumns:  []string{"id","uid","content","money","time"},
		QueryColumns:   []string{"id","uid","content","money","time"},
		InSertFields:   insertTaskFields,
		QueryField2Obj: queryTaskField2Obj,
	}
	return m.GetModel(sc)
}

func insertTaskFields(obj interface{}) []interface{} {
	task :=obj.(Task)
	return []interface{}{
		task.Id, task.Uid, task.Content, task.Money, task.Time,
	}
}
func queryTaskField2Obj(fields []interface{}) interface{} {
	task :=Task{
		Id:      m.GetString(fields[0]),
		Uid:     m.GetString(fields[1]),
		Content: m.GetString(fields[2]),
		Money:   m.GetString(fields[3]),
		Time:    m.GetString(fields[4]),
	}
	return task
}