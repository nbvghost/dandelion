package extends

type Tag struct {
	Name  string `gorm:"column:Name"`
	Count int    `gorm:"column:Count"`
	Uri   string
}
