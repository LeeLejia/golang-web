package plugin

import (
	"./wx"
	"./email"
)

func Init(){
	// 微信插件
	wx.Init()
	// 邮箱插件
	email.Init()
}