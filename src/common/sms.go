/**
 * 短信发送模块
 */

package common

import (
	"fmt"
	"math/rand"
	"time"
	"net/url"
	"crypto/md5"
	"strconv"
	"io"
	"strings"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

var SMS map[string]string

func SendSMSCode(phone string) {
	if SMS == nil {
		SMS = map[string]string{}
	}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	SMS[phone] = vcode

	crutime := time.Now()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(crutime.Unix(), 10))
	nonce := fmt.Sprintf("%x", h.Sum(nil))
	loc, err := time.LoadLocation("GMT")
	if err != nil {
		fmt.Println(err.Error())
	}
	stime := time.Date(crutime.Year(),crutime.Month(),crutime.Day(),crutime.Hour(),crutime.Minute(),crutime.Second(),crutime.Nanosecond(),loc)
	fmt.Println(stime.Location())
	accessKeyId := "LTAI1NsQHLR9jncU"
	accessKeySecret := "2lLDdD2MQdbFGCy3Uknr9QX5i2Nk5w"

	params := map[string]interface{}{}
	params["SignatureMethod"] = "HMAC-SHA1"
	params["SignatureNonce"] = nonce
	params["AccessKeyId"] = accessKeyId
	params["Timestamp"] = stime.Local().Format("2006-01-02T15:04:05Z")
	params["SignatureVersion"] = "1.0"
	params["Format"] = "xml"

	params["Action"]= "SendSms"
	params["Version"]= "2017-05-25"
	params["RegionId"]= "cn-hangzhou"
	params["PhoneNumbers"]= phone
	params["SignName"]= "阿里云短信测试专用"
	params["TemplateParam"]= "{\"customer\":\"test\"}"
	params["TemplateCode"]= "SMS_85630051"
	params["OutId"]= "123"

	i := 0
	res := "http://127.0.0.1:8888?"
	for k,v:=range params{
		if i == 0 {
			res += fmt.Sprintf("%s=%v",k,v)
		}else{
			res += fmt.Sprintf("&%s=%v",k,v)
		}
		i++
	}
	a,err:=url.Parse(res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result :=a.Query().Encode()
	strings.Replace(result,"+", "%20",-1)
	strings.Replace(result,"*", "%2A",-1)
	strings.Replace(result,"-", "%7E",-1)
	key := []byte(accessKeySecret+"&")
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(result))
	input := fmt.Sprintf("%x", mac.Sum(nil))
	final := base64.StdEncoding.EncodeToString([]byte(input))
	params["Signature"] = final

	res="http://127.0.0.1:8888?"
	for k,v:=range params{
		if i == 0 {
			res += fmt.Sprintf("%s=%v",k,v)
		}else{
			res += fmt.Sprintf("&%s=%v",k,v)
		}
		i++
	}
	a,err=url.Parse(res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	allResult :=a.Query().Encode()
	shit := "http://dysmsapi.aliyuncs.com/?"+allResult
	fmt.Println(shit)
}
