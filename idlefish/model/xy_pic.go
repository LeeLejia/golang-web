package model

import (
	m "github.com/cjwddz/fast-model"
	//"time"
)
/*
DROP TABLE IF EXISTS xy_task;
CREATE TABLE xy_pic (
"file" varchar(64) COLLATE "default",
)
WITH (OIDS=FALSE);
 */
type Pic struct {
	Name       string `json:"file"`
}

func GetPicModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "xy_pic",
		InsertColumns:  []string{"file"},
		QueryColumns:   []string{"file"},
		InSertFields:   insertPicFields,
		QueryField2Obj: queryPicField2Obj,
	}
	return m.GetModel(sc)
}

func insertPicFields(obj interface{}) []interface{} {
	pic :=obj.(Pic)
	return []interface{}{
		pic.Name,
	}
}
func queryPicField2Obj(fields []interface{}) interface{} {
	pic :=Pic{
		Name:        m.GetString(fields[0]),
	}
	return pic
}