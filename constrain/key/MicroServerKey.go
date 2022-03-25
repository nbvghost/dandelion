package key

type MicroServerKey string

const (
	MicroServerKeySSO     MicroServerKey = "sso"
	MicroServerKeyOSS     MicroServerKey = "oss"
	MicroServerKeyADMIN   MicroServerKey = "dandelion.admin"
	MicroServerKeySITE    MicroServerKey = "dandelion.site"
	MicroServerKeyMANAGER MicroServerKey = "dandelion.manager"
)
