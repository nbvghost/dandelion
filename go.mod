module github.com/nbvghost/dandelion

go 1.13

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/validator/v10 v10.9.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang/protobuf v1.5.2
	github.com/nbvghost/gpa v0.0.0-20210616142117-afb9b836a1c4
	github.com/nbvghost/gweb v1.2.16
	github.com/nbvghost/tool v0.0.0-20210205100218-d99aeb6cf016
	github.com/stretchr/testify v1.7.0 // indirect
	go.etcd.io/etcd/client/v3 v3.5.1
	golang.org/x/net v0.0.0-20211015174653-db2dff38ab41
	golang.org/x/sys v0.0.0-20211013075003-97ac67df715c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211015135405-485ec31e706e // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/nbvghost/gweb => ../gweb

replace github.com/nbvghost/gpa => ../gpa

//replace github.com/nbvghost/gweb v1.2.13 => C:\Users\nbvghost\Desktop\gweb
//replace github.com/nbvghost/gweb v1.2.14 => /Users/nbvghost/Desktop/gweb
//replace github.com/nbvghost/gweb v1.2.14 => /home/nbvghost/datas/projects/gweb
