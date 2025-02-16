package main

import (
	"home-media/sys"
	"home-media/sys/concat"
	"log"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			return
		}
	}()

	cfg, err := sys.NewConfig("../")
	if err != nil {
		panic(err)
	}

	rds := sys.NewClient(cfg)
	defer rds.Close()

	sys.NewQueue(sys.QueueOptions[concat.CQItem]{
		Capacity:   2,
		Throttle:   500,
		Periodic:   concat.PeriodicHandler(cfg, rds),
		Consume:    concat.ConsumeHandler(cfg, rds),
		OnConsumed: concat.OnConsumedHandler(cfg, rds),
	}).Start()
}
