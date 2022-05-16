package result

import (
	"fmt"

	"github.com/nbvghost/dandelion/constrain"
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

//var dbMap *gorp.DbMap
type ActionResult struct {
	Code    ActionResultCode
	Message string
	Data    interface{}
}

func (as *ActionResult) Apply(ctx constrain.IContext) {
	NewJsonResult(as).Apply(ctx)
}
func (as *ActionResult) Error() string {

	return fmt.Sprintf("Code=%v,Message=%v,Data=%v", as.Code, as.Message, as.Data)
}
func New(err error, msg string) *ActionResult {
	if err == nil {
		return &ActionResult{
			Code:    Success,
			Message: msg,
			Data:    nil,
		}
	} else {
		return &ActionResult{
			Code:    Error,
			Message: err.Error(),
			Data:    nil,
		}
	}
}
func NewSuccess(msg string) *ActionResult {
	return &ActionResult{
		Code:    Success,
		Message: msg,
		Data:    nil,
	}
}
func NewError(err error) *ActionResult {
	if err == nil {
		return &ActionResult{
			Code:    Success,
			Message: "OK",
			Data:    nil,
		}
	} else {
		return &ActionResult{
			Code:    Error,
			Message: err.Error(),
			Data:    nil,
		}
	}
}
func NewErrorText(text string) *ActionResult {
	return &ActionResult{
		Code:    Error,
		Message: text,
		Data:    nil,
	}
}
func NewCodeWithError(code ActionResultCode, err error) *ActionResult {
	if err == nil {
		if code == Success {
			return &ActionResult{
				Code:    Success,
				Message: "OK",
				Data:    nil,
			}
		} else {
			return &ActionResult{
				Code:    code,
				Message: "执行失败",
				Data:    nil,
			}
		}

	} else {
		return &ActionResult{
			Code:    code,
			Message: err.Error(),
			Data:    nil,
		}
	}
}
func NewData(data interface{}) *ActionResult {
	return &ActionResult{
		Code:    Success,
		Message: "OK",
		Data:    data,
	}
}
func NewErrorList(errs []error) *ActionResult {
	if len(errs) == 0 {
		return &ActionResult{
			Code:    Success,
			Message: "OK",
			Data:    nil,
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
		})
	}
	as.Data = asList
	return as
}

func (as *ActionResult) SmartSuccessData(data interface{}) *ActionResult {
	as.Message = "SUCCESS"
	as.Code = Success
	as.Data = data
	return as
}
func (as *ActionResult) SmartError(err error, successTxt string, data interface{}) *ActionResult {

	if err == nil {
		as.Message = successTxt
		as.Code = Success
		as.Data = data
	} else {
		as.Message = err.Error()
		as.Code = Fail
		as.Data = data
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
	return as
}
