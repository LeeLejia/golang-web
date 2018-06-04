package app

import (
	"../common"
	"net/http"
	"../common/log"
	fm "github.com/cjwddz/fast-model"
	"../model"
	"fmt"
	"time"
	"strconv"
	"github.com/bitly/go-simplejson"
)

/**
获取商品
 */
func GetGoods(sess *common.Session, w http.ResponseWriter, r *http.Request){
	cond:=fm.DbCondition{}.And2(r,"=","i_type").And2(r,"=","s_channel").And2(r,"like","s_name").And2(r,"=","i_state")
	total,err:=GoodModel.Count(cond)
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("GetGoods出错",sess.User.Email,err.Error())
		return
	}
	result,err:=GoodModel.Query(cond.Limit2(r,"start","count").Order("order by created_at desc"))
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("GetGoods出错",sess.User.Email,err.Error())
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"goods": result,"total":total})
	log.N("GetGoods",sess.User.Email,fmt.Sprintf("listCount=%d,total=%d",len(result),total))
	return
}

/**
添加商品
 */
func AddGood(sess *common.Session, w http.ResponseWriter, r *http.Request){
	if !common.IsRole(sess.User.Role,model.USER_ROLE_ADMIN){
		log.W(common.ACTION_VIOLENCE,sess.User.Email,"该用户roles=%s,尝试添加商品.",sess.User.Role)
		common.ReturnEFormat(w,common.CODE_ROLE_INVADE, "抱歉,您没有添加商品的权限.")
		return
	}
	r.ParseForm()
	name:=r.Form.Get("name")
	channel:= r.Form.Get("channel")
	_price:=r.Form.Get("price")
	_state:=r.Form.Get("state")
	_count:=r.Form.Get("count")
	var state int
	if _state == "invalid"{
		state = model.GOOD_TYPE_INVALID
	}else{
		state = model.GOOD_TYPE_VALID
	}
	count,err:=strconv.Atoi(_count)
	if err!=nil || count< -1 {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "请设置有效数量")
		return
	}
	price,err:=strconv.Atoi(_price)
	if err!=nil || price<=0 {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "请设置有效价格")
		return
	}
	expend:= simplejson.New()
	expend.Set("desc",r.Form.Get("desc"))
	expend.Set("picture",r.Form.Get("picture"))
	good:=model.T_Goods{
		Channel:   channel,
		Name:      name,
		Price:     price,
		State:     state,
		Owner:     sess.User.Email,
		Count:     count,
		Expend:    expend,
		UpdateAt:  time.Now(),
		CreatedAt: time.Now(),
	}
	err=GoodModel.Insert(good)
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("AddGood出错",sess.User.Email,err.Error())
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"msg":"添加商品成功!"})
	log.N("AddGood成功",sess.User.Email,"添加商品:channel=%s,name=%s,price=%s",channel,name,_price)
}