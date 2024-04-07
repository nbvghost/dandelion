package model

import (
	"database/sql/driver"
	"encoding/json"
)

type ITableType interface {
	~*Specification | ~*GoodsAttributes | ~*GoodsSkuLabel | ~*GoodsSkuLabelData
}

type Array[T ITableType] []T

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *Array[T]) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		bytes = []byte("[]")
	}
	err := json.Unmarshal(bytes, &j)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j *Array[T]) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
