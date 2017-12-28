package app

import (
	"net/http"
	"../common"
	"../model"
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
	if icon == "" {
		icon="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADAAAAAwCAYAAABXAvmHAAADhUlEQVRoQ+1aTUwTQRj9ZkvZtkARhGi0RBHDT6IBE0JMSKhCOMihSBC5ICoxJurVxJP8yJEYL0QTD0qACxKCePBCFCXpBQ4UgyEQFVSKVNTYSmmXdnfMNE6zbFpYoGVbXG7tPmbe+973vkk6iyDO/1Ao/o39XDWAUBRb2hhbdx07JOW0TkBDv+cMI8BTQHA0tsj/Y4NhXmDgam+d/g3lFxQQII9hJCaJS0gJCM5SEUEBjX2euZitvLSqGOa76/XZ5OuAgIY+7jyDhMF4qD7lKGCmpreefR4Q0NjnaQUELfEkADC0ddfrW1UBirmmOqBY6YNngpoBZT1QM6Bs/UE9B5Q24P92oOFUwigCwD0TfrPUCUuBxmpKZXzS7987BP24nS9YXQMjeSYXF9bp7U6hY+lotqWCzSULP7CuTdoWhULxJp0WdjKFRYV/ODyJEML0WXIiFGEMznuvOcenXzhXLi7iAu6YtaM56cxBsrBjRVi+O+wrDSXgcr933d4GLTg7LewKAOCmAc5EBWyGi7iArgs61/Qyb7O7MFQeTyi7MeR10bYgm4UjRp7dOq19W5KlMd9+6bW3VLA/iFNSAVLcshsOhxSxnRYyZzNjTcWJJW2vuNkVDgwdVaxp+IN/tHfCX0Y3kSOAkJaLi6gD96vYsRQW9l0f5AIZIJ/T9HCItMRmAkh2msvZA2s8dpD/DycgMwnsHed0yRQXMQGGRHA9qtYZxRUnk6T2hLZUHGZKzMVhG908SQsZGgaZMAZX+wj37eNPnCcXFzEBZHSSnuf8eIbjwUMXNrKo6Mtv3krDTImNL/DBnz4IdmpJ0I8t8PmrPkgVZ2UzXMQEPKll7X4B3O+WhEXxovmZTBoJIw3zRr0tZ1rJPuG3EmI6+wemfNYX0/y6sZmzH800l7N5tLViUkB7pdaalao5eXPIi2kLiCv1uIadTWDAsNF8l1ZWrtAdtxAJ70OLDn918lPSQ4suTvNBwnytWAvh5rsiLURO0LxMZs7uxBnf3Tg4LsVkxJgUHbiNLPJMLG784zBpPTm4HTsgO1S7DdxKiHebm6z9VAGyyhRFkOpAFIsra2nVAVlliiJIdSCKxZW19J5yIO4v+Yhll5555hHAEVn2KQzCAJ97LuoDl/F756KbqCG39QhDV6w6QSqPEVwJ+aqBuCsCmYixlz0EYGzkYlvavSHfVlG4xbe0/V/rXqJPNUQmQgAAAABJRU5ErkJggg=="
	}
	if file==""{
		common.ReturnEFormat(w, common.CODE_VERIFY_FAIL, "请先上传App")
		return
	}
	// todo 检测文件有效性
	if src!=""{
		// todo 检测文件有效性
	}
	app:=&model.T_app{
		Icon:icon,
		AppId:common.GetRandomString(16),
		Name:name,
		Version:version,
		Describe:describe,
		Developer:user.Id,
		Valid:valid,
		File:file,
		Src:src,
		Expend:simplejson.New(), // TODO
		DownloadCount:0,
		CreatedAt:time.Now(),
	}
	err:=model.AppModel.Insert(app)
	if err!=nil{
		common.ReturnEFormat(w,common.CODE_DB_RW_ERR, "数据库插入错误")
		return
	}
	common.ReturnFormat(w,200,nil,"操作成功")
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
	name := r.PostFormValue("name")
	role := r.PostFormValue("role")

	if user.Role==model.USER_ROLE_DEVELOPER{
		cond=fmt.Sprintf("where developer=%d",user.Id)
	}

	apps,err:= model.AppModel.Query(cond,"","")
	if err!=nil{
		fmt.Println(err.Error())
	}
	common.ReturnFormat(w, 200, apps, "success")
}

/**
删除App
 */
func DeleteApps(w http.ResponseWriter, r *http.Request, user *model.T_user){

}