package result

import (
	"github.com/nbvghost/dandelion/constrain"
)

type fileDownloadViewResult struct {
	Data     []byte
	Filename string
}

func (m *fileDownloadViewResult) GetName() string {
	return ""
}

func (m *fileDownloadViewResult) GetResult(context constrain.IContext, viewHandler constrain.IViewHandler) constrain.IResult {
	return NewFileDownloadResult(m.Data, m.Filename)
}

func (m *fileDownloadViewResult) GetContentType() string {
	return ""
}
func NewFileDownloadViewResult(d []byte, filename string) constrain.IViewResult {
	return &fileDownloadViewResult{Data: d, Filename: filename}
}
