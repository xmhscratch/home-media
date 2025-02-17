package main

import (
	"home-media/sys"
	"home-media/sys/download"
	"log"

	"github.com/sanity-io/litter"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			return
		}
	}()

	cfg, err := sys.NewConfig("./")
	// cfg, err := sys.NewConfig("../")
	if err != nil {
		panic(err)
	}

	rds := sys.NewClient(cfg)
	defer rds.Close()

	sys.NewQueue(sys.QueueOptions[download.DQItem]{
		Capacity:   2,
		Throttle:   500,
		Periodic:   download.PeriodicHandler(cfg, rds),
		Consume:    download.ConsumeHandler(cfg, rds),
		OnConsumed: download.OnConsumedHandler(cfg, rds),
		OnError: func(err error) {
			litter.D(err)
		},
	}).Start()
}
