package model

import (
	"time"
	"github.com/bitly/go-simplejson"
	m "github.com/cjwddz/fast-model"
)

/*
DROP TABLE IF EXISTS t_publish;
CREATE TABLE t_publish (
"id" serial NOT NULL,
"state" int4 DEFAULT -1,
"owner" text,
"name" varchar(128) COLLATE "default",
"describe" varchar(511) COLLATE "default",
"money_lower" int4 DEFAULT -1,
"money_upper" int4 DEFAULT -1,
"outsourcing" bool DEFAULT TRUE,
"labels" varchar(128) COLLATE "default",
"commission" varchar(15) COLLATE "default",
"need_code" bool DEFAULT FALSE,
"annex" jsonb NOT NULL,
"from_time" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
"to_time" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
"update_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */
type T_publish struct {
	ID          int64            `json:"id"`
	Owner       string           `json:"owner"`
	Name        string           `json:"name"`
	Describe    string           `json:"describe"`
	MoneyLow    int              `json:"money_lower"`
	MoneyUp     int              `json:"money_upper"`
	OutSourcing bool             `json:"outsourcing"`
	Labels      string           `json:"labels"`
	Commission  string           `json:"commission"`
	NeedCode    bool             `json:"need_code"`
	Annex       *simplejson.Json `json:"annex"`
	FromTime    time.Time        `json:"from_time"`
	ToTime      time.Time        `json:"to_time"`
	UpdateTime  time.Time        `json:"update_at"`
	CreatedTime time.Time        `json:"created_at"`
}

func GetPublishModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_publish",
		InsertColumns:  []string{"owner","name","describe","money_lower","money_upper","outsourcing","labels","commission","need_code","annex","from_time","to_time","update_at","created_at"},
		QueryColumns:   []string{"id","owner","name","describe","money_lower","money_upper","outsourcing","labels","commission","need_code","annex","from_time","to_time","update_at","created_at"},
		InSertFields:   insertPublishFields,
		QueryField2Obj: queryPublishField2Obj,
	}
	return m.GetModel(sc)
}

func insertPublishFields(obj interface{}) []interface{} {
	publish:=obj.(T_publish)
	Annex := []byte("{}")
	if publish.Annex != nil {
		bs, err := publish.Annex.MarshalJSON()
		if err==nil{
			Annex = bs
		}
	}
	return []interface{}{
		publish.Owner, publish.Name, publish.Describe, publish.MoneyLow, publish.MoneyUp, publish.OutSourcing, publish.Labels, publish.Commission, publish.NeedCode, Annex, publish.FromTime, publish.ToTime, publish.UpdateTime, publish.CreatedTime,
	}
}
func queryPublishField2Obj(fields []interface{}) interface{} {
	annex, _ := simplejson.NewJson(m.GetByteArr(fields[10]))
	publish := T_publish{
		ID:          m.GetInt64(fields[0], 0),
		Owner:       m.GetString(fields[1]),
		Name:        m.GetString(fields[2]),
		Describe:    m.GetString(fields[3]),
		MoneyLow:    m.GetInt(fields[4], 0),
		MoneyUp:     m.GetInt(fields[5], 0),
		OutSourcing: m.GetBool(fields[6], true),
		Labels:      m.GetString(fields[7]),
		Commission:  m.GetString(fields[8]),
		NeedCode:    m.GetBool(fields[9], false),
		Annex:       annex,
		FromTime:    m.GetTime(fields[11], time.Now()),
		ToTime:      m.GetTime(fields[12], time.Now()),
		UpdateTime:  m.GetTime(fields[13], time.Now()),
		CreatedTime: m.GetTime(fields[14], time.Now()),
	}
	return publish
}