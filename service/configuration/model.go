package configuration

type ActonType string

const (
	ActonTypeWebview ActonType = "WEBVIEW"
	ActonTypePage    ActonType = "PAGE"
)

type Image struct {
	Src       string
	Url       string
	ActonType ActonType
	Title     string
	Show      bool
}

type Header struct {
	Image
	Style string
}

type QuickLink struct {
	Image
}

type ShowPositionType string

const (
	ShowPositionTypePop    ShowPositionType = "POP"
	ShowPositionTypeBanner ShowPositionType = "BANNER"
)

type Pop struct {
	Matching []string
	Type     ShowPositionType
	Images   []Image
}

type Advert struct {
	Matching []string
	Type     ShowPositionType
	Images   []Image
}
type BrokerageType string

const (
	BrokeragePRODUCT BrokerageType = "PRODUCT"
	BrokerageCUSTOM  BrokerageType = "CUSTOM"
)

type Brokerage struct {
	Type  BrokerageType //PRODUCT,CUSTOM
	Leve1 uint
	Leve2 uint
	Leve3 uint
	Leve4 uint
	Leve5 uint
	Leve6 uint
}
