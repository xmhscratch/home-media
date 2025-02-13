package duration

import (
	"home-media/sys"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"home-media/sys/session"
	"os"
	"path/filepath"
)

func Update(cfg *sys.Config, msg *session.DQMessage) error {
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
