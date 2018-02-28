package app

import (
	"net/http"
	"../common"
	fm "github.com/cjwddz/fast-model"
	"../model"
	"fmt"
	"time"
	"github.com/bitly/go-simplejson"
)

/**
添加App
 */
func AddApp(sess *common.Session, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	icon := r.FormValue("icon")
	version := r.FormValue("version")
	describe := r.FormValue("describe")
	name:=r.FormValue("name")
	file:=r.FormValue("file")
	src:=r.FormValue("src")
	_valid :=r.FormValue("valid")
	valid:=true
	if _valid =="" || _valid =="false"{
		valid=false
	}
	if name == ""{
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "请输入App名称")
		return
	}
	if version==""{
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "版本号不能为空")
		return
	}
	if describe==""{
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "App描述不能为空")
		return
	}
	if file==""{
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "请先上传App")
		return
	}
	if src!=""{
		// todo 检测文件有效性
	}
	app:=&model.T_app{
		Icon:icon,
		AppId:common.GetRandomString(16),
		Name:name,
		Version:version,
		Describe:describe,
		Developer:sess.User.Account,
		Valid:valid,
		File:file,
		Src:src,
		Expend:simplejson.New(), // TODO
		DownloadCount:0,
		CreatedAt:time.Now(),
	}
	err:=AppModel.Insert(app)
	if err!=nil{
		common.ReturnEFormat(w,common.CODE_DB_RW_ERR, "数据库插入错误")
		return
	}
	common.ReturnFormat(w,common.CODE_OK,map[string]interface{}{"msg":"操作成功！"})
}

/**
获取App
 */
func ListApps(sess *common.Session,w http.ResponseWriter, r *http.Request){
	if sess.User.Role!=model.USER_ROLE_DEVELOPER && sess.User.Role!=model.USER_ROLE_SUPER && sess.User.Role!=model.USER_ROLE_ADMIN{
		common.ReturnEFormat(w, common.CODE_ROLE_INVADE,"你没有获取APP列表的权限！")
		return
	}
	r.ParseForm()
	cdt:=fm.DbCondition{}
	if sess.User.Role==model.USER_ROLE_DEVELOPER{
		cdt=cdt.And("=","developer",sess.User.Id)
	}
	cdt=cdt.And2(r,"=","b_valid").And2(r,"like","s_name")
	cdt=cdt.Order(r.FormValue("order")).Limit2(r,"start","len")
	apps,err:= AppModel.Query(cdt)
	if err!=nil{
		fmt.Println(err.Error())
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"apps":apps,"msg":"操作成功！"})
}

/**
删除App
 */
func DeleteApps(w http.ResponseWriter, r *http.Request, user *model.T_user){

}