package main

import (
	"home-media/sys"
	"home-media/sys/segment"
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

	sys.NewQueue(sys.QueueOptions[segment.SQItem]{
		Capacity:   1,
		Throttle:   500,
		Periodic:   segment.PeriodicHandler(cfg, rds),
		Consume:    segment.ConsumeHandler(cfg, rds),
		OnConsumed: segment.OnConsumedHandler(cfg, rds),
	}).Start()
}
