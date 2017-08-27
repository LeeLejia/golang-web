/**
 * 公共配置文件读取
 *
 * @author zhangqichao
 * Created on 2017-08-18
 */
package conf

import (
	"github.com/BurntSushi/toml"
	"fmt"
	"mypack/log"
	"github.com/labstack/echo/middleware"
)

type tomlFile struct {
	DBHost         string `toml:"DBHost"`
	DBPort         string `toml:"DBPort"`
	DBUser         string `toml:"DBUser"`
	DBPassword     string `toml:"DBPassword"`
	DBName         string `toml:"DBName"`
	DBNameTest     string `toml:"DBNameTest"`
	Static         string `toml:"Static"`
	FileUploadPath string `toml:"fileUploadPath"`
	ServerHost     string `toml:"serverHost"`
	ServerPort     string `toml:"serverPort"`
}

var App *tomlFile

func init() {
	App = new(tomlFile)
}

func Init(filePath string) {
	var path=RealFilePath(filePath)
	_, err := toml.DecodeFile(path, App)
	if err != nil {
		log.E("读取配置文件出错！路径：%s",path)
		panic(err)
	}
	fmt.Println(App)
}
