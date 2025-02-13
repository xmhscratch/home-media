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

	sys.NewQueue(sys.QueueOptions{
		Capacity:       1,
		TickDelay:      1000 * 1,
		LoopDelay:      1,
		PeriodicPush:   concat.PeriodicPushHandler(cfg, rds),
		OnPushed:       concat.OnPushedHandler(cfg, rds),
		PeriodicRemove: concat.PeriodicRemoveHandler(cfg, rds),
		OnRemoved:      concat.OnRemovedHandler(cfg, rds),
	}).Start()
}
