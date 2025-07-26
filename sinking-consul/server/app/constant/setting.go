package constant

const (
	WebGroup = "web"              //网站配置组
	WebName  = WebGroup + ".name" //网站名称

	LoginGroup    = "login"                  //登录配置组
	LoginAccount  = LoginGroup + ".account"  //登录账号
	LoginPassword = LoginGroup + ".password" //登录密码
	LoginToken    = LoginGroup + ".token"    //登录token
	LoginExpire   = LoginGroup + ".expire"   //登录token过期时间
)
