package result

import "encoding/json"

var _ Result = (*JsonResult)(nil)

type Result interface {
	Apply() ([]byte, error)
}
type JsonResult struct {
	error
	Data interface{}
	///sync.RWMutex
}

func (r *JsonResult) Apply() ([]byte, error) {
	b, err := json.Marshal(r.Data)
	return b, err
}
