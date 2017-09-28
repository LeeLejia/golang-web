package app

import (
	"../model"
	"net/http"
	"../common"
	"fmt"
	"io/ioutil"
	"../common/conf"
	"strings"
	"path"
	"os"
	"io"
)

/**
上传文件
 */
func UploadFile(w http.ResponseWriter, r *http.Request, user *model.T_user){
	file,fheader,err:=r.FormFile("file")
	if err!=nil{
		fmt.Println(err.Error())
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"上传文件失败")
		return
	}
	fmt.Println(fheader)
	bs,err:=ioutil.ReadAll(file)
	ioutil.WriteFile(conf.App.FileUploadPath+"/"+fheader.Filename,bs,0777)
}

/**
获取文件
只有管理员和超级管理员可以获取文件列表
其他成员只能下载单个文件
todo
 */
func FindFiles(w http.ResponseWriter, r *http.Request, user *model.T_user){
	if user.Role!=model.USER_ROLE_SUPER && user.Role!=model.USER_ROLE_ADMIN{
		common.ReturnEFormat(w, common.CODE_ROLE_INVADE,"你没有获取文件列表！")
		return
	}
	cond:=""
	apps,err:= model.FindFiles(cond,"","")
	if err!=nil{
		fmt.Println(err.Error())
	}
	common.ReturnFormat(w, 200, apps, "success")
}


func UploadPicture(w http.ResponseWriter, r *http.Request) {
	f, h, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	filename := h.Filename
	fileSuffix := strings.ToLower(path.Ext(filename))
	if fileSuffix != ".jpg" && fileSuffix != ".png" {
		common.ReturnEFormat(w, 500, fmt.Sprintf("不支持的图片格式'%s'", fileSuffix))
		return
	}
	defer f.Close()
	filePath := conf.App.FileUploadPath + filename
	t, err := os.Create("." + filePath)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	common.ReturnFormat(w, 200, map[string]interface{}{"filePath": conf.App.ServerHost + filePath}, "SUCCESS")
}

