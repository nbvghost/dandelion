package model

import (
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/entity/base"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/gpa/types"
)

//Content   ContentType  ContentSubType

// ContentItem Menus
type ContentItem struct {
	types.Entity
	OID             types.PrimaryKey             `gorm:"column:OID;index"`
	Uri             string                       `gorm:"column:Uri"`
	Name            string                       `gorm:"column:Name"`
	Sort            int                          `gorm:"column:Sort"`
	ContentTypeID   types.PrimaryKey             `gorm:"column:ContentTypeID"`
	Type            ContentTypeType              `gorm:"column:Type"`
	TemplateName    string                       `gorm:"column:TemplateName"` //使用模板的文件名，如多文章列表，
	Introduction    string                       `gorm:"column:Introduction"` //主类介绍
	Image           string                       `gorm:"column:Image"`        //主类图片
	Hide            bool                         `gorm:"column:Hide"`
	CustomizeFields []sqltype.CustomizeFieldList `gorm:"column:CustomizeFields;type:JSON"`
}

func (ContentItem) TableName() string {
	return "ContentItem"
}

type ContentSubType struct {
	types.Entity
	OID                    types.PrimaryKey `gorm:"column:OID;index"`
	Uri                    string           `gorm:"column:Uri"`
	Name                   string           `gorm:"column:Name"`
	ContentItemID          types.PrimaryKey `gorm:"column:ContentItemID"`
	ParentContentSubTypeID types.PrimaryKey `gorm:"column:ParentContentSubTypeID"`
	Sort                   int              `gorm:"column:Sort"`
}

func (ContentSubType) TableName() string {
	return "ContentSubType"
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

type Content struct {
	types.Entity
	OID              types.PrimaryKey `gorm:"column:OID;index"`                     //
	Uri              string           `gorm:"column:Uri"`                           //
	Title            string           `gorm:"column:Title"`                         //
	Summary          string           `gorm:"column:Summary"`                       //
	Content          string           `gorm:"column:Content"`                       //
	Picture          string           `gorm:"column:Picture"`                       //
	ContentItemID    types.PrimaryKey `gorm:"column:ContentItemID"`                 //
	ContentSubTypeID types.PrimaryKey `gorm:"column:ContentSubTypeID"`              //
	FromUrl          string           `gorm:"column:FromUrl"`                       //
	Author           string           `gorm:"column:Author"`                        //
	CountView        int              `gorm:"column:CountView"`                     //
	CountLike        int              `gorm:"column:CountLike"`                     //
	Tags             pq.StringArray   `gorm:"column:Tags;type:text[];default:'{}'"` //
}

func (Content) TableName() string {
	return "Content"
}

type ContentConfig struct {
	types.Entity
	OID                 types.PrimaryKey            `gorm:"column:OID;unique"`
	Name                string                      `gorm:"column:Name"`
	Logo                string                      `gorm:"column:Logo"`
	FaviconIco          string                      `gorm:"column:FaviconIco"`
	SocialAccount       sqltype.SocialAccountList   `gorm:"column:SocialAccount;type:JSON"`
	CustomerService     sqltype.CustomerServiceList `gorm:"column:CustomerService;type:JSON"`
	EnableHTMLCache     bool                        `gorm:"column:EnableHTMLCache"`
	EnableMultiLanguage bool                        `gorm:"column:EnableMultiLanguage"`
	FocusPicture        sqltype.FocusPictureList    `gorm:"column:FocusPicture;type:JSON"`
}

func (ContentConfig) TableName() string {
	return "ContentConfig"
}
