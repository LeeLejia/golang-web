# Golang后台

    该项在于实现一个较为完整的Golang后台.
    
## 目录结构
    
    ├── app                             ----------- 接口实现
    │   ├── apps.go     
    │   ├── blog.go
    │   ├── file.go
    │   ├── manifest.go
    │   ├── publish.go
    │   ├── shop.go
    │   └── users.go
    ├── bak                              ----------- 配置文件备份
    │   ├── app.toml
    │   ├── init.sql
    │   └── public.sql
    ├── common                           ----------- 共用库
    │   ├── common.go
    │   ├── conf                         ----------- 配置模块
    │   │   ├── conf.go
    │   │   └── path.go
    │   ├── handle.go                    ----------- 请求接口模块
    │   ├── log                          ----------- 日志模块
    │   │   ├── logger.go
    │   │   └── logger_test.go
    │   │   
    │   ├── RESULT_CODES.go              ----------- 返回码常量
    │   ├── session.go                   ----------- Session模块
    │   ├── sms.go                       ----------- 短信模块
    │   ├── test_utils.go
    │   └── utils.go
    ├── main.go                          ----------- 主程序
    ├── model                            ----------- 数据模型
    │   ├── t_app.go
    │   ├── t_code.go
    │   ├── test
    │   │   ├── t_app_test.go
    │   │   ├── t_code_test.go
    │   │   ├── t_file_test.go
    │   │   ├── t_log_test.go
    │   │   ├── t_user_test.go
    │   │   └── t_vlog_test.go
    │   │   
    │   ├── t_file.go
    │   ├── t_log.go
    │   ├── t_publish.go
    │   ├── t_user.go
    │   ── t_vlog.go
    └─ plugin                            ----------- 插件
       ├── email                         ----------- 邮箱插件
       │   └── email.go
       ├── plugin.go
       └── wx                            ----------- 微信机器人插件
           ├── plugins
           │   └── plugin.go
           └── robot.go
    

## 添加接口
路由配置,其中Check指定是否做用户验证
```go
routers:=[]common.BH{
		// 用户登录/注销/注册
		{Url:"/api/login",Check:false,Handle:app.Login},
		{Url:"/api/logout",Check:true,Handle:app.Logout},
		{Url:"/api/register",Check:false,Handle:app.Register},
		{Url:"/api/setUserAvatar",Check:true,Handle:app.SetUserAvatar},
		// 发布任务
		{Url:"/api/publish",Check:true,Handle:app.Publish},
		{Url:"/api/getTask",Check:false,Handle:app.GetTask},
		// 文件校验/上传图片/上传文件
		{Url:"/api/checkSha256",Check:true,Handle:app.CheckSha256},
		{Url:"/api/uploadFile",Check:true,Handle:app.UploadFile},
		{Url:"/api/listFiles",Check:true,Handle:app.ListFiles},
		{Url:"/api/deleteFile",Check:true,Handle:app.DeleteFile},
		// App添加/删除/列表获取
		{Url:"/api/developer/add-app",Check:true,Handle:app.AddApp},
		{Url:"/api/developer/list-apps",Check:true,Handle:app.ListApps},
	}
```

## 功能实现

- [x] session
- [x] 登录/注册/注销
- [x] 日志系统
- [x] 文件秒传
- [x] 文件管理
- [x] 发布任务
- [ ] 任务管理
- [ ] 个人信息修改

- [ ] 邮箱系统/注册验证/任务推送
- [ ] 微信机器人
- [ ] redis 接口数据缓存
- [ ] 支付系统
- [ ] 充值系统

