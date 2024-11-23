package sqltype

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbvghost/dandelion/library/dao"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type IMapKey interface {
	string | dao.PrimaryKey | int | uint
}

type Object[K IMapKey, V IMetaType] map[K]V

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *Object[K, V]) Scan(value any) error {
	if value == nil {
		*m = make(Object[K, V])
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		if len(v) > 0 {
			bytes = make([]byte, len(v))
			copy(bytes, v)
		}
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	err := json.Unmarshal(bytes, m)
	//*m = JSON(result)
	return err
}

func (m Object[K, V]) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "{}", nil
	}
	return b, err
	//return string(m), nil
}

// MarshalJSON to output non base64 encoded []byte
/*func (m Array[T]) MarshalJSON() ([]byte, error) {
	return json.RawMessage(m).MarshalJSON()
}*/

// UnmarshalJSON to deserialize []byte
/*func (m *Array[T]) UnmarshalJSON(b []byte) error {
	result := json.RawMessage{}
	err := result.UnmarshalJSON(b)
	*m = JSON(result)
	return err
}*/

func (m Object[K, V]) String() string {
	return fmt.Sprintf("%#v", m)
}

// GormDataType gorm common data type
func (Object[K, V]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (Object[K, V]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "JSON"
}

func (m Object[K, V]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(m) == 0 {
		return gorm.Expr("?", "{}")
	}
	b, err := json.Marshal(m)
	if err != nil {
		b = []byte("{}")
	}
	return gorm.Expr("?", string(b))
}
