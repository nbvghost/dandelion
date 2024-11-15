package sqltype

type CustomerService struct {
	Photo         string
	Name          string
	Title         string
	Hide          bool
	SocialAccount []SocialAccount
}

/*type CustomerServiceList []CustomerService

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *CustomerServiceList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j CustomerServiceList) Value() (driver.Value, error) {
	if j == nil {
		j = make(CustomerServiceList, 0)
	}
	return json.Marshal(&j)
}*/
