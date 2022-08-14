package extends

import (
	"github.com/nbvghost/dandelion/constrain"
)

type ViewBase struct {
	Name string `json:"-"`
	//PkgPath string `json:"-"`
	//Global Global `json:"-"`
	HtmlMeta         *HtmlMeta
	HtmlMetaCallback func(viewBase ViewBase, meta *HtmlMeta) error `json:"-"`
}

func (m ViewBase) GetResult(context constrain.IContext, viewHandler constrain.IViewHandler) constrain.IResult {
	return nil
}
func (m ViewBase) GetName() string {
	return m.Name
}

/*func (m *ViewBase) SetName(name string) {
	m.Name = name
}
func (m ViewBase) GetPkgPath() string {
	return m.PkgPath
}
func (m *ViewBase) SetPkgPath(path string) {
	m.PkgPath = path
}
func (m *ViewBase) SetGlobal(v interface{}) {
	m.Global = v.(Global)
}
*/
