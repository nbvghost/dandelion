package model

import (
	"encoding/json"
	"github.com/lib/pq"
	"github.com/nbvghost/dandelion/entity/sqltype"
	"github.com/nbvghost/dandelion/library/dao"
)

//Content   ContentType  ContentSubType

// ContentItem Menus
type ContentItem struct {
	dao.Entity
	OID           dao.PrimaryKey  `gorm:"column:OID;index"`
	Uri           string          `gorm:"column:Uri"`
	Name          string          `gorm:"column:Name"`
	Sort          int             `gorm:"column:Sort"`
	ContentTypeID dao.PrimaryKey  `gorm:"column:ContentTypeID"`
	Type          ContentTypeType `gorm:"column:Type"`
	TemplateName  string          `gorm:"column:TemplateName"` //使用模板的文件名，如多文章列表，
	Introduction  string          `gorm:"column:Introduction"` //主类介绍
	Image         string          `gorm:"column:Image"`        //主类图片
	Badge         string          `gorm:"column:Badge"`        //徽章
	//Hide            bool                         `gorm:"column:Hide"`
	ShowAtMenu      bool                                  `gorm:"column:ShowAtMenu"`
	ShowAtHome      bool                                  `gorm:"column:ShowAtHome"`
	CustomizeFields sqltype.Array[sqltype.CustomizeField] `gorm:"column:CustomizeFields;type:JSON"`
	Config          string                                `gorm:"column:Config;type:JSON"`
}

func (ContentItem) TableName() string {
	return "ContentItem"
}

type ContentSubType struct {
	dao.Entity
	OID                    dao.PrimaryKey `gorm:"column:OID;index"`
	Uri                    string         `gorm:"column:Uri"`
	Name                   string         `gorm:"column:Name"`
	ContentItemID          dao.PrimaryKey `gorm:"column:ContentItemID"`
	ParentContentSubTypeID dao.PrimaryKey `gorm:"column:ParentContentSubTypeID"`
	Sort                   int64          `gorm:"column:Sort"`
}

func (ContentSubType) TableName() string {
	return "ContentSubType"
}

type ContentTypeType string

func (m ContentTypeType) String() string {
	return string(m)
}

const (
	ContentTypeContents ContentTypeType = "contents"
	ContentTypeContent  ContentTypeType = "content"
	ContentTypeIndex    ContentTypeType = "index"
	ContentTypeGallery  ContentTypeType = "gallery"
	ContentTypeProducts ContentTypeType = "products"
	ContentTypeBlog     ContentTypeType = "blog"
	ContentTypePage     ContentTypeType = "page"
)

// MenuType
type ContentType struct {
	dao.Entity
	Label string          `gorm:"column:Label"`
	Type  ContentTypeType `gorm:"column:Type;unique"`
}

func (ContentType) TableName() string {
	return "ContentType"
}

type Content struct {
	dao.Entity
	OID              dao.PrimaryKey `gorm:"column:OID;index"`          //
	Uri              string         `gorm:"column:Uri"`                //
	Title            string         `gorm:"column:Title"`              //
	Summary          string         `gorm:"column:Summary"`            //
	Content          string         `gorm:"column:Content"`            //
	Picture          string         `gorm:"column:Picture"`            //
	Images           pq.StringArray `gorm:"column:Images;type:text[]"` //多图展示
	ContentItemID    dao.PrimaryKey `gorm:"column:ContentItemID"`      //
	ContentSubTypeID dao.PrimaryKey `gorm:"column:ContentSubTypeID"`   //
	FromUrl          string         `gorm:"column:FromUrl"`            //
	Author           string         `gorm:"column:Author"`             //

	Keywords    string `gorm:"column:Keywords"`    //
	Description string `gorm:"column:Description"` //

	IsStickyTop  bool           `gorm:"column:IsStickyTop"`                      //
	CountView    int            `gorm:"column:CountView"`                        //
	CountLike    int            `gorm:"column:CountLike"`                        //
	CountShare   int            `gorm:"column:CountShare"`                       //
	FieldGroupID dao.PrimaryKey `gorm:"column:FieldGroupID"`                     //
	FieldData    string         `gorm:"column:FieldData;type:JSON;default:'{}'"` //json
	Tags         pq.StringArray `gorm:"column:Tags;type:text[];default:'{}'"`    //
}

func (m *Content) GetFieldData() []any {
	var arr []any
	json.Unmarshal([]byte(m.FieldData), &arr)
	return arr
}

func (Content) TableName() string {
	return "Content"
}

type ContentConfig struct {
	dao.Entity
	OID                 dao.PrimaryKey                         `gorm:"column:OID;unique"`
	Name                string                                 `gorm:"column:Name"`
	Logo                string                                 `gorm:"column:Logo"`
	FaviconIco          string                                 `gorm:"column:FaviconIco"`
	SocialAccount       sqltype.Array[sqltype.SocialAccount]   `gorm:"column:SocialAccount;type:JSON"`
	CustomerService     sqltype.Array[sqltype.CustomerService] `gorm:"column:CustomerService;type:JSON"`
	EnableHTMLCache     bool                                   `gorm:"column:EnableHTMLCache"`
	EnableMultiLanguage bool                                   `gorm:"column:EnableMultiLanguage"`
	FocusPicture        sqltype.Array[sqltype.FocusPicture]    `gorm:"column:FocusPicture;type:JSON"`
	TemplateConfig      string                                 `gorm:"column:TemplateConfig;type:JSON;default:'{}'"`
}

func (ContentConfig) TableName() string {
	return "ContentConfig"
}
