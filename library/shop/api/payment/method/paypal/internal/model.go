package internal

import (
	"fmt"
	"github.com/nbvghost/dandelion/entity/model"
)

type ItemTotal struct {
	CurrencyCode string `json:"currency_code,omitempty"`
	Value        string `json:"value,omitempty"`
}

type Breakdown struct {
	ItemTotal ItemTotal `json:"item_total,omitempty"`
}

type Amount struct {
	CurrencyCode string    `json:"currency_code,omitempty"`
	Value        string    `json:"value,omitempty"`
	Breakdown    Breakdown `json:"breakdown,omitempty"`
}
type Name struct {
	GivenName string `json:"given_name,omitempty"` //名
	Surname   string `json:"surname,omitempty"`    //姓
	FullName  string `json:"full_name,omitempty"`
}

func (m *Name) GetFullName() string {
	return m.GivenName + " " + m.Surname
}

type Shipping struct {
	Name    *Name    `json:"name,omitempty"`
	Type    string   `json:"type,omitempty"` //SHIPPING
	Address *Address `json:"address,omitempty"`
}
type CheckoutOrdersUnitItemUnitAmount struct {
	CurrencyCode string `json:"currency_code,omitempty"`
	Value        string `json:"value,omitempty"`
}
type CheckoutOrdersUnitItem struct {
	Name        string                           `json:"name,omitempty"`
	Quantity    string                           `json:"quantity,omitempty"`
	Description string                           `json:"description,omitempty"`
	Sku         string                           `json:"sku,omitempty"`
	Url         string                           `json:"url,omitempty"`
	Category    string                           `json:"category,omitempty"` //DIGITAL_GOODS,PHYSICAL_GOODS,DONATION
	ImageUrl    string                           `json:"image_url,omitempty"`
	UnitAmount  CheckoutOrdersUnitItemUnitAmount `json:"unit_amount,omitempty"`
}
type CheckoutOrdersUnit struct {
	ReferenceId string                   `json:"reference_id,omitempty"`
	Description string                   `json:"description,omitempty"`
	Amount      Amount                   `json:"amount,omitempty"`
	Shipping    *Shipping                `json:"shipping,omitempty"`
	Items       []CheckoutOrdersUnitItem `json:"items,omitempty"`
}
type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

type Payer struct {
	PayerId      string `json:"payer_id,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	Name         Name   `json:"name,omitempty"`
	Phone        struct {
		PhoneType   string `json:"phone_type,omitempty"`
		PhoneNumber struct {
			NationalNumber string `json:"national_number,omitempty"`
		} `json:"phone_number,omitempty"`
	} `json:"phone,omitempty"`
	BirthDate string `json:"birth_date,omitempty"`
	TaxInfo   struct {
		TaxId     string `json:"tax_id,omitempty"`
		TaxIdType string `json:"tax_id_type,omitempty"`
	} `json:"tax_info,omitempty"`
	Address Address `json:"address,omitempty"`
}
type Address struct {
	AddressLine1 string `json:"address_line_1,omitempty"`
	AddressLine2 string `json:"address_line_2,omitempty"`
	AdminArea2   string `json:"admin_area_2,omitempty"`
	AdminArea1   string `json:"admin_area_1,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	CountryCode  string `json:"country_code,omitempty"`
}

func (m *Address) SetAddress(address *model.Address) *Address {
	if address == nil {
		return nil
	}
	m.AddressLine1 = address.Detail
	if len(address.Company) > 0 {
		m.AddressLine2 = fmt.Sprintf("(%s)", address.Company)
	}
	m.AdminArea1 = address.CountyName + "." + address.ProvinceName
	m.AdminArea2 = address.CityName
	m.PostalCode = address.PostalCode
	m.CountryCode = address.CountyCode
	return m
}

type CheckoutOrdersCard struct {
	Name           string  `json:"name,omitempty"`
	Number         string  `json:"number,omitempty"`
	SecurityCode   string  `json:"security_code,omitempty"`
	Expiry         string  `json:"expiry,omitempty"`
	BillingAddress Address `json:"billing_address,omitempty"`
}
type CheckoutOrdersPaymentSource struct {
	Card CheckoutOrdersCard `json:"card,omitempty"`
}
