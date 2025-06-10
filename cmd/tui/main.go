package main

import (
	"fmt"
	"home-media/sys"
	"home-media/sys/tui"
	"log"
	"os"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			return
		}
	}()

	// tui.NewTest()
	var (
		err error
		m   *tui.TuiManager
	)

	cfg, err := sys.NewConfig("../")
	if err != nil {
		panic(err)
	}

	m, err = tui.NewTuiManager(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go m.ListenToSocket()
	if err = m.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// echo "hello" | socat - unix:///run/tuid.sock
// echo "4"$'\x1E'"$(cat ./cmd/tui/sample.txt)" | socat - unix:///run/tuid.sock
