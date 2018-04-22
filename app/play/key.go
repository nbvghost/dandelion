package play

const (
	SessionAdmin    string = "SESSION_ADMIN"   //商家后台
	SessionManager  string = "SESSION_MANAGER" //系统管理
	SessionUser     string = "SESSION_USER"    //前台用户
	SessionShop     string = "SESSION_SHOP"    //商铺
	SessionCompany  string = "SESSION_COMPANY"
	SessionAction   string = "SESSION_ACTION"
	SessionUserID   string = "SESSION_USERID"
	SessionOpenID   string = "SESSION_OPENID"
	SessionRedirect string = "SESSION_REDIRECT"
	SessionCart     string = "SESSION_CART"
	SessionCaptcha  string = "SESSION_CAPTCHA"

	ActionKey_add    string = "add"
	ActionKey_save   string = "save"
	ActionKey_change string = "change"
	ActionKey_get    string = "get"
	ActionKey_one    string = "one"
	ActionKey_list   string = "list"
	ActionKey_del    string = "del"
	Paging           int    = 10

	ConfigurationKey_component_verify_ticket uint64 = 1001
)
