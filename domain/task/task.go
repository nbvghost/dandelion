package task

import (
	"context"
	"github.com/nbvghost/dandelion/domain/task/job"
	"log"
	"time"
)

func Start() {
	go func() {
		t := job.NewQueryOrdersTask(context.TODO())
		ticker := time.NewTicker(time.Minute * 5)
		for range ticker.C {
			log.Println("NewQueryOrdersTask[* 5 0]")
			err := t.Run()
			if err != nil {
				log.Println(err)
			}
		}
	}()
	go func() {
		t := job.NewQueryExpressCompanyTask(context.TODO())
		log.Println("NewQueryExpressCompanyTask")
		err := t.Run()
		if err != nil {
			log.Println(err)
		}
		ticker := time.NewTicker(time.Hour * 24)
		for range ticker.C {
			log.Println("NewQueryExpressCompanyTask[24 0 0]")
			err := t.Run()
			if err != nil {
				log.Println(err)
			}
		}
	}()
}
