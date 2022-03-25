package serviceobject

type ServerDesc struct {
	Name string
	Port int
	IP   string
}

func NewServerDesc(name string, port int, ip string) *ServerDesc {
	return &ServerDesc{
		Name: name,
		Port: port,
		IP:   ip,
	}
}
