package metadata

import (
	"home-media/sys"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"home-media/sys/session"
	"os"
	"path/filepath"
)

func UpdateDuration(cfg *sys.Config, msg *session.DQMessage) error {
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

		reader.WriteVar("ExecBin", filepath.Join(cfg.RootPath, "./ci/duration.sh"))
		reader.WriteVar("Input", filepath.Join(cfg.RootPath, cfg.DataDir, msg.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func UpdateSubtitles(cfg *sys.Config, msg *session.DQMessage) error {
	var exitCode chan int = make(chan int)

	reader := command.NewCommandReader()
	writer := command.NewSessionWriter(cfg, msg.SessionId, "subtitles")

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  reader,
			Stdout: writer,
			Stderr: os.Stderr,

			Args: os.Args,

			Main: Main,
		}

		reader.WriteVar("ExecBin", filepath.Join(cfg.RootPath, "./ci/subtitle.sh"))
		reader.WriteVar("Input", filepath.Join(cfg.RootPath, cfg.DataDir, msg.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func UpdateDubs(cfg *sys.Config, msg *session.DQMessage) error {
	var exitCode chan int = make(chan int)

	reader := command.NewCommandReader()
	writer := command.NewSessionWriter(cfg, msg.SessionId, "dubs")

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  reader,
			Stdout: writer,
			Stderr: os.Stderr,

			Args: os.Args,

			Main: Main,
		}

		reader.WriteVar("ExecBin", filepath.Join(cfg.RootPath, "./ci/dub.sh"))
		reader.WriteVar("Input", filepath.Join(cfg.RootPath, cfg.DataDir, msg.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}
