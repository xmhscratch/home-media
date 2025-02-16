package extract

import (
	"home-media/sys"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"home-media/sys/session"
	"os"
	"path/filepath"
)

func ExtractVideo(cfg *sys.Config, msg *session.DQMessage) error {
	func(exitCode chan int) int {
		defer close(exitCode)

		go func() {
			stdin := command.NewCommandReader()
			stdout := command.NewNullWriter()
			stderr := command.NewNullWriter()

			shell := runtime.Shell{
				PID: os.Getpid(),

				Stdin:  stdin,
				Stdout: stdout,
				Stderr: stderr,

				Args: os.Args,

				Main: Main,
			}

			stdin.WriteVar("ExecBin", filepath.Join(cfg.BinPath, "./extract-vid.sh"))
			stdin.WriteVar("Input", filepath.Join(cfg.DataPath, msg.SavePath))
			stdin.WriteVar("StreamIndex", "0")
			stdin.WriteVar("LangCode", "default")
			stdin.WriteVar("Output", filepath.Join(
				cfg.DataPath,
				filepath.Dir(msg.SavePath),
				session.GetFileKeyName(msg.SavePath),
			))

			exitCode <- shell.Run()
		}()
		return <-exitCode
	}(make(chan int))

	return nil
}
