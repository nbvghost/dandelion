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

func Connect(etcd constrain.IEtcd, dbName string) error {
	d, err := NewDB(etcd, dbName)
	if err != nil {
		return err
	}
	instance.pq = d
	return nil
}

func NewDB(etcd constrain.IEtcd, dbName string) (postgres.IPostgres, error) {
	var err error
	var dns string
	if dns, err = etcd.ObtainPostgresql(dbName); err != nil {
		return nil, err
	}
	pq := postgres.New(dns)
	return pq, nil
}
