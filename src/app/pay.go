package app

import (
	"../common"
	"net/http"
)

/**
购买支付接口
 */
func Pay(_ *common.Session, w http.ResponseWriter, r *http.Request) {
	// return 后会带上参数:money=0.11&name=二手书一本&out_trade_no=123456&pid=7266&trade_no=2018060222421456619&trade_status=TRADE_SUCCESS&type=alipay&sign=ad76003fb5de12af6531e79090733715&sign_type=MD5
	// 父子窗通信

	url,err:=common.GetAliPayLink(0.1056,"二手书一本","https://www.cjwddz.cn/blog/")
	if err!=nil{
		common.ReturnEFormat(w,common.CODE_SERVICE_ERR, err.Error())
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"url":url,"msg":"操作成功！"})
}