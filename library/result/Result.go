package result

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/nbvghost/dandelion/constrain"
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
	Data interface{}
}

func (r *JsonResult) Apply(context constrain.IContext) ([]byte, error) {
	b, err := json.Marshal(r.Data)
	if err != nil {
		return nil, err
	}
	packageBytes, err := MarshalResult(b, &Head{Mine: MIME_APPLICATION_JSON})
	if err != nil {
		return nil, err
	}
	return packageBytes, nil
}

type TextResult struct {
	Data string
}

func (r *TextResult) Apply(context constrain.IContext) ([]byte, error) {
	packageBytes, err := MarshalResult([]byte(r.Data), &Head{Mine: MIME_TEXT_PLAIN})
	if err != nil {
		return nil, err
	}
	return packageBytes, nil
}
