package result

import (
	"github.com/nbvghost/dandelion/constrain"
	"time"
)

type ActionResultCode int

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
func NewSuccess(msg string) *ActionResult {
	return &ActionResult{
		Code:    Success,
		Message: msg,
		Data:    nil,
		Now:     time.Now().UnixMilli(),
	}
}
func NewError(err error) *ActionResult {
	if err == nil {
		return &ActionResult{
			Code:    Success,
			Message: "OK",
			Data:    nil,
			Now:     time.Now().UnixMilli(),
		}
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
