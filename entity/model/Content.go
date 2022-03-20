package model

import (
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/gpa/types"
)

//Content   ContentType  ContentSubType

//Menus
type ContentItem struct {
	base.BaseModel
	OID           types.PrimaryKey `gorm:"column:OID"`
	Name          string           `gorm:"column:Name"`
	Sort          int              `gorm:"column:Sort"`
	ContentTypeID types.PrimaryKey `gorm:"column:ContentTypeID"`
	TemplateName  string           `gorm:"column:TemplateName"` //使用模板的文件名，如多文章列表，
	Type          ContentTypeType  `gorm:"column:Type"`
	Hide          bool             `gorm:"column:Hide"`
}

func (ContentItem) TableName() string {
	return "ContentItem"
}

type ContentTypeType string

const (
	ContentTypeContents ContentTypeType = "contents"
	ContentTypeContent  ContentTypeType = "content"
	ContentTypeIndex    ContentTypeType = "index"
	ContentTypeGallery  ContentTypeType = "gallery"
	ContentTypeProducts ContentTypeType = "products"
	ContentTypeBlog     ContentTypeType = "blog"
	ContentTypePage     ContentTypeType = "page"
)

//MenuType
type ContentType struct {
	base.BaseModel
	Label string          `gorm:"column:Label"`
	Type  ContentTypeType `gorm:"column:Type;unique"`
}

func (ContentType) TableName() string {
	return "ContentType"
}

type ContentSubType struct {
	types.Entity
	OID                    types.PrimaryKey `gorm:"column:OID"`
	Name                   string           `gorm:"column:Name"`
	ContentItemID          types.PrimaryKey `gorm:"column:ContentItemID"`
	ParentContentSubTypeID types.PrimaryKey `gorm:"column:ParentContentSubTypeID"`
	Sort                   int              `gorm:"column:Sort"`
}

func (ContentSubType) TableName() string {
	return "ContentSubType"
}

type Content struct {
	types.Entity
	OID              types.PrimaryKey `gorm:"column:OID"`
	Title            string           `gorm:"column:Title"`
	Keywords         string           `gorm:"column:Keywords"`
	Description      string           `gorm:"column:Description"`
	Summary          string           `gorm:"column:Summary"`
	Content          string           `gorm:"column:Content"`
	Picture          string           `gorm:"column:Picture"`
	ContentItemID    types.PrimaryKey `gorm:"column:ContentItemID"`
	ContentSubTypeID types.PrimaryKey `gorm:"column:ContentSubTypeID"`
	FromUrl          string           `gorm:"column:FromUrl"`
	Author           string           `gorm:"column:Author"`
	CountView        int              `gorm:"column:CountView"`
	CountLike        int              `gorm:"column:CountLike"`
	Tags             pq.StringArray   `gorm:"column:Tags;type:text[];default:'{}'"` //
}

func (Content) TableName() string {
	return "Content"
}

type ContentConfig struct {
	types.Entity
	OID             types.PrimaryKey            `gorm:"column:OID;unique"`
	Name            string                      `gorm:"column:Name"`
	Logo            string                      `gorm:"column:Logo"`
	SocialAccount   sqltype.SocialAccountList   `gorm:"column:SocialAccount;type:JSON"`
	CustomerService sqltype.CustomerServiceList `gorm:"column:CustomerService;type:JSON"`
	EnableHTMLCache bool                        `gorm:"column:EnableHTMLCache"`
	FocusPicture    sqltype.FocusPictureList    `gorm:"column:FocusPicture;type:JSON"`
}

func (ContentConfig) TableName() string {
	return "ContentConfig"
}
