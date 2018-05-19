package app

import (
	"net/http"
	"../common"
	"../common/log"
	"fmt"
)



/**
添加App
 */
func Order(sess *common.Session, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//icon := r.FormValue("icon")
	log.N("GetTask","",fmt.Sprintf("abc%d",0))
}