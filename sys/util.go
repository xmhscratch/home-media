package sys

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Random(a int, z int) int {
	var (
		min int = a
		max int = z
	)
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	return rng.Intn(max-min+1) + min
}

func CheckFileExists(filePath string) bool {
	if f, err := os.Open(filePath); err != nil {
		return false
	} else {
		f.Close()
	}
	return true
}

func WaitTermination() {
	exit := make(chan struct{})
	SignalC := make(chan os.Signal, 4)

	signal.Notify(
		SignalC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for s := range SignalC {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
				return
			}
		}
	}()

	<-exit
	os.Exit(0)
}
