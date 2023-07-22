package sqltype

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/nbvghost/dandelion/library/dao"
)

type StringArray []string

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		bytes = []byte("[]")
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j StringArray) Value() (driver.Value, error) {
	if j == nil {
		j = make(StringArray, 0)
	}

	return json.Marshal(j)
}

type Number interface {
	int64 | dao.PrimaryKey
}

type PrimaryKeyArray []dao.PrimaryKey

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *PrimaryKeyArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		bytes = []byte("[]")
	}
	err := json.Unmarshal(bytes, &j)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j PrimaryKeyArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
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
