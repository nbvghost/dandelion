package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/nbvghost/dandelion/entity/sqltype"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type GoodsSku struct {
	GoodsSkuLabel     *GoodsSkuLabel
	GoodsSkuLabelData *GoodsSkuLabelData
}

type ITableType interface {
	Specification | GoodsAttributes | GoodsSkuLabel | GoodsSkuLabelData | GoodsSku
}

type Array[T ITableType] []T

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *Array[T]) Scan(value interface{}) error {
	if value == nil {
		*m = make(Array[T], 0)
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

func (m Array[T]) Value() (driver.Value, error) {
	if len(m) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "[]", nil
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

func (m Array[T]) String() string {
	return fmt.Sprintf("%#v", m)
}

// GormDataType gorm common data type
func (Array[T]) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (Array[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return "JSON"
}

func (m Array[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	if len(m) == 0 {
		return gorm.Expr("?", "[]")
	}

	b, err := json.Marshal(m)
	if err != nil {
		b = []byte("[]")
	}

	return gorm.Expr("?", string(b))
}

var _ gorm.Valuer = sqltype.Array[sqltype.CustomerService]{}
var _ driver.Value = sqltype.Array[sqltype.CustomerService]{}
var _ schema.GormDataTypeInterface = sqltype.Array[sqltype.CustomerService]{}

var _ gorm.Valuer = Array[Specification]{}
var _ driver.Value = Array[Specification]{}
var _ schema.GormDataTypeInterface = Array[Specification]{}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
/*func (j *Array[T]) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		bytes = []byte("[]")
	}
	err := json.Unmarshal(bytes, &j)
	return err
}*/

// Value 实现 driver.Valuer 接口，Value 返回 json value
/*func (j Array[T]) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	if err != nil {
		return "[]", nil
	}
	return b, err
}*/
