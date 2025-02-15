package extract

import (
	"home-media/sys"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"home-media/sys/session"
	"os"
	"path/filepath"
	"strconv"
)

func ExtractDubs(cfg *sys.Config, msg *session.DQMessage) error {
	var (
		err error
		ss  *session.Session[session.FileSourceType]
	)

	if ss, _, err = msg.FileType.InitSession(
		cfg,
		msg.SessionId,
	); err != nil {
		return err
	}

	for _, dub := range ss.Dubs {
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

				stdin.WriteVar("ExecBin", "/export/bin/extract-dub.sh")
				stdin.WriteVar("Input", filepath.Join(cfg.DataPath, msg.SavePath))
				stdin.WriteVar("StreamIndex", strconv.FormatInt(dub.StreamIndex, 5<<1))
				stdin.WriteVar("LangCode", sys.BuildString(dub.LangCode, strconv.FormatInt(dub.StreamIndex, 5<<1)))
				stdin.WriteVar("Output", filepath.Join(
					cfg.DataPath,
					filepath.Dir(msg.SavePath),
					session.GetFileKeyName(msg.SavePath),
				))

				exitCode <- shell.Run()
			}()
			return <-exitCode
		}(make(chan int))
	}

	return nil
}
