package sqltype

type FocusPicture struct {
	Url       string
	Title     string
	Introduce string
	Link      string
	Hide      bool
}

/*type FocusPictureList []FocusPicture

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *FocusPictureList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (j FocusPictureList) Value() (driver.Value, error) {
	if j == nil {
		j = make(FocusPictureList, 0)
	}
	return json.Marshal(j)
}*/
