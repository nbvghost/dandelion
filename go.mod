module github.com/nbvghost/dandelion

go 1.17

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/go-playground/validator/v10 v10.9.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/nbvghost/glog v1.0.17
	github.com/nbvghost/gpa v0.0.0-20210616142117-afb9b836a1c4
	github.com/nbvghost/gweb v1.2.16
	github.com/nbvghost/tool v0.0.0-20210205100218-d99aeb6cf016
	go.etcd.io/etcd/client/v3 v3.5.1
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	go.etcd.io/etcd/api/v3 v3.5.1 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.1 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20211015174653-db2dff38ab41 // indirect
	golang.org/x/sys v0.0.0-20211013075003-97ac67df715c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211015135405-485ec31e706e // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/nbvghost/gweb => ../../framework/gweb

replace github.com/nbvghost/gpa => ../../framework/gpa

//replace github.com/nbvghost/gweb v1.2.13 => C:\Users\nbvghost\Desktop\gweb
//replace github.com/nbvghost/gweb v1.2.14 => /Users/nbvghost/Desktop/gweb
//replace github.com/nbvghost/gweb v1.2.14 => /home/nbvghost/datas/projects/gweb
