package sqltype

import (
	"database/sql/driver"
	"encoding/json"
)

type Map map[string]any

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *Map) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		bytes = []byte("{}")
	}
	if j == nil {
		*j = make(Map)
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j Map) Value() (driver.Value, error) {
	if j == nil {
		j = make(Map)
	}

	return json.Marshal(j)
}

/*type InterfaceArray []interface{}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j InterfaceArray) ConvertValue(v interface{}) (driver.Value, error) {
	if j == nil {
		j = make(InterfaceArray, 0)
	}

	return "1,3", nil
}
func (j InterfaceArray) Value() (driver.Value, error) {
	if j == nil {
		j = make(InterfaceArray, 0)
	}

	return "1,3", nil
}
*/
