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
	Style   string
	Version string
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
type BaiduTranslateConfiguration struct {
	URL         string
	SecurityKey string
	Appid       string
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
	Leve1 float64
	Leve2 float64
	Leve3 float64
	Leve4 float64
	Leve5 float64
	Leve6 float64
}
