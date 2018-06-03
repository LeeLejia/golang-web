package model

import (
	"time"
	m "github.com/cjwddz/fast-model"
	"github.com/bitly/go-simplejson"
)

/*
DROP TABLE IF EXISTS t_goods;
CREATE TABLE t_goods (
"id" serial NOT NULL,
"type" varchar(16) UNIQUE COLLATE "default",
"name" varchar(128) COLLATE "default",
"price" numeric(2),
"state" int,
"owner" varchar(128) COLLATE "default",
"count" int DEFAULT -1,
"expend" jsonb NOT NULL,
"update_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */
type T_Goods struct {
	ID        int64            `json:"id"`
	Type      string           `json:"type"`
	Name      string           `json:"name"`
	Price     float32          `json:"price"`
	State     int              `json:"state"`
	Owner     string           `json:"owner,omitempty"`
	Count     int              `json:"count"`
	Expend    *simplejson.Json `json:"expend"`
	UpdateAt  time.Time        `json:"update_at"`
	CreatedAt time.Time        `json:"created_at"`
}

func GetGoodModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_goods",
		InsertColumns:  []string{"type","name","price","state","owner","count","expend","update_at","created_at"},
		QueryColumns:   []string{"id","type","name","price","state","owner","count","expend","update_at","created_at"},
		InSertFields:   insertGoodFields,
		QueryField2Obj: queryGoodField2Obj,
	}
	return m.GetModel(sc)
}

func insertGoodFields(obj interface{}) []interface{} {
	file:=obj.(T_Goods)
	return []interface{}{
		file.Type, file.Name, file.Price, file.State, file.Owner, file.Count,file.Expend,file.UpdateAt,file.CreatedAt,
	}
}
func queryGoodField2Obj(fields []interface{}) interface{} {
	expend,_:=simplejson.NewJson(m.GetByteArr(fields[7]))
	good := T_Goods{
		ID:        m.GetInt64(fields[0], 0),
		Type:      m.GetString(fields[1]),
		Name:      m.GetString(fields[2]),
		Price:     m.GetFloat(fields[3], 0),
		State:     m.GetInt(fields[4], 0),
		Owner:     m.GetString(fields[5]),
		Count:     m.GetInt(fields[6],0),
		Expend:    expend,
		UpdateAt:  m.GetTime(fields[8], time.Now()),
		CreatedAt: m.GetTime(fields[9], time.Now()),
	}
	return good
}