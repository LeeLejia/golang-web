package app

import (
	"fmt"
	"../src/common"
	"../src/common/conf"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

func UploadPicture(w http.ResponseWriter, r *http.Request) {
	common.SetContent(w)
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
