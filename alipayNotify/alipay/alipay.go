package alipay

import (
	"net/http"
	"io/ioutil"
	"time"
)

var client = http.Client{}
var qrCodeCache = []byte("尚未初始化二维码页面！请再次刷新.")
var qrCodeScanOk = false
const(
	qrCode = "https://auth.alipay.com/login/homeB.htm?redirectType=parent"
)
func init(){
	http.HandleFunc("/alipay", func(writer http.ResponseWriter, r *http.Request) {
		if qrCodeScanOk{
			writer.Write([]byte("支付宝监控系统正常运行中.."))
			return
		}
		preToScanQRCode()
		writer.Write(qrCodeCache)
		go listenScanEven()
	})
}
// 准备支付宝扫描页面
func preToScanQRCode(){
	req, err := http.NewRequest(http.MethodGet, qrCode, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.89 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		qrCodeCache = []byte(err.Error())
		return
	}
	qrCodeCache, err = ioutil.ReadAll(resp.Body)
}
// 监听扫描事件
func listenScanEven(){
	t:=0
	for !qrCodeScanOk{
		select{
			case <-time.After(2000):
				// 2秒检查一次扫描结果
				if t++;t%30==0{
					// 更新页面，防止过时
					preToScanQRCode()
					continue
				}
				// 监听扫描事件
				// todo
			}
		}
	}
