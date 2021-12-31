package base

type ViewBase struct {
	StructName string `json:"-"`
}

func (m ViewBase) GetName() string {
	return m.StructName
}
func (m *ViewBase) SetName(name string) {
	m.StructName = name
}
