package model

import (
	"time"
	m "github.com/cjwddz/fast-model"
	"github.com/bitly/go-simplejson"
)

/*
DROP TABLE IF EXISTS t_orders;
CREATE TABLE t_orders (
"id" serial NOT NULL,
"type" varchar(16) COLLATE "default",
"channel" varchar(16) COLLATE "default",
"order_id" varchar(128) UNIQUE COLLATE "default",
"name" varchar(128) COLLATE "default",
"price" int,
"state" int,
"owner" varchar(128) COLLATE "default",
"expend" jsonb NOT NULL,
"created_at" timestamp(6) DEFAULT CURRENT_TIMESTAMP
)
WITH (OIDS=FALSE);
 */
type T_Order struct {
	ID        int64            `json:"id"`
	Type      string           `json:"type"`
	Channel	  string           `json:"channel"`
	OrderId	  string		   `json:"order_id"`
	Name      string           `json:"name"`
	Price     int              `json:"price"`
	State     int              `json:"state"`
	Owner     string           `json:"owner,omitempty"`
	Expend    *simplejson.Json `json:"expend"`
	CreatedAt time.Time        `json:"created_at"`
}
const (
	GOOD_ORDER_STATE_SUCCESS = 1
	GOOD_ORDER_STATE_WAITTING_PAY = 2
	GOOD_ORDER_STATE_PAY_FAILED = -1
)
func GetOrderModel() (m.DbModel, error){
	sc:=m.SqlController {
		TableName:      "t_orders",
		InsertColumns:  []string{"type","channel","order_id","name","price","state","owner","expend","created_at"},
		QueryColumns:   []string{"id","type","channel","order_id","name","price","state","owner","expend","created_at"},
		InSertFields:   insertOrderFields,
		QueryField2Obj: queryOrderField2Obj,
	}
	return m.GetModel(sc)
}

func insertOrderFields(obj interface{}) []interface{} {
	order :=obj.(T_Order)
	expend := []byte("{}")
	if order.Expend != nil {
		bs, err := order.Expend.MarshalJSON()
		if err==nil{
			expend = bs
		}
	}
	return []interface{}{
		order.Type, order.Channel, order.OrderId, order.Name, order.Price, order.State, order.Owner, expend, order.CreatedAt,
	}
}

func queryOrderField2Obj(fields []interface{}) interface{} {
	expend,_:=simplejson.NewJson(m.GetByteArr(fields[8]))
	order := T_Order{
		ID:        m.GetInt64(fields[0], 0),
		Type:      m.GetString(fields[1]),
		Channel:   m.GetString(fields[2]),
		OrderId:   m.GetString(fields[3]),
		Name:      m.GetString(fields[4]),
		Price:     m.GetInt(fields[5], 0),
		State:     m.GetInt(fields[6], 0),
		Owner:     m.GetString(fields[7]),
		Expend:    expend,
		CreatedAt: m.GetTime(fields[9], time.Now()),
	}
	return order
}