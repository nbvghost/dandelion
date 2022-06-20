package key

type serverType string

const (
	ServerTypeHttp serverType = "http"
	ServerTypeGrpc serverType = "grpc"
)

type MicroServer struct {
	Name       string
	ServerType serverType
}

var (
	MicroServerSSO        = MicroServer{Name: "sso", ServerType: ServerTypeHttp}
	MicroServerOSS        = MicroServer{Name: "oss", ServerType: ServerTypeHttp}
	MicroServerADMIN      = MicroServer{Name: "dandelion.admin", ServerType: ServerTypeHttp}
	MicroServerSITE       = MicroServer{Name: "dandelion.site", ServerType: ServerTypeHttp}
	MicroServersShop      = MicroServer{Name: "dandelion.shop", ServerType: ServerTypeHttp}
	MicroServerMANAGER    = MicroServer{Name: "dandelion.manager", ServerType: ServerTypeHttp}
	MicroServerMimiServer = MicroServer{Name: "mimi.server", ServerType: ServerTypeHttp}
)
