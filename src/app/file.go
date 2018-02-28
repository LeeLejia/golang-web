package app

import (
	"../model"
	"net/http"
	"../common"
	"../common/conf"
	"../common/log"
	"os"
	"io"
	fm "github.com/cjwddz/fast-model"
	"time"
	"strings"
	"path"
	"fmt"
)

/**
获取文件
只有管理员和超级管理员可以获取文件列表
其他成员只能下载单个文件
 */
func ListFiles(ssesion * common.Session,w http.ResponseWriter, r *http.Request){
	if ssesion.User.Role!=model.USER_ROLE_SUPER && ssesion.User.Role!=model.USER_ROLE_ADMIN{
		common.ReturnEFormat(w, common.CODE_ROLE_INVADE,"你没有获取文件列表权限！")
		return
	}
}

/**
检查sha256,判断文件是否存在
 */
func CheckSha256(_ *common.Session, w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	// 获取sha256
	sha256:=r.Form.Get("sha256")
	if sha256==""{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"参数缺少sha256！")
		return
	}else if len(sha256)!=64{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"sha256不合法！")
		return
	}
	// 判断文件是否存在
	rs,err:=FileModel.Query(fm.DbCondition{}.And("=","key",sha256).Limit(1,-1))
	if err!=nil || rs==nil || len(rs)==0{
		common.ReturnFormat(w,common.CODE_SUCCESS,map[string]interface{}{"exist":false,"msg":"文件不存在！"})
		return
	}
	common.ReturnFormat(w,common.CODE_SUCCESS,map[string]interface{}{"exist":true,"sha256":sha256})
	return
}

/**
	上传小文件
 */
func UploadFile(sess *common.Session, w http.ResponseWriter, r *http.Request) {
	f, h, err := r.FormFile("file")
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
		log.E("UploadFile出错",sess.User.Account,err.Error())
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "内部服务出错！")
		log.E("UploadFile出错",sess.User.Account,err.Error())
		return
	}
	file:=model.T_File{
		Key:fileKey,
		Type:fileType,
		Name:h.Filename,
		Owner:sess.User.Account,
		CreatedAt:time.Now(),
	}
	err=FileModel.Insert(file)
	if err!=nil{
		log.E("UploadFile出错",sess.User.Account,err.Error())
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "数据库写入失败！")
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"key": fileKey,"msg":"success"})
	log.N("UploadFile上传文件成功",sess.User.Account,h.Filename)
}
/**
上传图片
 */
func UploadPicture(sess *common.Session, w http.ResponseWriter, r *http.Request){
	f, h, err := r.FormFile("img")
	if err != nil {
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"请提交文件")
		return
	}
	defer f.Close()
	filename := h.Filename
	fileSuffix := strings.ToLower(path.Ext(filename))
	if fileSuffix != ".jpg" && fileSuffix != ".png" && fileSuffix!=".gif" && fileSuffix!=".jpeg"{
		common.ReturnEFormat(w, 500, fmt.Sprintf("不支持的图片格式'%s'", fileSuffix))
		return
	}
	fileKey := r.FormValue("sha256")
	if len(fileKey)!=64{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"文件指纹缺失！")
		return
	}
	filePath := conf.App.PathFile +"/"+ fileKey
	t, err := os.Create(filePath)
	if err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "内部服务出错！")
		log.E("UploadPicture出错",sess.User.Account,err.Error())
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "内部服务出错！")
		log.E("UploadPicture出错",sess.User.Account,err.Error())
		return
	}
	file:=model.T_File{
		Key:fileKey,
		Type:"picture-normal",
		Name:h.Filename,
		Owner:sess.User.Account,
		CreatedAt:time.Now(),
	}
	err=FileModel.Insert(file)
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "数据库写入失败！")
		log.E("UploadPicture出错",sess.User.Account,err.Error())
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"key": fileKey,"msg":"success"})
	log.N("UploadPicture上传图片成功",sess.User.Account,h.Filename)
}

