package main

import (
	"home-media/sys"
	"home-media/sys/download"
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
		Capacity:       2,
		TickDelay:      1000 * 1,
		LoopDelay:      500,
		PeriodicPush:   download.PeriodicPushHandler(cfg, rds),
		OnPushed:       download.OnPushedHandler(cfg, rds),
		PeriodicRemove: download.PeriodicRemoveHandler(cfg, rds),
		OnRemoved:      download.OnRemovedHandler(cfg, rds),
	}).Start()
}
