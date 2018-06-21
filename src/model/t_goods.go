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
"type" varchar(128),
"channel" varchar(16) COLLATE "default",
"name" varchar(128) COLLATE "default",
"price" int,
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
	Channel	  string           `json:"channel"`
	Name      string           `json:"name"`
	Price     int              `json:"price"`
	State     int              `json:"state"`
	Owner     string           `json:"owner,omitempty"`
	Count     int              `json:"count"`
	Expend    *simplejson.Json `json:"expend"`
	UpdateAt  time.Time        `json:"update_at"`
	CreatedAt time.Time        `json:"created_at"`
}
const (
	GOOD_TYPE_INVALID = -1
	GOOD_TYPE_VALID = 1
)
func GetGoodModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_goods",
		InsertColumns:  []string{"type", "channel","name","price","state","owner","count","expend","update_at","created_at"},
		QueryColumns:   []string{"id","type","channel","name","price","state","owner","count","expend","update_at","created_at"},
		InSertFields:   insertGoodFields,
		QueryField2Obj: queryGoodField2Obj,
	}
	return m.GetModel(sc)
}

func insertGoodFields(obj interface{}) []interface{} {
	good :=obj.(T_Goods)
	expend := []byte("{}")
	if good.Expend != nil {
		bs, err := good.Expend.MarshalJSON()
		if err==nil{
			expend = bs
		}
	}
	return []interface{}{
		good.Type, good.Channel, good.Name, good.Price, good.State, good.Owner, good.Count, expend, good.UpdateAt, good.CreatedAt,
	}
}
func queryGoodField2Obj(fields []interface{}) interface{} {
	expend,_:=simplejson.NewJson(m.GetByteArr(fields[8]))
	good := T_Goods{
		ID:        m.GetInt64(fields[0], 0),
		Type:      m.GetString(fields[1]),
		Channel:   m.GetString(fields[2]),
		Name:      m.GetString(fields[3]),
		Price:     m.GetInt(fields[4], 0),
		State:     m.GetInt(fields[5], 0),
		Owner:     m.GetString(fields[6]),
		Count:     m.GetInt(fields[7],0),
		Expend:    expend,
		UpdateAt:  m.GetTime(fields[9], time.Now()),
		CreatedAt: m.GetTime(fields[10], time.Now()),
	}
	return good
}