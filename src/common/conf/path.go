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
var _appPath string

// AppPath 项目根目录
func AppPath(appPath ...string) string {
	if len(appPath) > 0 {
		_appPath = appPath[0]
	}

	if _appPath == "" {
		_appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		if !fileExists(filepath.Join(_appPath, "conf", "app.toml")) {
			workPath, _ := os.Getwd()
			workPath, _ = filepath.Abs(workPath)
			_appPath = workPath
		}
	}
	return _appPath
}

// RealFilePath 返回绝对路径
// =>
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
