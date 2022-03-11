package result

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/library/contexext"
	"net/http"
)

const MIME_APPLICATION_JSON byte = 1
const MIME_TEXT_PLAIN byte = 2

var _ constrain.IResult = (*JsonResult)(nil)

type Head struct {
	Mine byte
}

func (m *Head) ToData(b []byte) error {
	buffer := bytes.NewBuffer(b)
	if mine, err := buffer.ReadByte(); err != nil {
		return err
	} else {
		m.Mine = mine
	}
	return nil
}
func (m *Head) ToBytes() []byte {
	return []byte{m.Mine}
}

func UnmarshalResult(b []byte) ([]byte, *Head, error) {
	buffer := bytes.NewBuffer(b)
	var headLen uint64
	err := binary.Read(buffer, binary.BigEndian, &headLen)
	if err != nil {
		return nil, nil, err
	}

	headBytes := make([]byte, headLen)
	if _, err := buffer.Read(headBytes); err != nil {
		return nil, nil, err
	}

	var head Head
	if err := head.ToData(headBytes); err != nil {
		return nil, nil, err
	}

	dataBytes := make([]byte, buffer.Len())
	if _, err := buffer.Read(dataBytes); err != nil {
		return nil, nil, err
	}
	return dataBytes, &head, nil
}
func MarshalResult(b []byte, head *Head) ([]byte, error) {
	headBytes := head.ToBytes()
	buffer := bytes.NewBuffer(nil)
	var headLen = uint64(len(headBytes))
	err := binary.Write(buffer, binary.BigEndian, &headLen)
	if err != nil {
		return nil, err
	}
	buffer.Write(headBytes)
	buffer.Write(b)
	return buffer.Bytes(), nil
}

type JsonResult struct {
	error
	Data interface{}
	///sync.RWMutex
}

func (r *JsonResult) Apply(context constrain.IContext) {
	v := contexext.FromContext(context)

	var b []byte
	var err error

	b, err = json.Marshal(r.Data)
	if err != nil {
		(&ErrorResult{Error: err}).Apply(context)
		return
	}
	//return buffer.Bytes(), err
	//b, err = json.Marshal(r.Data)
	//b = buffer.Bytes()

	v.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	v.Response.WriteHeader(http.StatusOK)
	//context.Response.Header().Add("Content-Type", "application/json")
	v.Response.Write(b)
}

type TextResult struct {
	Data string
}

func (r *TextResult) Apply(context constrain.IContext) {
	v := contexext.FromContext(context)
	v.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	v.Response.WriteHeader(http.StatusOK)
	v.Response.Write([]byte(r.Data))
}

type ErrorResult struct {
	Error error
}

func (r *ErrorResult) Apply(context constrain.IContext) {
	v := contexext.FromContext(context)
	if r.Error != nil {
		http.Error(v.Response, r.Error.Error(), http.StatusNotFound)
	} else {
		http.Error(v.Response, "error", http.StatusNotFound)
	}
}
func NewErrorResult(err error) *ErrorResult {
	return &ErrorResult{Error: err}
}

type EmptyResult struct {
}

func (r *EmptyResult) Apply(context constrain.IContext) {
	v := contexext.FromContext(context)
	v.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	v.Response.WriteHeader(http.StatusOK)
	v.Response.Write([]byte("{}"))
}

type ImageBytesResult struct {
	Data        []byte
	ContentType string //: image/png
}

func (r *ImageBytesResult) Apply(context constrain.IContext) {
	v := contexext.FromContext(context)
	//context.Response.Header().Add()
	v.Response.Header().Set("Content-Type", r.ContentType)
	v.Response.Write(r.Data)

}
