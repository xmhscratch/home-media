package duration

import (
	"home-media/streaming/core"
	"home-media/streaming/core/command"
	"home-media/streaming/core/runtime"
	"home-media/streaming/core/session"
	"os"
	"path/filepath"
)

func Update(cfg *core.Config, msg *session.DQMessage) error {
	var exitCode chan int = make(chan int)

	reader := command.NewCommandReader()
	writer := command.NewSessionWriter(cfg, msg.SessionId, "duration")

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  reader,
			Stdout: writer,
			Stderr: os.Stderr,

			Args: os.Args,

			Main: Main,
		}

		reader.WriteVar("ExecBin", "ffprobe")
		reader.WriteVar("Input", filepath.Join(cfg.RootPath, cfg.DataDir, msg.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}
