package sqltype

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type ObjectArray []map[string]interface{}

func (j *ObjectArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, j)
	return err
}

func (j ObjectArray) Value() (driver.Value, error) {
	if j == nil {
		j = make(ObjectArray, 0)
	}
	return json.Marshal(&j)
}
