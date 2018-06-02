/*
	通用类
 */
package common

import (
	"crypto/md5"
	"fmt"
	"../common/conf"
)

/**
获取支付链接
 */
func GetAliPayLink(money float32,commodity string,return_url string) (url string, err error) {
	notify_url:="https://www.cjwddz.cn/api/register"
	out_trade_no:="123456"
	pid:=conf.App.PayPid
	key:=conf.App.PayKey
	sitename:=conf.App.SiteName
	payType:="alipay"
	rawParams := fmt.Sprintf("money=%0.2f&name=%s&notify_url=%s&out_trade_no=%s&pid=%d&return_url=%s&sitename=%s&type=%s",money,commodity,notify_url,out_trade_no,pid,return_url,sitename,payType)
	ciphertext := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s",rawParams,key))))
	return fmt.Sprintf("http://pay.sddyun.cn/submit.php?%s&sign_type=MD5&sign=%s",rawParams,ciphertext),nil
}

