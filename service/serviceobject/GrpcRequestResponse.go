package serviceobject

type GrpcRequest struct {
	Route      string
	HttpMethod string
	Body       map[string]interface{}
}
type GrpcResponse struct {
	Error int
	Data  interface{}
}
