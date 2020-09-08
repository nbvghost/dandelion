package result

type ActionResultCode int

const ActionOK ActionResultCode = 0
const ActionFail ActionResultCode = 1
const ActionSQLError ActionResultCode = 100
const ActionJSONError ActionResultCode = 200

//var dbMap *gorp.DbMap
type ActionResult struct {
	Code    ActionResultCode
	Message string
	Data    interface{}
}

func (as *ActionResult) SmartSuccessData(data interface{}) *ActionResult {
	as.Message = "SUCCESS"
	as.Code = ActionOK
	as.Data = data
	return as
}
func (as *ActionResult) SmartError(err error, successTxt string, data interface{}) *ActionResult {

	if err == nil {
		as.Message = successTxt
		as.Code = ActionOK
		as.Data = data
	} else {
		as.Message = err.Error()
		as.Code = ActionFail
		as.Data = data
	}
	return as
}
func (as *ActionResult) Smart(code ActionResultCode, s string, f string) *ActionResult {
	as.Code = code
	if code == ActionOK {
		as.Message = s
	} else {
		as.Message = f
	}
	return as
}
func (as *ActionResult) SmartData(code ActionResultCode, s string, f string, data interface{}) *ActionResult {
	as.Code = code
	if code == ActionOK {
		as.Message = s
		as.Data = data
	} else {
		as.Message = f
	}
	return as
}
