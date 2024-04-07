package sqltype

type Route struct {
	Path        string
	Permissions []string //拥有的权限['edit', 'add', 'delete']
	Children    []Route
}
/*type RouteArray []Route

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (m *RouteArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, m)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (m *RouteArray) Value() (driver.Value, error) {
	if m == nil {
		c:=make(RouteArray, 0)
		m = &c
	}
	return json.Marshal(m)
}
*/