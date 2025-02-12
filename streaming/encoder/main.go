package main

import (
	"home-media/streaming/core"
	"home-media/streaming/core/segment"
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
		PeriodicPush:   segment.PeriodicPushHandler(cfg, rds),
		OnPushed:       segment.OnPushedHandler(cfg, rds),
		PeriodicRemove: segment.PeriodicRemoveHandler(cfg, rds),
		OnRemoved:      segment.OnRemovedHandler(cfg, rds),
	}).Start()
}
