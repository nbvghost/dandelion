module github.com/nbvghost/dandelion

go 1.13

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/nbvghost/gweb v1.2.16
	go.etcd.io/etcd/client/v3 v3.5.1
	golang.org/x/net v0.0.0-20211015174653-db2dff38ab41
	golang.org/x/sys v0.0.0-20211013075003-97ac67df715c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211015135405-485ec31e706e // indirect
	google.golang.org/grpc v1.41.0
)

replace github.com/nbvghost/gweb => ../gweb

//replace github.com/nbvghost/gweb v1.2.13 => C:\Users\nbvghost\Desktop\gweb
//replace github.com/nbvghost/gweb v1.2.14 => /Users/nbvghost/Desktop/gweb
//replace github.com/nbvghost/gweb v1.2.14 => /home/nbvghost/datas/projects/gweb
