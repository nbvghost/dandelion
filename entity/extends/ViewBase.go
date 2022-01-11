package extends

type ViewBase struct {
	Name    string `json:"-"`
	PkgPath string `json:"-"`
	Global  Global `json:"-"`
}

func (m ViewBase) GetName() string {
	return m.Name
}
func (m *ViewBase) SetName(name string) {
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
