package sqltype

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type SocialAccount struct {
	Type SocialType
	Hide bool
	Account string
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *SocialAccount) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j *SocialAccount) Value() (driver.Value, error) {
	if j == nil {
		j = &SocialAccount{}
	}

	return json.Marshal(j)
}
/*
type SocialAccountList []SocialAccount

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *SocialAccountList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j SocialAccountList) Value() (driver.Value, error) {
	if j == nil {
		j = make(SocialAccountList, 0)
	}

	return json.Marshal(j)
}*/
