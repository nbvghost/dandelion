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
	"github.com/nbvghost/dandelion/library/environments"
	"github.com/nbvghost/gpa"
)

type server struct {
	dsn     string
	ormOnce sync.Once
	gpaOnce sync.Once

	orm *gorm.DB
	gpa gpa.IDataBase

	logger logger.Interface
}

func (p *server) Orm() *gorm.DB {
	p.ormOnce.Do(func() {
		_database, err := gorm.Open(postgres.Open(p.dsn), &gorm.Config{
			Logger: p.logger,
		})
		if err != nil {
			log.Println(err)
		}
		p.orm = _database
	})
	return p.orm
}

func (p *server) GPA() gpa.IDataBase {
	p.gpaOnce.Do(func() {
		p.gpa = gpa.Open("postgres", p.dsn, nil)
	})

	return p.gpa
}

func New(dsn string) IPostgres {
	var newLogger logger.Interface
	if environments.Release() {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Millisecond * 500, // Slow SQL threshold
				LogLevel:                  logger.Warn,            // Log level
				IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                   // Disable color
			},
		)
	} else {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Millisecond * 100, // Slow SQL threshold
				LogLevel:                  logger.Info,            // Log level
				IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                   // Disable color
			},
		)
	}
	return &server{dsn: dsn, logger: newLogger}
}
