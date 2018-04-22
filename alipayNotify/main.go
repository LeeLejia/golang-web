package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"io"
)

var cookies string

func main() {
	http.HandleFunc("/alipay", func(writer http.ResponseWriter, request *http.Request) {
		resp,err:=http.Get("https://auth.alipay.com/login/homeB.htm?redirectType=parent")
		if err!=nil{
			writer.Write([]byte("访问支付宝网页失败！"))
			return
		}
		io.Copy(writer,resp.Body)
	})
	cookies = `mobileSendTime=-1; credibleMobileSendTime=-1; ctuMobileSendTime=-1; riskMobileBankSendTime=-1; riskMobileAccoutSendTime=-1; riskMobileCreditSendTime=-1; riskCredibleMobileSendTime=-1; riskOriginalAccountMobileSendTime=-1; cna=fMdvErTK/j8CATsp/PcXbLCe; alipay="K1iSL1mnW5fOx3otoaxG8uvgc9fx+cet2myBCzMsxA4Ra1GwRw=="; iw.userid="K1iSL1mnW5fOx3otoaxG8g=="; UM_distinctid=16121834e7eed7-0d3cf757797ae4-3a75045d-1fa400-16121834e7f859; ctoken=M-BPXXGBzBRpPhKu; LoginForm=alipay_login_auth; CLUB_ALIPAY_COM=2088022627241688; ali_apache_tracktmp="uid=2088022627241688"; session.cookieNameId=ALIPAYJSESSIONID; ALIPAYJSESSIONID=RZ24XBGt15lIaYT6v67dc2MsJi0T0iauthRZ13GZ00; zone=GZ00C; JSESSIONID=A8038BC601D10A0E8974F92DA2CE283B; spanner=oTUd4hWCyU1K82Ra2mljNLzbC+rmhcSDXt2T4qEYgj0=; rtk=/4Rq1JGsddfSYgLx84QuBSbNHLwXRd1ztzBp015XRpdSOKlTisn`
	HttpSetup()
	Run()
}

func HttpSetup() {
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		tradeNo := r.URL.Query().Get("tradeno")
		if len(tradeNo) != 6 {
			w.Write([]byte("tradeno error"))
			w.WriteHeader(302)
			return
		}
		for _, v := range TransferMap {
			if v.TradeNo[len(v.TradeNo)-6:] == tradeNo {
				//time out
				if time.Now().Sub(v.Time) > time.Minute*5 {
					w.Write([]byte("time out"))
					//has checked
				} else if v.Examed {
					w.Write([]byte("checked"))
					//normal
				} else {
					v.Examed = true
					w.Write([]byte(fmt.Sprint(v.Amount)))
				}
				return
			}
		}
		w.Write([]byte("No result"))
	})

	http.HandleFunc("/exam", func(w http.ResponseWriter, r *http.Request) {
		tel := r.URL.Query().Get("tel")
		email := r.URL.Query().Get("email")
		// amountStr := r.URL.Query().Get("amount")

		if len(tel) == 0 && len(email) == 0 {
			w.Write([]byte("Param error!"))
			return
		}

		if r.Method == http.MethodGet {
			for _, v := range TransferMap {
				if v.Examed == true {
					continue
				}
				if len(tel) >= 11 {
					if len(v.TelHead) != 0 && len(v.TelHead) != 0 {
						if v.TelHead == tel[:3] && v.TelTail == tel[7:11] {
							if time.Now().Sub(v.Time) < time.Hour*2 {
								w.Write([]byte(fmt.Sprint(v.Amount)))
								v.Examed = true
								return
							}
						}
					}
				} else if len(email) > 5 {
					if len(v.Email) == 0 {
						continue
					}
					m, n := strings.IndexByte(v.Email, '*'), strings.LastIndexByte(v.Email, '*')
					log.Println(m, n)
					part1 := v.Email[:m]
					part2 := v.Email[n+1:]

					log.Println(part1, part2)
					//campare v.Email and email
					if part1 == email[:len(part1)] && part2 == email[len(email)-len(part2):] {
						if time.Now().Sub(v.Time) < time.Hour*2 {
							w.Write([]byte(fmt.Sprint(v.Amount)))
							v.Examed = true
							return
						}
					}
				}

			}
		}

		w.Write([]byte("No result!"))
	})

	go http.ListenAndServe(":2048", nil)
}

func Run() {
	for {
		GetTransfer(cookies)
		time.Sleep(time.Second * 2)
	}
}
