package main

import (
	"fmt"
	"home-media/sys/tui"
	"os"
)

func main() {
	// tui.NewTest()
	var (
		err error
		m   *tui.TuiManager
	)

	m, err = tui.NewTuiManager()
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
// echo "4|$(cat ./cmd/tui/sample.txt)" | socat - unix:///run/tuid.sock
