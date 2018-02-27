/**
 * 一些读取配置用到的公共函数
 */
package conf

import (
	"os"
	"path/filepath"
	"strings"
)

// 项目根目录
var AppRootPath string

// AppPath 项目根目录
func AppPath(appPath ...string) string {
	if len(appPath) > 0 {
		AppRootPath = appPath[0]
	}
	if AppRootPath == "" {
		AppRootPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		if !fileExists(filepath.Join(AppRootPath, "conf", "task.toml")) {
			workPath, _ := os.Getwd()
			workPath, _ = filepath.Abs(workPath)
			AppRootPath = workPath
		}
	}
	return AppRootPath
}

// RealFilePath 返回绝对路径
func RealFilePath(relFilename string) string {
	if strings.HasPrefix(relFilename, "/") || relFilename[1]==':'{
		return relFilename
	}
	return filepath.Join(AppPath(), relFilename)
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
