package common


const(
	// 请求参数错误
	CODE_PARAMS_INVALID=500
	// 系统错误
	CODE_SERVICE_ERR=503
	// 校验错误
	CODE_VERIFY_FAIL=404
	// SESSION登录超时
	CODE_NOT_ALLOW = 406
	//需要重新登录
	CODE_NEET_LOGIN_AGAIN = 1024
	//角色不符合
	CODE_ROLE_INVADE = 512
	//数据库读写错误
	CODE_DB_RW_ERR = 128
)
