/**
	配置文件
 */
package conf

import (
	"github.com/BurntSushi/toml"
	"fmt"
)

type tomlFile struct {
	DBHost     string `toml:"DBHost"`
	DBPort     string `toml:"DBPort"`
	DBUser     string `toml:"DBUser"`
	DBPassword string `toml:"DBPassword"`
	DBName     string `toml:"DBName"`

	TcpAddress string `toml:"Address"`
}

var App *tomlFile

func init() {
	App = new(tomlFile)
}

func Init(filePath string) {
	var path=RealFilePath(filePath)
	fmt.Println(fmt.Sprintf("\x1b[%dm加载配置文件:%s\x1b[0m",uint8(93),path))
	_, err := toml.DecodeFile(path, App)
	if err != nil {
		fmt.Println(fmt.Sprintf("\x1b[%dm文件不存在，请检查指定路径是否存在配置文件。！detail:%s\x1b[0m",uint8(91),err.Error()))
		panic(err)
	}
}
