package extends

import "strings"

type Address struct {
	Name         string
	ProvinceName string
	CityName     string
	CountyName   string
	Detail       string
	PostalCode   string
	Tel          string
}

func (addr Address) IsEmpty() bool {

	return strings.EqualFold(addr.Name, "") || strings.EqualFold(addr.Tel, "") || strings.EqualFold(addr.Detail, "")
}
