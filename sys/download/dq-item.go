package download

import (
	"home-media/sys"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"home-media/sys/session"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sanity-io/litter"
)

func (ctx DQItem) Index() int {
	return int(time.Now().Unix())
}

func (ctx DQItem) Key() string {
	return session.GetFileKeyName(ctx.dm.SavePath)
}

func (ctx *DQItem) StartDownload() error {
	litter.D("start downloading:", ctx.dm)

	var exitCode chan int = make(chan int)

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

			Main: DownloadShell,
		}

		stdin.WriteVar("ExecBin", filepath.Join(ctx.cfg.BinPath, "./download.sh"))
		stdin.WriteVar("DownloadURL", ctx.dm.DownloadURL)
		stdin.WriteVar("Output", ctx.dm.SavePath)
		stdin.WriteVar("BaseURL", ctx.cfg.StreamApiURL)
		stdin.WriteVar("RootDir", filepath.Join(ctx.cfg.DataPath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func (ctx *DQItem) UpdateDuration() error {
	litter.D("update duration...")

	var exitCode chan int = make(chan int)

	stdin := command.NewCommandReader()
	stdout := command.NewSessionInfoWriter(ctx.cfg, ctx.dm.SessionId, "duration")
	stderr := command.NewNullWriter()

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,

			Args: os.Args,

			Main: MetadataShell,
		}

		stdin.WriteVar("ExecBin", filepath.Join(ctx.cfg.BinPath, "./duration.sh"))
		stdin.WriteVar("Input", filepath.Join(ctx.cfg.DataPath, ctx.dm.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func (ctx *DQItem) UpdateSubtitles() error {
	litter.D("update subtitle...")

	var exitCode chan int = make(chan int)

	stdin := command.NewCommandReader()
	stdout := command.NewSessionFileWriter(ctx.cfg, ctx.dm.SessionId, ctx.Key(), "subtitles")
	stderr := command.NewNullWriter()

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,

			Args: os.Args,

			Main: MetadataShell,
		}

		stdin.WriteVar("ExecBin", filepath.Join(ctx.cfg.BinPath, "./subtitle.sh"))
		stdin.WriteVar("Input", filepath.Join(ctx.cfg.DataPath, ctx.dm.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func (ctx *DQItem) UpdateDubs() error {
	litter.D("update dub...")

	var exitCode chan int = make(chan int)

	stdin := command.NewCommandReader()
	stdout := command.NewSessionFileWriter(ctx.cfg, ctx.dm.SessionId, ctx.Key(), "dubs")
	stderr := command.NewNullWriter()

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,

			Args: os.Args,

			Main: MetadataShell,
		}

		stdin.WriteVar("ExecBin", filepath.Join(ctx.cfg.BinPath, "./dub.sh"))
		stdin.WriteVar("Input", filepath.Join(ctx.cfg.DataPath, ctx.dm.SavePath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func (ctx *DQItem) ExtractVideo() error {
	litter.D("extract video...")

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

				Main: ExtractShell,
			}

			stdin.WriteVar("ExecBin", filepath.Join(ctx.cfg.BinPath, "./extract-vid.sh"))
			stdin.WriteVar("Input", filepath.Join(ctx.cfg.DataPath, ctx.dm.SavePath))
			stdin.WriteVar("StreamIndex", "0")
			stdin.WriteVar("LangCode", "default")
			stdin.WriteVar("Output", filepath.Join(
				ctx.cfg.DataPath,
				filepath.Dir(ctx.dm.SavePath),
				session.GetFileKeyName(ctx.dm.SavePath),
			))

			exitCode <- shell.Run()
		}()
		return <-exitCode
	}(make(chan int))

	return nil
}

func (ctx *DQItem) ExtractDubs() error {
	litter.D("extract audio...")

	for _, dub := range ctx.dm.FileMeta.Dubs {
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

					Main: ExtractShell,
				}

				stdin.WriteVar("ExecBin", filepath.Join(ctx.cfg.BinPath, "./extract-dub.sh"))
				stdin.WriteVar("Input", filepath.Join(ctx.cfg.DataPath, ctx.dm.SavePath))
				stdin.WriteVar("StreamIndex", strconv.FormatInt(dub.StreamIndex, 5<<1))
				stdin.WriteVar("LangCode", sys.BuildString(dub.LangCode, strconv.FormatInt(dub.StreamIndex, 5<<1)))
				stdin.WriteVar("Output", filepath.Join(
					ctx.cfg.DataPath,
					filepath.Dir(ctx.dm.SavePath),
					session.GetFileKeyName(ctx.dm.SavePath),
				))

				exitCode <- shell.Run()
			}()
			return <-exitCode
		}(make(chan int))
	}

	return nil
}

func (ctx *DQItem) ExtractSubtitles() error {
	litter.D("extract subtitle...")

	for _, sub := range ctx.dm.FileMeta.Subtitles {
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

					Main: ExtractShell,
				}

				stdin.WriteVar("ExecBin", filepath.Join(ctx.cfg.BinPath, "./extract-sub.sh"))
				stdin.WriteVar("Input", filepath.Join(ctx.cfg.DataPath, ctx.dm.SavePath))
				stdin.WriteVar("StreamIndex", strconv.FormatInt(sub.StreamIndex, 5<<1))
				stdin.WriteVar("LangCode", sys.BuildString(sub.LangCode, strconv.FormatInt(sub.StreamIndex, 5<<1)))
				stdin.WriteVar("Output", filepath.Join(
					ctx.cfg.DataPath,
					filepath.Dir(ctx.dm.SavePath),
					session.GetFileKeyName(ctx.dm.SavePath),
				))

				exitCode <- shell.Run()
			}()
			return <-exitCode
		}(make(chan int))
	}

	return nil
}

// func (ctx *DQItem) Complete(cfg *sys.Config) {
// 	rds := sys.NewClient(ctx.cfg)
// 	defer rds.Close()

// 	// _, err := rds.HSet(
// 	// 	context.TODO(),
// 	// 	GetKeyName(ctx.SessionId, ":files"),
// 	// 	[]string{ctx.AttrName, strings.TrimSpace(string(p))},
// 	// ).Result()

// 	// return
// }
