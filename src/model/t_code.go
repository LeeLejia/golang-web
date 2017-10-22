package model

import (
	"time"
	"fmt"
	"../pdb"
	"github.com/bitly/go-simplejson"
)

/*
DROP TABLE IF EXISTS t_code;
CREATE TABLE t_code (
"id" serial NOT NULL,
"code" varchar(16) UNIQUE COLLATE "default",
"app_id" varchar(16) COLLATE "default",
"developer" int4,
"consumer" jsonb NOT NULL,
"describe" varchar(255) COLLATE "default",
"valid" bool DEFAULT TRUE,
"machine_count" int4 DEFAULT -1,
"enable_time" bool DEFAULT FALSE,
"start_time" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
"end_time" timestamp(6) DEFAULT CURRENT_TIMESTAMP,
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);

INSERT INTO "public"."t_code" VALUES ('1', '1111', 'old1001', '1001', '{}', 'NO', 't', '-1', 'f', '2017-10-22 08:45:20', '2017-10-22 08:45:23', '2017-03-16 09:53:04');
INSERT INTO "public"."t_code" VALUES ('2', '1523', 'old1001', '1001', '{}', 'weixin_5-15', 't', '-1', 'f', '2017-10-22 08:45:26', '2017-10-22 08:45:28', '2017-10-21 10:01:01');
INSERT INTO "public"."t_code" VALUES ('3', '12586', 'old1001', '1001', '{}', 'No', 't', '-1', 'f', '2017-10-22 08:45:32', '2017-10-22 08:45:34', '2017-03-16 10:01:56');
INSERT INTO "public"."t_code" VALUES ('4', '112235', 'old1001', '1001', '{}', 'wechat_dayuxing', 't', '-1', 'f', '2017-10-22 08:45:37', '2017-10-22 08:45:39', '2017-03-19 10:03:21');
INSERT INTO "public"."t_code" VALUES ('5', '112237', 'old1001', '1001', '{}', 'wechatextract_kabeinanren', 't', '-1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('6', '1688', 'old1002', '1001', '{}', 'qq2602447159王客户', 't', '-1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('7', '1ab2c', 'old1002', '1001', '{}', 'dayuxing', 't', '-1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('8', '2bc34', 'old1002', '1001', '{}', 'wechat_some', 't', '1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('9', '454hk', 'old1002', '1001', '{}', 'qq2634667590', 't', '-1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('10', 'ap12s5', 'old1002', '1001', '{}', 'qq736749803', 't', '1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('11', 'de98ag8', 'old1002', '1001', '{}', 'QQ95415081', 't', '-1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('12', 'eks9', 'old1002', '1001', '{}', 'qq927631708', 't', '-1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('13', 'n34ik', 'old1002', '1001', '{}', 'QQ519050901,咖啡男人', 't', '-1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('14', 'n3i45', 'old1002', '1001', '{}', 'QQ454999156', 'f', '-1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');
INSERT INTO "public"."t_code" VALUES ('15', 's23i21', 'old1002', '1001', '{}', '', 't', '-1', 'f', '2017-10-22 08:45:42', '2017-10-22 08:45:45', '2017-03-19 10:04:25');

 */
type T_code struct {
	ID           int64            `json:"id"`
	Code         string           `json:"code"`
	AppId        string           `json:"app_id"`
	Developer    int              `json:"developer"`
	Consumer     *simplejson.Json `json:"consumer"`
	Describe     string           `json:"describe"`
	Valid        bool             `json:"valid"`
	MachineCount int              `json:"machine_count"`
	EnableTime   bool             `json:"enable_time"`
	StartTime    time.Time        `json:"start_time"`
	EndTime      time.Time        `json:"end_time"`
	CreatedAt    time.Time        `json:"created_at"`
}
func CodeTableName() string {
	return "t_code"
}

func (c *T_code) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(code,app_id,developer,consumer,describe,valid,machine_count,enable_time,start_time,end_time,created_at) "+
			  "VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)", CodeTableName()))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	c.CreatedAt = time.Now()
	consumer:="{}"
	if c.Consumer!=nil{
		bs,err:=c.Consumer.MarshalJSON()
		if err!=nil{
			return err
		}
		consumer=string(bs)
	}
	_, err = stmt.Exec(c.Code,c.AppId, c.Developer,consumer,c.Describe,c.Valid,c.MachineCount,c.EnableTime,c.StartTime,c.EndTime,c.CreatedAt)
	return
}

func FindCodes(condition, limit, order string) (result []T_code,err error) {
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,code,app_id,developer,consumer,describe,valid,machine_count," +
			  "enable_time,start_time,end_time,created_at FROM %s %s %s %s", CodeTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}
	for rows.Next() {
		tmp := T_code{}
		bs:=new([]byte)
		err = rows.Scan(&tmp.ID, &tmp.Code,&tmp.AppId, &tmp.Developer, bs,&tmp.Describe,&tmp.Valid,&tmp.MachineCount,&tmp.EnableTime,&tmp.StartTime,&tmp.EndTime,&tmp.CreatedAt)
		tmp.Consumer,_=simplejson.NewJson(*bs)
		if err==nil {
			result = append(result, tmp)
		}
	}
	return result, err
}

func CountCodes(condition string) (count int, err error) {
	count = 0
	err = pdb.Session.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s %s", CodeTableName(), condition)).Scan(&count)
	return
}

func UpdateCodes(update, condition string) (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("UPDATE %s %s %s", CodeTableName(), update, condition))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = stmt.Exec()
	return
}
