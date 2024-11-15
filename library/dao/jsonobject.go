package dao

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type IJSONObjectType interface {
	any
}
type JSONObject[T IJSONObjectType] struct {
	Data *T
}


// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *JSONObject[T]) Scan(value interface{}) error {
	if value == nil {
		m.Data = new(T)
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

	if m.Data==nil{
		m.Data = new(T)
	}
	err := json.Unmarshal(bytes, m.Data)
	//*m = JSON(result)
	return err
}

func (m JSONObject[T]) Value() (driver.Value, error) {
	if m.Data==nil {
		return "{}", nil
	}
	b, err := json.Marshal(m.Data)
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

func (m JSONObject[T]) String() string {
	return fmt.Sprintf("%#v", m.Data)
}

// GormDataType gorm common data type
func (JSONObject[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (JSONObject[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "JSON"
}

func (m JSONObject[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if m.Data==nil {
		return gorm.Expr("?", "{}")
	}

	b, err := json.Marshal(m.Data)
	if err != nil {
		b = []byte("{}")
	}

	return gorm.Expr("?", string(b))
}