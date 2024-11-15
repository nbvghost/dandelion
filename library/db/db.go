package db

import (
	"errors"
	"github.com/nbvghost/dandelion/constrain"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/server/postgres"
)

var instance struct {
	pq postgres.IPostgres
}

func Orm() *gorm.DB {
	if instance.pq == nil {
		panic(errors.New("请初化数据库"))
	}
	return instance.pq.Orm()
}
func ConnectWithout(dsn string) error {
	pq := postgres.New(dsn)
	instance.pq = pq
	return nil
}
func Connect(etcd constrain.IEtcd, dbName string) error {
	var err error
	var dsn string
	if dsn, err = etcd.ObtainPostgresql(dbName); err != nil {
		return err
	}
	pq := postgres.New(dsn)
	instance.pq = pq
	return nil
}
