package dao

type ContentItemContentSubType struct {
	ContentItem    `gorm:"column:ContentItem"`
	ContentSubType `gorm:"column:ContentSubType"`
}
