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

func ExtractSubtitles(cfg *sys.Config, msg *session.DQMessage) error {
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

	for _, sub := range ss.Subtitles {
		func(exitCode chan int) int {
			defer close(exitCode)

			go func() {
				reader := command.NewCommandReader()

				shell := runtime.Shell{
					PID: os.Getpid(),

					Stdin:  reader,
					Stdout: os.Stdout,
					Stderr: os.Stderr,

					Args: os.Args,

					Main: Main,
				}

				reader.WriteVar("ExecBin", filepath.Join(cfg.RootPath, "./ci/extract-sub.sh"))
				reader.WriteVar("Input", filepath.Join(cfg.RootPath, cfg.DataDir, msg.SavePath))
				reader.WriteVar("StreamIndex", strconv.FormatInt(sub.StreamIndex, 5<<1))
				reader.WriteVar("LangCode", sys.BuildString(sub.LangCode, strconv.FormatInt(sub.StreamIndex, 5<<1)))
				reader.WriteVar("Output", filepath.Join(
					cfg.RootPath, cfg.DataDir,
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
				reader := command.NewCommandReader()

				shell := runtime.Shell{
					PID: os.Getpid(),

					Stdin:  reader,
					Stdout: os.Stdout,
					Stderr: os.Stderr,

					Args: os.Args,

					Main: Main,
				}

				reader.WriteVar("ExecBin", filepath.Join(cfg.RootPath, "./ci/extract-dub.sh"))
				reader.WriteVar("Input", filepath.Join(cfg.RootPath, cfg.DataDir, msg.SavePath))
				reader.WriteVar("StreamIndex", strconv.FormatInt(dub.StreamIndex, 5<<1))
				reader.WriteVar("LangCode", sys.BuildString(dub.LangCode, strconv.FormatInt(dub.StreamIndex, 5<<1)))
				reader.WriteVar("Output", filepath.Join(
					cfg.RootPath, cfg.DataDir,
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
