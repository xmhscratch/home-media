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

	stdin := command.NewCommandReader()
	stdout := command.NewSessionWriter(cfg, msg.SessionId, "duration")
	stderr := command.NewNullWriter()

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,

			Args: os.Args,

			Main: Main,
		}

		stdin.WriteVar("ExecBin", "/bin/home-media/duration.sh")
		stdin.WriteVar("Input", filepath.Join(cfg.DataPath, msg.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func UpdateSubtitles(cfg *sys.Config, msg *session.DQMessage) error {
	var exitCode chan int = make(chan int)

	stdin := command.NewCommandReader()
	stdout := command.NewSessionWriter(cfg, msg.SessionId, "subtitles")
	stderr := command.NewNullWriter()

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,

			Args: os.Args,

			Main: Main,
		}

		stdin.WriteVar("ExecBin", "/bin/home-media/subtitle.sh")
		stdin.WriteVar("Input", filepath.Join(cfg.DataPath, msg.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func UpdateDubs(cfg *sys.Config, msg *session.DQMessage) error {
	var exitCode chan int = make(chan int)

	stdin := command.NewCommandReader()
	stdout := command.NewSessionWriter(cfg, msg.SessionId, "dubs")
	stderr := command.NewNullWriter()

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,

			Args: os.Args,

			Main: Main,
		}

		stdin.WriteVar("ExecBin", "/bin/home-media/dub.sh")
		stdin.WriteVar("Input", filepath.Join(cfg.DataPath, msg.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}
