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
	"github.com/bitly/go-simplejson"
)
var GoodLock sync.Mutex
/**
购买支付接口
 */
func Pay(sess *common.Session, w http.ResponseWriter, r *http.Request) {
	// return 后会带上参数:money=0.11&name=二手书一本&out_trade_no=123456&pid=7266&trade_no=2018060222421456619&trade_status=TRADE_SUCCESS&type=alipay&sign=ad76003fb5de12af6531e79090733715&sign_type=MD5
	// 父子窗通信
	r.ParseForm()
	returnUrl:=r.Form.Get("returnUrl")
	goodType:=r.Form.Get("goodType")
	count,err:= strconv.Atoi(r.Form.Get("count"))
	contact:=r.Form.Get("contact")
	if contact == ""{
		contact = sess.User.Email
	}
	if err!=nil || count<0{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID, "数量输入非法")
		return
	}
	if count== 0{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID, "选择商品数量不能为0!")
		return
	}
	if goodType == ""{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID, "请提供商品类型")
		return
	}
	result,err:=GoodModel.Query(fm.DbCondition{}.And("=","type",goodType).Limit(1,0))
	if err!=nil{
		log.E("Pay",sess.User.Email,"reason:%s",err.Error())
		common.ReturnEFormat(w,common.CODE_DB_RW_ERR, "数据库查询出错")
		return
	}
	if result==nil || len(result)<=0{
		log.W("Pay",sess.User.Email,"商品类型不存在,goodType=%s",goodType)
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID, "商品类型不存在")
		return
	}
	good:= result[0].(model.T_Goods)
	if good.State == model.GOOD_TYPE_INVALID{
		log.W("Pay",sess.User.Email,"商品类型不可用,goodType=%s,goodName=%s,channel=%s",good.Type,good.Name,good.Channel)
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID, "该商品已经下架")
		return
	}
	orderId:=good.Type + common.BytesToInt([]byte(strconv.FormatInt(time.Now().Unix(),32))) + common.GetRandomInt(4)
	if good.Count != - 1 && good.Count < count {
		log.W("Pay",sess.User.Email,"商品库存数量不足,goodType=%s,goodName=%s,channel=%s,库存数量:%d,购买数量:%d",good.Type,good.Name,good.Channel,good.Count,count)
		common.ReturnEFormat(w,common.CODE_RESOURCE_SHORT, "商品数量不足")
		return
	}
	expend:= simplejson.New()
	expend.Set("count",count)
	expend.Set("contact",contact)
	// 提取用户附加信息
	order := model.T_Order{
		Type:      good.Type,
		Channel:   good.Channel,
		OrderId:   orderId,
		Name:      good.Name,
		Price:     good.Price * count,
		State:     model.GOOD_ORDER_STATE_WAITTING_PAY,
		Owner:     sess.User.Email,
		Expend:    expend,
		CreatedAt: time.Now(),
	}
	err=OrderModel.Insert(order)
	if err!=nil{
		log.E("Pay",sess.User.Email,"系统添加订单失败,reason:%s",err.Error())
		common.ReturnEFormat(w,common.CODE_DB_RW_ERR, "生成订单失败!请稍后重试.")
		return
	}
	// 提取用户附加信息
	url,err:=common.GetAliPayLink(float32(good.Price)/100,good.Name,orderId,returnUrl)
	if err!=nil{
		log.E("Pay",sess.User.Email,"reason:%s",err.Error())
		common.ReturnEFormat(w,common.CODE_SERVICE_ERR, "服务器出现错误.")
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"url":url,"msg":"操作成功！"})
}