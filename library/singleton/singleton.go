package singleton

import (
	"sync"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/service/etcd"
	"github.com/nbvghost/dandelion/service/postgres"
	"github.com/nbvghost/gpa"
)

var once sync.Once
var instance struct {
	pq postgres.IPostgres
}

func Orm() *gorm.DB {
	return instance.pq.Orm()
}
func GPA() gpa.IDataBase {

	return instance.pq.GPA()
}
func Init(etcd etcd.IEtcd, serverName string) {
	once.Do(func() {
		var dns string
		var err error
		if dns, err = etcd.ObtainPostgresql(serverName); err != nil {
			panic(err)
		}
		instance.pq = postgres.New(dns)

	})
}
