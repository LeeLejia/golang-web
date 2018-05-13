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
	"fmt"
	"strconv"
)
/**
获取文件列表
 */
func ListFiles(sess * common.Session,w http.ResponseWriter, r *http.Request){
	cond:=fm.DbCondition{}.And2(r,"like","s_type").And2(r,"like","s_name").And("=","owner",sess.User.Email)
	total,err:=FileModel.Count(cond)
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("ListFiles出错",sess.User.Email,err.Error())
		return
	}
	result,err:=FileModel.Query(cond.Limit2(r,"start","count").Order("order by created_at desc"))
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("ListFiles出错",sess.User.Email,err.Error())
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"files": result,"total":total})
	log.N("ListFiles",sess.User.Email,fmt.Sprintf("listCount=%d,total=%d",len(result),total))
	return
}
/**
删除文件
 */
func DeleteFile(sess * common.Session,w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	hash:=r.Form.Get("hash")
	if hash == "" {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "需要提供文件hash")
		return
	}
	idStr:=r.Form.Get("id")
	if idStr == "" {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "需要文件id")
		return
	}
	id,err:=strconv.Atoi(idStr)
	if err!= nil {
		common.ReturnEFormat(w, common.CODE_PARAMS_INVALID, "非法id")
		log.W("DeleteFile",sess.User.Email,err.Error())
		return
	}
	cond:=fm.DbCondition{}.And("=","owner",sess.User.Email).And("=","id",id).And("=","key",hash)
	err = FileModel.Delete(cond)
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("DeleteFile",sess.User.Email,err.Error())
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"msg": "删除文件成功","id":id,"hash":hash})
	log.N("DeleteFile",sess.User.Email,fmt.Sprintf("id=%d,hash=%s",id,hash))
	return
}
/**
检查sha256,判断文件是否存在
 */
func CheckSha256(sess *common.Session, w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	// 获取sha256
	sha256:=r.Form.Get("hash")
	if sha256==""{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"请提供hash值！")
		return
	}else if len(sha256)!=64{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"hash不合法！")
		return
	}
	// 判断文件是否存在
	rs,err:=FileModel.Query(fm.DbCondition{}.And("=","key",sha256).Limit(1,-1))
	if err!=nil || rs==nil || len(rs)==0{
		common.ReturnFormat(w,common.CODE_SUCCESS,map[string]interface{}{"exist":false,"msg":"文件不存在！"})
		return
	}
	newfile := rs[0].(model.T_File)
	newfile.Owner = sess.User.Email
	newfile.CreatedAt = time.Now()
	err=FileModel.Insert(newfile)
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("UploadFile出错",sess.User.Email,err.Error())
		return
	}
	log.N("CheckSha256插入新数据成功",sess.User.Email,newfile.Name)
	common.ReturnFormat(w,common.CODE_SUCCESS,map[string]interface{}{"exist":true,"key":sha256})
	return
}

/**
	上传文件
 */
func UploadFile(sess *common.Session, w http.ResponseWriter, r *http.Request) {
	f, h, err := r.FormFile("file")
	if err != nil {
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"请提交文件")
		return
	}
	defer f.Close()
	fileKey := r.FormValue("hash")
	if len(fileKey)!=64{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"文件指纹缺失！")
		return
	}
	if h.Size> 512*1024*1024{
		common.ReturnEFormat(w,common.CODE_PARAMS_INVALID,"上传的文件不能超过512M")
		return
	}

	filePath := conf.App.PathFile +"/"+ fileKey
	t, err := os.Create(filePath)
	if err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "服务器内部出错！")
		log.E("UploadFile出错",sess.User.Email,err.Error())
		return
	}
	defer t.Close()
	if _, err := io.Copy(t, f); err != nil {
		common.ReturnEFormat(w, common.CODE_SERVICE_ERR, "服务器内部出错！")
		log.E("UploadFile出错",sess.User.Email,err.Error())
		return
	}
	file:=model.T_File{
		Key:fileKey,
		Type:h.Header.Get("Content-Type"),
		Size:h.Size,
		Name:h.Filename,
		Owner:sess.User.Email,
		CreatedAt:time.Now(),
	}
	err=FileModel.Insert(file)
	if err!=nil{
		common.ReturnEFormat(w, common.CODE_DB_RW_ERR, "服务器内部出错！")
		log.E("UploadFile出错",sess.User.Email,err.Error())
		return
	}
	common.ReturnFormat(w, common.CODE_OK, map[string]interface{}{"key": fileKey,"msg":"success"})
	log.N("UploadFile上传文件成功",sess.User.Email,h.Filename)
	return
}
