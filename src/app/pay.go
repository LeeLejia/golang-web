package app

import (
	"../common"
	"net/http"
	"sync"
	"../model"
	"../common/log"
	fm "github.com/cjwddz/fast-model"
	"time"
	"strconv"
	"fmt"
)
var GoodLock sync.Mutex
/**
购买支付接口
 */
func Pay(sess *common.Session, w http.ResponseWriter, r *http.Request) {
	// return 后会带上参数:money=0.11&name=二手书一本&out_trade_no=123456&pid=7266&trade_no=2018060222421456619&trade_status=TRADE_SUCCESS&type=alipay&sign=ad76003fb5de12af6531e79090733715&sign_type=MD5
	// 父子窗通信
	r.ParseForm()
	goodType:=r.Form.Get("goodType")
	if goodType == ""{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID, "请提供商品类型")
		return
	}
	result,err:=GoodModel.Query(fm.DbCondition{}.And("=","type",goodType).Limit(1,0))
	if err!=nil{
		log.E("Pay数据库查询出错",sess.User.Email,"reason:%s",err.Error())
		common.ReturnEFormat(w,common.CODE_DB_RW_ERR, "数据库查询出错")
		return
	}
	if result==nil || len(result)<=0{
		log.W("Pay商品类型不存在",sess.User.Email,"商品类型不存在")
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID, "商品类型不存在")
		return
	}
	good:= result[0].(model.T_Goods)
	orderId:=good.Type + common.BytesToInt([]byte(strconv.FormatInt(time.Now().Unix(),32)))
	// 提取用户附加信息
	fmt.Printf("orderId:%s",orderId)
	// check
	url,err:=common.GetAliPayLink(good.Price,good.Name,orderId,"https://www.cjwddz.cn/blog/")
	if err!=nil{
		common.ReturnEFormat(w,common.CODE_SERVICE_ERR, err.Error())
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"url":url,"msg":"操作成功！"})
}