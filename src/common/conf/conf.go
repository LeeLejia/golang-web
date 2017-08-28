/**
	配置文件
 */
package conf

import (
	"github.com/BurntSushi/toml"
	"fmt"
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
		fmt.Println(fmt.Sprintf("\x1b[%dm读取配置文件失败！detail:%s\x1b[0m",91,err.Error()))
		panic(err)
	}
	fmt.Println(App)
}
