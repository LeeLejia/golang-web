package model

import (
	"time"
	"fmt"
	"../pdb"
	"github.com/bitly/go-simplejson"
)

type T_code struct {
	ID           int64    `json:"id"`
	Code         string   `json:"code"`
	AppId        string   `json:"app_id"`
	Developer    int   `json:"developer"`
	Consumer     * simplejson.Json   `json:"consumer"`
	Describe     string   `json:"describe"`
	Valid        bool     `json:"valid"`
	MachineCount int      `json:"machine_count"`
	MostCount    int      `json:"most_count"`
	EnableTime   bool       `json:"enable_time"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	CreatedAt    time.Time `json:"created_at"`
}
func CodeTableName() string {
	return "t_code"
}

func (c *T_code) Insert() (err error) {
	stmt, err := pdb.Session.Prepare(fmt.Sprintf("INSERT INTO %s(code,app_id,developer,consumer,describe,valid,machine_count,most_count,enable_time,start_time,end_time,created_at) "+
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
	_, err = stmt.Exec(c.Code,c.AppId, c.Developer,consumer,c.Describe,c.Valid,c.MachineCount,c.MostCount,c.EnableTime,c.StartTime,c.EndTime,c.CreatedAt)
	return
}

func FindCodes(condition, limit, order string) (result []T_code,err error) {
	rows, err := pdb.Session.Query(fmt.Sprintf("SELECT id,code,app_id,developer,consumer,describe,valid,machine_count,most_count," +
			  "enable_time,start_time,end_time,created_at FROM %s %s %s %s", CodeTableName(), condition, order, limit))
	if err != nil {
		return result, err
	}
	for rows.Next() {
		tmp := T_code{}
		bs:=new([]byte)
		err = rows.Scan(&tmp.ID, &tmp.Code,&tmp.AppId, &tmp.Developer, bs,&tmp.Describe,&tmp.Valid,&tmp.MachineCount,&tmp.MostCount,&tmp.EnableTime,&tmp.StartTime,&tmp.EndTime,&tmp.CreatedAt)
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
