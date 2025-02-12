package main

import (
	"home-media/streaming/core"
	"home-media/streaming/core/concat"
	"log"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			return
		}
	}()

	cfg, err := core.NewConfig("../")
	if err != nil {
		panic(err)
	}

	rds := core.NewClient(cfg)
	defer rds.Close()

	core.NewQueue(core.QueueOptions{
		Capacity:       1,
		TickDelay:      1000 * 1,
		LoopDelay:      1,
		PeriodicPush:   concat.PeriodicPushHandler(cfg, rds),
		OnPushed:       concat.OnPushedHandler(cfg, rds),
		PeriodicRemove: concat.PeriodicRemoveHandler(cfg, rds),
		OnRemoved:      concat.OnRemovedHandler(cfg, rds),
	}).Start()
}
