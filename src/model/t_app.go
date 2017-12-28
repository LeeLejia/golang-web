package model

import (
	"time"
	"fmt"
	"github.com/bitly/go-simplejson"
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
"developer" int4,
"valid" bool DEFAULT TRUE,
"file" varchar(256) COLLATE "default",
"src" varchar(256) COLLATE "default",
"expend" jsonb NOT NULL,
"download_count" int4 DEFAULT -1,
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */
type T_app struct {
	ID            int64            `json:"id"`
	Icon          string           `json:"icon"`
	AppId         string           `json:"app_id"`
	Name          string           `json:"name"`
	Version       string           `json:"version"`
	Describe      string           `json:"describe"`
	Developer     int              `json:"developer"`
	Valid         bool             `json:"valid"`
	File          string           `json:"file"`
	Src           string           `json:"src"`
	Expend        *simplejson.Json `json:"expend"`
	DownloadCount int              `json:"download_count"`
	CreatedAt     time.Time        `json:"created_at"`
}
var AppModel DbModel
func InitApp(){
	sc:=SqlController {
		TableName:      "t_app",
		InsertColumns:  []string{"icon","app_id","version","name","describe","developer","valid","file","src","expend","download_count","created_at"},
		QueryColumns:   []string{"id","icon","app_id","version","name","describe","developer","valid","file","src","expend","download_count","created_at"},
		InSertFields:   InsertFields,
		QueryField2Obj: QueryField2Obj,
	}
	AppModel,err:=GetModel(sc)
	if err!=nil{
		fmt.Println(err.Error())
		return
	}
}
func InsertFields(obj interface{}) []interface{} {
	fmt.Println(obj)
	app:=obj.(*T_app)
	expend := []byte{}
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
func QueryField2Obj(fields []interface{}) interface{} {
	expend,_:=simplejson.NewJson(GetByteArr(fields[10]))
	app:=&T_app{
		ID:GetInt64(fields[0],0),
		Icon:GetString(fields[1]),
		AppId:GetString(fields[2]),
		Name:GetString(fields[4]),
		Version:GetString(fields[3]),
		Describe:GetString(fields[5]),
		Developer:GetInt(fields[6],-1),
		Valid:GetBool(fields[7],true),
		File:GetString(fields[8]),
		Src:GetString(fields[9]),
		Expend:expend,
		DownloadCount:GetInt(fields[11],0),
		CreatedAt:GetTime(fields[12],time.Now()),
	}
	return app
}