/**
	配置文件
 */
package conf

import (
	"github.com/BurntSushi/toml"
	"fmt"
	"os"
)

type tomlFile struct {
	DBHost     string `toml:"DBHost"`
	DBPort     string `toml:"DBPort"`
	DBUser     string `toml:"DBUser"`
	DBPassword string `toml:"DBPassword"`
	DBName     string `toml:"DBName"`
	DBNameTest string `toml:"DBNameTest"`
	PathFile   string `toml:"Path_Upload_file"`
	PathPic	   string `toml:"Path_Upload_pic"`
	ServerHost string `toml:"serverHost"`
	ServerPort string `toml:"serverPort"`
	StaticPath string `toml:"staticPath"`
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
	fmt.Println(fmt.Sprintf("\x1b[%dm配置:%s\x1b[0m",uint8(93),App))
	if _,err:=os.Stat(App.StaticPath);err!=nil{
		os.MkdirAll(App.StaticPath,0755)
		fmt.Println(fmt.Sprintf("\x1b[%dm路径不存在%s,即将重新创建。detail:%s\x1b[0m",uint8(91),App.StaticPath,err.Error()))
	}
	if _,err:=os.Stat(App.PathPic);err!=nil{
		os.MkdirAll(App.PathPic,0755)
		fmt.Println(fmt.Sprintf("\x1b[%dm路径不存在%s,即将重新创建。detail:%s\x1b[0m",uint8(91),App.PathPic,err.Error()))
	}
	if _,err:=os.Stat(App.PathFile);err!=nil{
		os.MkdirAll(App.PathFile,0755)
		fmt.Println(fmt.Sprintf("\x1b[%dm路径不存在%s,即将重新创建。detail:%s\x1b[0m",uint8(91),App.PathFile,err.Error()))
	}
}
