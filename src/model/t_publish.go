package model

import (
	"time"
	"github.com/bitly/go-simplejson"
	m "github.com/cjwddz/fast-model"
)

/*
DROP TABLE IF EXISTS t_app;
CREATE TABLE t_app (
"id" serial NOT NULL,
"icon" text,
"app_id" varchar(16) COLLATE "default",
"name" varchar(128) COLLATE "default",
"version" varchar(16) COLLATE "default",
"describe" varchar(255) COLLATE "default",
"developer" varchar(128) COLLATE "default",
"valid" bool DEFAULT TRUE,
"file" varchar(256) COLLATE "default",
"src" varchar(256) COLLATE "default",
"expend" jsonb NOT NULL,
"download_count" int4 DEFAULT -1,
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */
type T_publish struct {
	ID            int64            `json:"id"`
	Icon          string           `json:"icon"`
	AppId         string           `json:"app_id"`
	Name          string           `json:"name"`
	Version       string           `json:"version"`
	Describe      string           `json:"describe"`
	Developer     string           `json:"developer"`
	Valid         bool             `json:"valid"`
	File          string           `json:"file"`
	Src           string           `json:"src"`
	Expend        *simplejson.Json `json:"expend"`
	DownloadCount int              `json:"download_count"`
	CreatedAt     time.Time        `json:"created_at"`
}

func GetPublishModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_app",
		InsertColumns:  []string{"icon","app_id","version","name","describe","developer","valid","file","src","expend","download_count","created_at"},
		QueryColumns:   []string{"id","icon","app_id","version","name","describe","developer","valid","file","src","expend","download_count","created_at"},
		InSertFields:   insertAppFields,
		QueryField2Obj: queryAppField2Obj,
	}
	return m.GetModel(sc)
}

func insertPublishFields(obj interface{}) []interface{} {
	app:=obj.(T_app)
	expend := []byte("{}")
	if app.Expend != nil {
		bs, err := app.Expend.MarshalJSON()
		if err==nil{
			expend = bs
		}
	}
	return []interface{}{
		app.Icon,app.AppId,app.Version,app.Name,app.Describe,app.Developer,app.Valid,app.File,app.Src,expend,app.DownloadCount,app.CreatedAt,
	}
}
func queryPublishField2Obj(fields []interface{}) interface{} {
	expend,_:=simplejson.NewJson(m.GetByteArr(fields[10]))
	app:=T_app{
		ID:m.GetInt64(fields[0],0),
		Icon:m.GetString(fields[1]),
		AppId:m.GetString(fields[2]),
		Name:m.GetString(fields[4]),
		Version:m.GetString(fields[3]),
		Describe:m.GetString(fields[5]),
		Developer:m.GetString(fields[6]),
		Valid:m.GetBool(fields[7],true),
		File:m.GetString(fields[8]),
		Src:m.GetString(fields[9]),
		Expend:expend,
		DownloadCount:m.GetInt(fields[11],0),
		CreatedAt:m.GetTime(fields[12],time.Now()),
	}
	return app
}