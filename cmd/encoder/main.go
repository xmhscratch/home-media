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
		Capacity:     1,
		LoopDelay:    500,
		PeriodicPush: segment.PeriodicPushHandler(cfg, rds),
		OnPushed:     segment.OnPushedHandler(cfg, rds),
		OnRemoved:    segment.OnRemovedHandler(cfg, rds),
	}).Start()
}
