package user

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/gpa/types"
	"log"
)

type Address struct {
	Organization *model.Organization `mapping:""`
	Post         struct {
		CityName        string //`form:"CityName"`
		Company         string //`form:"Company"`
		CountyName      string //`form:"CountyName"`
		DefaultBilling  bool   //`form:"DefaultBilling"`
		DefaultShipping bool   //`form:"DefaultShipping"`
		Detail          string //`form:"Detail"`
		FirstName       string //`form:"FirstName"`
		LastName        string //`form:"LastName"`
		PostalCode      string //`form:"PostalCode"`
		ProvinceName    string //`form:"ProvinceName"`
		Tel             string //`form:"Tel"`
	} `method:"Post"`
	Put struct {
		ID              types.PrimaryKey
		CityName        string //`form:"CityName"`
		Company         string //`form:"Company"`
		CountyName      string //`form:"CountyName"`
		DefaultBilling  bool   //`form:"DefaultBilling"`
		DefaultShipping bool   //`form:"DefaultShipping"`
		Detail          string //`form:"Detail"`
		FirstName       string //`form:"FirstName"`
		LastName        string //`form:"LastName"`
		PostalCode      string //`form:"PostalCode"`
		ProvinceName    string //`form:"ProvinceName"`
		Tel             string //`form:"Tel"`
	} `method:"Put"`
	Get struct {
		ID types.PrimaryKey `form:"id"`
	} `method:"Get"`
	Delete struct {
		ID types.PrimaryKey `form:"id"`
	} `method:"Delete"`
}

func (m *Address) HandleDelete(context constrain.IContext) (constrain.IResult, error) {
	where := dao.NewWhere()
	where.Eq(`"UserID"`, context.UID())

	err := dao.DeleteBy(singleton.Orm(), &model.Address{}, map[string]any{"UserID": context.UID(), "ID": m.Delete.ID})
	if err != nil {
		return nil, err
	}

	addressList := dao.Find(singleton.Orm(), &model.Address{}).Where(where.String()).List()
	return result.NewData(map[string]any{"AddressList": addressList}), nil
}

func (m *Address) Handle(context constrain.IContext) (r constrain.IResult, err error) {
	where := dao.NewWhere()
	where.Eq(`"UserID"`, context.UID())
	if m.Get.ID > 0 {
		where.Eq(`"ID"`, m.Get.ID)
	}
	addressList := dao.Find(singleton.Orm(), &model.Address{}).Where(where.String()).List()
	return result.NewData(map[string]any{"AddressList": addressList}), nil
}

func (m *Address) HandlePut(context constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(context)
	log.Println(contextValue)

	address := map[string]any{
		"ID":              m.Put.ID,
		"UserID":          context.UID(),
		"Name":            m.Put.LastName + " " + m.Put.FirstName,
		"CountyName":      m.Put.CountyName,
		"ProvinceName":    m.Put.ProvinceName,
		"CityName":        m.Put.CityName,
		"Detail":          m.Put.Detail,
		"PostalCode":      m.Put.PostalCode,
		"Tel":             m.Put.Tel,
		"Company":         m.Put.Company,
		"DefaultBilling":  m.Put.DefaultBilling,
		"DefaultShipping": m.Put.DefaultShipping,
	}

	tx := singleton.Orm().Begin()
	if m.Put.DefaultBilling || m.Put.DefaultShipping {
		changeValue := map[string]any{}
		if m.Put.DefaultBilling {
			changeValue["DefaultBilling"] = false
		}
		if m.Put.DefaultShipping {
			changeValue["DefaultShipping"] = false
		}
		if len(changeValue) > 0 {
			err := dao.UpdateBy(tx, &model.Address{}, changeValue, map[string]any{"UserID": context.UID()})
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	err := dao.UpdateByPrimaryKey(tx, &model.Address{}, m.Put.ID, address)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return result.NewSuccess("OK"), nil
}
func (m *Address) HandlePost(context constrain.IContext) (constrain.IResult, error) {
	contextValue := contexext.FromContext(context)
	log.Println(contextValue)

	address := &model.Address{
		UserID:          context.UID(),
		Name:            m.Post.LastName + " " + m.Post.FirstName,
		CountyName:      m.Post.CountyName,
		ProvinceName:    m.Post.ProvinceName,
		CityName:        m.Post.CityName,
		Detail:          m.Post.Detail,
		PostalCode:      m.Post.PostalCode,
		Tel:             m.Post.Tel,
		Company:         m.Post.Company,
		DefaultBilling:  m.Post.DefaultBilling,
		DefaultShipping: m.Post.DefaultShipping,
	}

	tx := singleton.Orm().Begin()
	if m.Post.DefaultBilling || m.Post.DefaultShipping {
		changeValue := map[string]any{}
		if m.Post.DefaultBilling {
			changeValue["DefaultBilling"] = false
		}
		if m.Post.DefaultShipping {
			changeValue["DefaultShipping"] = false
		}
		if len(changeValue) > 0 {
			err := dao.UpdateBy(tx, &model.Address{}, changeValue, map[string]any{"UserID": context.UID()})
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	err := dao.Create(tx, address)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return result.NewSuccess("OK"), nil
}
