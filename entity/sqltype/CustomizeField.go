package sqltype

type CustomizeField struct {
	Name string
	Type string
}

/*type CustomizeFieldList []CustomizeField

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *CustomizeFieldList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j CustomizeFieldList) Value() (driver.Value, error) {
	if j == nil {
		j = make(CustomizeFieldList, 0)
	}
	return json.Marshal(&j)
}*/
