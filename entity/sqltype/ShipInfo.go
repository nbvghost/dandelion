package sqltype

/*type ShipInfo struct {

}*/

/*type ShipInfo struct {
	No   string //快递单号
	Name string //快递
	Key  string //快递编号
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *ShipInfo) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j ShipInfo) Value() (driver.Value, error) {
	return json.Marshal(j)
}

*/
/*
// RefundInfo 退货信息
type RefundInfo struct {
	Status   RefundStatus
	ShipInfo ShipInfo
	HasGoods bool
	Reason   string //原因
	AskTime  time.Time
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (j *RefundInfo) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	err := json.Unmarshal(bytes, j)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (j RefundInfo) Value() (driver.Value, error) {

	return json.Marshal(j)
}*/
