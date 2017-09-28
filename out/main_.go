/**

	前端请求测试程序
 */
package main
import (
	"net/http"
	"fmt"
	"mypack/utils"
	"github.com/bitly/go-simplejson"
	"regexp"
	"strings"
)

func main(){
	fmt.Println("正在监听端口：88")
	path:=utils.GetProDir()+"/微信采集.txt"
	fmt.Println("文件目录："+path)
	reg, _ := regexp.Compile("(13|15|18)[0-9]{9}")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		r.ParseMultipartForm(32 << 20)
		data:=r.PostForm.Get("sns")
		js,err:=simplejson.NewJson([]byte(data))
		if err==nil{
			content,_:=js.Get("content").String()
			if !reg.MatchString(content){
				w.Write([]byte("{\"code\":\"no\";}"))
				return
			}
			fmt.Println(content)
			content=strings.Replace(strings.Replace(strings.Replace(content," ","",-1),"\r\n","",-1),"\n","",-1)
			utils.AppendToFile(path,[]byte(content+"\n"))
		}
		w.Write([]byte("{\"code\":\"ok\";}"))
	})
	http.ListenAndServe(":88",nil)
}
