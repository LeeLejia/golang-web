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
	checkDirs()
}

func checkDirs(){
	for _,path:=range []string{App.StaticPath,App.PathPic,App.PathFile} {
		rp:=RealFilePath(path)
		if _,err:=os.Stat(rp);err!=nil{
			fmt.Println(fmt.Sprintf("\x1b[%dm路径不存在%s,即将重新创建。detail:%s\x1b[0m",uint8(91),rp,err.Error()))
			err:=os.MkdirAll(rp,0755)
			if err!=nil{
				fmt.Println(fmt.Sprintf("\x1b[%dm创建路径失败:%s。detail:%s\x1b[0m",uint8(91),rp,err.Error()))
			}
		}
	}
}
