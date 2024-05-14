package task

import (
	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/service/internal/task/job"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

func Start(context constrain.IWithoutSessionContext) {
	group := errgroup.Group{}
	group.Go(func() error {
		t := job.NewQueryOrdersTask(context)
		run := func() {
			log.Println("NewQueryOrdersTask[* 5 0]")
			err := t.Run()
			if err != nil {
				log.Println(err)
			}
		}
		run()
		ticker := time.NewTicker(time.Minute * 5)
		for range ticker.C {
			run()
		}
		return nil
	})
	group.Go(func() error {
		t := job.NewQueryExpressCompanyTask(context)
		log.Println("NewQueryExpressCompanyTask")
		err := t.Run()
		if err != nil {
			return err
		}
		ticker := time.NewTicker(time.Hour * 24)
		for range ticker.C {
			log.Println("NewQueryExpressCompanyTask[24 0 0]")
			err := t.Run()
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})

	group.Go(func() error {
		t := job.NewQuerySupplyOrdersTask(context)
		log.Println("NewQuerySupplyOrdersTask")
		err := t.Run()
		if err != nil {
			return err
		}
		ticker := time.NewTicker(time.Second * 60)
		for range ticker.C {
			log.Println("NewQuerySupplyOrdersTask[0 0 60]")
			err = t.Run()
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})

	group.Go(func() error {
		t := job.NewQueryTransfersTask(context)
		log.Println("NewQueryTransfersTask")
		err := t.Run()
		if err != nil {
			return err
		}
		ticker := time.NewTicker(time.Second * 15)
		for range ticker.C {
			log.Println("NewQueryTransfersTask[0 0 15]")
			err = t.Run()
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
	err := group.Wait()
	if err != nil {
		log.Println(err)
	}
}
