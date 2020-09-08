package dao

//Content   ContentType  ContentSubType

//Menus
type ContentItem struct {
	BaseModel
	OID           uint64          `gorm:"column:OID"`
	Name          string          `gorm:"column:Name"`
	Sort          int             `gorm:"column:Sort"`
	ContentTypeID uint64          `gorm:"column:ContentTypeID"`
	TemplateName  string          `gorm:"column:TemplateName"` //使用模板的文件名，如多文章列表，
	Type          ContentTypeType `gorm:"column:Type"`
	Hide          bool            `gorm:"column:Hide"`
}

func (ContentItem) TableName() string {
	return "ContentItem"
}

type ContentTypeType string

const ContentTypeContents ContentTypeType = "contents"
const ContentTypeContent ContentTypeType = "content"
const ContentTypeIndex ContentTypeType = "index"
const ContentTypeGallery ContentTypeType = "gallery"
const ContentTypeProducts ContentTypeType = "products"

//MenuType
type ContentType struct {
	BaseModel
	Label string          `gorm:"column:Label"`
	Type  ContentTypeType `gorm:"column:Type;unique"`
}

func (ContentType) TableName() string {
	return "ContentType"
}

//Classify
type ContentSubType struct {
	BaseModel
	Name                   string `gorm:"column:Name"`
	ContentItemID          uint64 `gorm:"column:ContentItemID"`
	ParentContentSubTypeID uint64 `gorm:"column:ParentContentSubTypeID"`
	Sort                   int    `gorm:"column:Sort"`
}

func (ContentSubType) TableName() string {
	return "ContentSubType"
}

type Content struct {
	BaseModel
	Title            string `gorm:"column:Title"`
	Keywords         string `gorm:"column:Keywords"`
	Description      string `gorm:"column:Description"`
	Introduce        string `gorm:"column:Introduce"`
	Content          string `gorm:"column:Content;type:LONGTEXT"`
	Picture          string `gorm:"column:Picture"`
	ContentItemID    uint64 `gorm:"column:ContentItemID"`
	ContentSubTypeID uint64 `gorm:"column:ContentSubTypeID"`
	FromUrl          string `gorm:"column:FromUrl"`
	Author           string `gorm:"column:Author"`
	Look             int    `gorm:"column:Look"`
}

func (Content) TableName() string {
	return "Content"
}
