package postgres

import (
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	_ "github.com/lib/pq"
	"github.com/nbvghost/glog"
	"github.com/nbvghost/gpa"
)

type server struct {
	dns     string
	ormOnce sync.Once
	gpaOnce sync.Once

	orm *gorm.DB
	gpa gpa.IDataBase

	logger logger.Interface
}

func (p *server) Orm() *gorm.DB {
	p.ormOnce.Do(func() {
		_database, err := gorm.Open(postgres.Open(p.dns), &gorm.Config{
			Logger: p.logger,
		})
		if err != nil {
			glog.Error(err)
		}
		p.orm = _database
	})
	return p.orm
}

func (p *server) GPA() gpa.IDataBase {
	p.gpaOnce.Do(func() {
		p.gpa = gpa.Open("postgres", p.dns, nil)
	})

	return p.gpa
}

func New(dns string) IPostgres {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	return &server{dns: dns, logger: newLogger}
}
