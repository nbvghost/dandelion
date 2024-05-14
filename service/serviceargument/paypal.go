package serviceargument

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
