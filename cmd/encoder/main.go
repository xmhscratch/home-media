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

	sys.NewQueue(sys.QueueOptions{
		Capacity:       1,
		TickDelay:      1000 * 1,
		LoopDelay:      1,
		PeriodicPush:   segment.PeriodicPushHandler(cfg, rds),
		OnPushed:       segment.OnPushedHandler(cfg, rds),
		PeriodicRemove: segment.PeriodicRemoveHandler(cfg, rds),
		OnRemoved:      segment.OnRemovedHandler(cfg, rds),
	}).Start()
}
