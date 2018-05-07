package app

import "os"

import (
	"../common"
	"../common/log"
	"net/http"
	"io"
	"../common/conf"
	"../model"
	"time"
)

/**
	上传小文件
 */
func UploadMarkdown(sess *common.Session, w http.ResponseWriter, r *http.Request) {
	// todo
	f, h, err := r.FormFile("markdown")
	if err != nil {
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"请提交文件")
		return
	}
	defer f.Close()
	fileKey := r.FormValue("sha256")
	if len(fileKey)!=64{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"文件指纹缺失！")
		return
	}
	fileType:=r.FormValue("type")
	if fileType==""{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"未指定文件类型！")
		return
	}
	filePath := conf.App.PathFile +"/"+ fileKey
	t, err := os.Create(filePath)
	if err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "内部服务出错！")
		log.E("UploadFile出错",sess.User.Email,err.Error())
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "内部服务出错！")
		log.E("UploadFile出错",sess.User.Email,err.Error())
		return
	}
	file:=model.T_File{
		Key:fileKey,
		Type:fileType,
		Name:h.Filename,
		Owner:sess.User.Email,
		CreatedAt:time.Now(),
	}
	err=FileModel.Insert(file)
	if err!=nil{
		log.E("UploadFile出错",sess.User.Email,err.Error())
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "数据库写入失败！")
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"key": fileKey,"msg":"success"})
	log.N("UploadFile上传文件成功",sess.User.Email,h.Filename)
}
