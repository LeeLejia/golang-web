package app

import (
	"../model"
	"net/http"
	"../common"
	"fmt"
	"strings"
	"../common/conf"
	"path"
	"os"
	"io"
	"time"
)

/**
获取文件
只有管理员和超级管理员可以获取文件列表
其他成员只能下载单个文件
todo
 */
func ListFiles(w http.ResponseWriter, r *http.Request, user *model.T_user){
	if user.Role!=model.USER_ROLE_SUPER && user.Role!=model.USER_ROLE_ADMIN{
		common.ReturnEFormat(w, common.CODE_ROLE_INVADE,"你没有获取文件列表权限！")
		return
	}
	cond:=""
	apps,err:= model.FindFiles(cond,"","")
	if err!=nil{
		fmt.Println(err.Error())
		common.ReturnFormat(w, common.CODE_SERVICE_ERR, apps, err.Error())
		return
	}
	common.ReturnFormat(w, 200, apps, "success")
}

/**
上传文件
 */
func UploadFile(w http.ResponseWriter, r *http.Request, user *model.T_user){
	f, h, err := r.FormFile("file")
	defer f.Close()
	fileKey := common.GetRandomString(16)
	filePath := conf.App.PathFile +"/"+ fileKey
	t, err := os.Create(filePath)
	if err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "内部服务出错！")
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "内部服务出错！")
		return
	}
	file:=model.T_File{
		FileKey:fileKey,
		FileName:h.Filename,
		FileType:model.FILE_TYPE_FILE,
		Owner:user.Id,
		Path:filePath,
		CreatedAt:time.Now(),
	}
	err=file.Insert()
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, err.Error())
		return
	}
	common.ReturnFormat(w, 200, map[string]interface{}{"key": fileKey}, "success")
}
/**
上传图片
 */
func UploadPicture(w http.ResponseWriter, r *http.Request, user *model.T_user) {
	f, h, err := r.FormFile("file")
	if err != nil {
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"请提交图片")
		return
	}
	filename := h.Filename
	fileSuffix := strings.ToLower(path.Ext(filename))
	if fileSuffix != ".jpg" && fileSuffix != ".png" && fileSuffix!=".gif"{
		common.ReturnEFormat(w, 500, fmt.Sprintf("不支持的图片格式'%s'", fileSuffix))
		return
	}
	defer f.Close()
	fileKey := common.GetRandomString(16)
	filePath := conf.App.PathPic +"/"+ fileKey
	t, err := os.Create(filePath)
	if err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "内部服务出错！")
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "内部服务出错！")
		return
	}
	file:=model.T_File{
		FileKey:fileKey,
		FileType:model.FILE_TYPE_PIC,
		FileName:h.Filename,
		Owner:user.Id,
		Path:filePath,
		CreatedAt:time.Now(),
	}
	err=file.Insert()
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, err.Error())
		return
	}
	common.ReturnFormat(w, 200, map[string]interface{}{"key": fileKey}, "success")
}

