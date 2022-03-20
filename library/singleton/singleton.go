package singleton

import (
	"github.com/nbvghost/dandelion/constrain"
	"sync"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/server/postgres"
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
func Init(etcd constrain.IEtcd, serverName string) error {
	var err error
	once.Do(func() {
		var dns string

		if dns, err = etcd.ObtainPostgresql(serverName); err != nil {
			panic(err)
		}
		instance.pq = postgres.New(dns)

		gpa.Init(func() gpa.IDataBase {
			return gpa.Open("postgres", dns, nil)
		})

	})

	return err
}
