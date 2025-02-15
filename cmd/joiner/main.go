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
		Capacity:     1,
		LoopDelay:    500,
		PeriodicPush: concat.PeriodicPushHandler(cfg, rds),
		OnPushed:     concat.OnPushedHandler(cfg, rds),
		OnRemoved:    concat.OnRemovedHandler(cfg, rds),
	}).Start()
}
