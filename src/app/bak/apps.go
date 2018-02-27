package bak

import (
	"net/http"
	"../../common"
	fm "github.com/cjwddz/fast-model"
	"../../model"
	"fmt"
	"time"
	"github.com/bitly/go-simplejson"
)

/**
添加App
 */
func AddApp(w http.ResponseWriter, r *http.Request, user *model.T_user) {
	icon := r.PostFormValue("icon")
	version := r.PostFormValue("version")
	describe := r.PostFormValue("describe")
	name:=r.PostFormValue("name")
	file:=r.PostFormValue("file")
	src:=r.PostFormValue("src")
	_valid :=r.PostFormValue("valid")
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
		Developer:user.Account,
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
func ListApps(w http.ResponseWriter, r *http.Request, user *model.T_user){
	if user.Role!=model.USER_ROLE_DEVELOPER && user.Role!=model.USER_ROLE_SUPER && user.Role!=model.USER_ROLE_ADMIN{
		common.ReturnEFormat(w, common.CODE_ROLE_INVADE,"你没有获取APP列表的权限！")
		return
	}
	condstr:= make([]string,0)
	c:=1
	valid := r.PostFormValue("valid")
	if valid!=""{
		condstr=append(condstr, fmt.Sprintf("valid = $%d",c))
		c++

	}
	cdt:=fm.DbCondition{}
	if user.Role==model.USER_ROLE_DEVELOPER{
		cdt=cdt.And("=","developer",user.Id)
	}
	cdt=cdt.And2(r,"=","b_valid").And2(r,"like","s_name")
	cdt=cdt.Order(r.PostFormValue("order")).Limit2(r,"start","len")
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