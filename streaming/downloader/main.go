package main

import (
	"home-media/streaming/core"
	"home-media/streaming/core/download"
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
		Capacity:       2,
		TickDelay:      1000 * 1,
		LoopDelay:      1,
		PeriodicPush:   download.PeriodicPushHandler(cfg, rds),
		OnPushed:       download.OnPushedHandler(cfg, rds),
		PeriodicRemove: download.PeriodicRemoveHandler(cfg, rds),
		OnRemoved:      download.OnRemovedHandler(cfg, rds),
	}).Start()
}
