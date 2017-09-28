/**

	前端请求测试程序
 */
package main
import (
	"net/http"
	"fmt"
	"mypack/utils"
	"time"
)

func main(){
	fmt.Println("正在监听端口：80")
	path:=utils.GetProDir()+"/log.txt"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		fmt.Println("========>")
		fmt.Println(fmt.Sprintf("[remoteAddr]   %s",r.RemoteAddr))
		fmt.Println(fmt.Sprintf("[url]   %s",r.URL.Path))
		fmt.Println(fmt.Sprintf("[header]"))
		for _,h:=range(r.Header){
			fmt.Println(fmt.Sprintf("   --%s",fmt.Sprint(h)))
		}
		r.ParseForm()
		fmt.Println(fmt.Sprintf("[form]"))
		fmt.Println(r.Form)
		utils.AppendToFile(path,[]byte(time.Now().String()+"  |  "+ fmt.Sprint(r.Form)+"\n"))
		r.ParseMultipartForm(32 << 20)
		fmt.Println(fmt.Sprintf("[postForm]"))
		fmt.Println(r.PostForm)
		w.Write([]byte("{\"code\":\"ok\";}"))
	})
	http.ListenAndServe(":80",nil)
}
