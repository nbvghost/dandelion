package result

import (
	"time"

	"github.com/nbvghost/dandelion/constrain"
)

type ActionResultCode int

type ActionResultCodeMessage struct {
	Code    ActionResultCode
	Message string
}

const Success ActionResultCode = 0
const NotFound ActionResultCode = 10
const IsExist ActionResultCode = 20
const HttpStatus ActionResultCode = 30
const Warning ActionResultCode = 40
const Error ActionResultCode = -100
const Fail ActionResultCode = -101
const TypeError ActionResultCode = -102
const SQLError ActionResultCode = -103
const AuthError ActionResultCode = -104
const HttpError ActionResultCode = -105
const JSONError ActionResultCode = -106
const IOError ActionResultCode = -107
const NIOError ActionResultCode = -108
const SessionTimeOut ActionResultCode = -400

var (
	TokenNotFound    ActionResultCodeMessage = ActionResultCodeMessage{Code: 40001, Message: "Token不存在"}
	TokenInvalid     ActionResultCodeMessage = ActionResultCodeMessage{Code: 40002, Message: "Token格式错误或无效"}
	TokenExpired     ActionResultCodeMessage = ActionResultCodeMessage{Code: 40003, Message: "Token已过期，请重新登录"}
	NotLogin         ActionResultCodeMessage = ActionResultCodeMessage{Code: 40004, Message: "未登录，请先登录"}
	AccountPwdErr    ActionResultCodeMessage = ActionResultCodeMessage{Code: 40005, Message: "账号或密码错误"}
	CaptchaErr       ActionResultCodeMessage = ActionResultCodeMessage{Code: 40006, Message: "验证码错误或已过期"}
	AccountForbidden ActionResultCodeMessage = ActionResultCodeMessage{Code: 40007, Message: "账号已被禁用"}
	NoPermission     ActionResultCodeMessage = ActionResultCodeMessage{Code: 40008, Message: "当前角色无操作权限"}
	DataForbidden    ActionResultCodeMessage = ActionResultCodeMessage{Code: 40009, Message: "越权访问，无权操作该数据"}

	// ---------- 参数校验 41xxx ----------
	ParamErr     ActionResultCodeMessage = ActionResultCodeMessage{Code: 41000, Message: "请求参数错误"}
	ParamMissing ActionResultCodeMessage = ActionResultCodeMessage{Code: 41001, Message: "缺少必填参数"}
	ParamFormat  ActionResultCodeMessage = ActionResultCodeMessage{Code: 41002, Message: "参数格式非法"}
	ParamRange   ActionResultCodeMessage = ActionResultCodeMessage{Code: 41003, Message: "参数值超出合法范围"}
	RepeatSubmit ActionResultCodeMessage = ActionResultCodeMessage{Code: 41004, Message: "请勿重复提交"}

	// ---------- 业务逻辑 42xxx ----------
	DataNotFound   ActionResultCodeMessage = ActionResultCodeMessage{Code: 42001, Message: "数据不存在"}
	DataExist      ActionResultCodeMessage = ActionResultCodeMessage{Code: 42002, Message: "数据已存在"}
	StatusNotAllow ActionResultCodeMessage = ActionResultCodeMessage{Code: 42003, Message: "当前数据状态不支持该操作"}
	StockNotEnough ActionResultCodeMessage = ActionResultCodeMessage{Code: 42004, Message: "库存不足"}
	BillFinished   ActionResultCodeMessage = ActionResultCodeMessage{Code: 42005, Message: "单据已完结，不可修改"}
	FileFail       ActionResultCodeMessage = ActionResultCodeMessage{Code: 42006, Message: "文件不存在或上传失败"}
	FileLimit      ActionResultCodeMessage = ActionResultCodeMessage{Code: 42007, Message: "文件大小或后缀不合法"}

	// ---------- 服务异常 50xxx ----------
	ServerErr      ActionResultCodeMessage = ActionResultCodeMessage{Code: 50000, Message: "服务器内部异常"}
	DBOperateErr   ActionResultCodeMessage = ActionResultCodeMessage{Code: 50001, Message: "数据库操作失败"}
	ThirdPartyErr  ActionResultCodeMessage = ActionResultCodeMessage{Code: 50002, Message: "调用第三方服务失败"}
	TimeoutErr     ActionResultCodeMessage = ActionResultCodeMessage{Code: 50003, Message: "请求处理超时"}
	RequestLimit   ActionResultCodeMessage = ActionResultCodeMessage{Code: 50004, Message: "请求过于频繁，请稍后重试"}
	SystemMaintain ActionResultCodeMessage = ActionResultCodeMessage{Code: 50005, Message: "系统维护中，暂不可用"}

	// ---------- WebSocket长连接 60xxx ----------
	WsAuthFail    ActionResultCodeMessage = ActionResultCodeMessage{Code: 60001, Message: "长连接鉴权失败"}
	WsSessionLost ActionResultCodeMessage = ActionResultCodeMessage{Code: 60002, Message: "会话已断开"}
	WsParseErr    ActionResultCodeMessage = ActionResultCodeMessage{Code: 60003, Message: "消息解析失败"}
	WsOffline     ActionResultCodeMessage = ActionResultCodeMessage{Code: 60004, Message: "目标客户端不在线"}
)

// ActionResult 接口返回的结构体
type ActionResult struct {
	Code    ActionResultCode
	Message string
	Data    interface{}
	Now     int64
}

func (as *ActionResult) Apply(ctx constrain.IContext) {
	NewJsonResult(as).Apply(ctx)
}
func (as *ActionResult) Error() string {
	return as.Message //fmt.Sprintf("Code=%v,Message=%v,Data=%v", as.Code, as.Message, as.Data)
}
func New(err error, msg string) *ActionResult {
	if err == nil {
		return &ActionResult{
			Code:    Success,
			Message: msg,
			Data:    nil,
			Now:     time.Now().UnixMilli(),
		}
	} else {
		return &ActionResult{
			Code:    Error,
			Message: err.Error(),
			Data:    nil,
			Now:     time.Now().UnixMilli(),
		}
	}
}
func NewNotFound() *ActionResult {
	return NewCodeWithMessage(NotFound, "not found")
}
func NewSuccess(msg string) *ActionResult {
	return &ActionResult{
		Code:    Success,
		Message: msg,
		Data:    nil,
		Now:     time.Now().UnixMilli(),
	}
}
func NewErrorWithCode(code ActionResultCodeMessage) *ActionResult {
	return &ActionResult{
		Code:    code.Code,
		Message: code.Message,
		Data:    nil,
		Now:     time.Now().UnixMilli(),
	}
}
func NewError(err error) *ActionResult {
	if err == nil {
		return NewSuccess("OK")
	} else {
		if v, ok := err.(*ActionResult); ok {
			return v
		} else {
			return &ActionResult{
				Code:    Error,
				Message: err.Error(),
				Data:    nil,
				Now:     time.Now().UnixMilli(),
			}
		}

	}
}
func NewErrorText(text string) *ActionResult {
	return &ActionResult{
		Code:    Error,
		Message: text,
		Data:    nil,
		Now:     time.Now().UnixMilli(),
	}
}
func NewCodeWithMessage(code ActionResultCode, message string) *ActionResult {
	return &ActionResult{
		Code:    code,
		Message: message,
		Data:    nil,
		Now:     time.Now().UnixMilli(),
	}
}

func NewDataMessage(data interface{}, msg string) *ActionResult {
	return &ActionResult{
		Code:    Success,
		Message: msg,
		Data:    data,
		Now:     time.Now().UnixMilli(),
	}
}
func NewData(data interface{}) *ActionResult {
	return &ActionResult{
		Code:    Success,
		Message: "OK",
		Data:    data,
		Now:     time.Now().UnixMilli(),
	}
}
func NewErrorList(errs []error) *ActionResult {
	if len(errs) == 0 {
		return &ActionResult{
			Code:    Success,
			Message: "OK",
			Data:    nil,
			Now:     time.Now().UnixMilli(),
		}
	}
	as := &ActionResult{}
	as.Code = Error

	asList := make([]*ActionResult, 0)
	for i := range errs {
		asList = append(asList, &ActionResult{
			Code:    Error,
			Message: errs[i].Error(),
			Data:    nil,
			Now:     time.Now().UnixMilli(),
		})
	}
	as.Data = asList
	return as
}

func (as *ActionResult) SmartSuccessData(data interface{}) *ActionResult {
	as.Message = "SUCCESS"
	as.Code = Success
	as.Data = data
	as.Now = time.Now().UnixMilli()
	return as
}
func (as *ActionResult) SmartError(err error, successTxt string, data interface{}) *ActionResult {

	if err == nil {
		as.Message = successTxt
		as.Code = Success
		as.Data = data
		as.Now = time.Now().UnixMilli()
	} else {
		as.Message = err.Error()
		as.Code = Fail
		as.Data = data
		as.Now = time.Now().UnixMilli()
	}
	return as
}
func (as *ActionResult) Smart(code ActionResultCode, s string, f string) *ActionResult {
	as.Code = code
	if code == Success {
		as.Message = s
	} else {
		as.Message = f
	}
	as.Now = time.Now().UnixMilli()
	return as
}
func (as *ActionResult) SmartData(code ActionResultCode, s string, f string, data interface{}) *ActionResult {
	as.Code = code
	if code == Success {
		as.Message = s
		as.Data = data
	} else {
		as.Message = f
	}
	as.Now = time.Now().UnixMilli()
	return as
}
