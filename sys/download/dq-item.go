package download

import (
	"fmt"
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
	return ctx.dm.FileKey
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
	close(exitCode)

	litter.D("file downloaded!")
	return nil
}

func (ctx *DQItem) UpdateDuration() error {
	litter.D("update duration...")

	var exitCode chan int = make(chan int)

	stdin := command.NewCommandReader()
	stdout := command.NewSessionFileWriter(ctx.cfg, ctx.dm.FileMeta, ctx.dm.SessionId, ctx.dm.FileKey, "duration")
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
	close(exitCode)

	litter.D("duration updated!")
	return nil
}

func (ctx *DQItem) UpdateSubtitles() error {
	litter.D("update subtitle...")

	var exitCode chan int = make(chan int)

	stdin := command.NewCommandReader()
	stdout := command.NewSessionFileWriter(ctx.cfg, ctx.dm.FileMeta, ctx.dm.SessionId, ctx.dm.FileKey, "subtitles")
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
	close(exitCode)

	litter.D("subtitle updated!")
	return nil
}

func (ctx *DQItem) UpdateDubs() error {
	litter.D("update dub...")

	var exitCode chan int = make(chan int)

	stdin := command.NewCommandReader()
	stdout := command.NewSessionFileWriter(ctx.cfg, ctx.dm.FileMeta, ctx.dm.SessionId, ctx.dm.FileKey, "dubs")
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
	close(exitCode)

	litter.D("dub updated!")
	return nil
}

func (ctx *DQItem) UpdateSourceReady(isReady bool) error {
	litter.D("update source ready...")

	var err error
	writer := command.NewSessionFileWriter(
		ctx.cfg,
		ctx.dm.FileMeta,
		ctx.dm.SessionId,
		ctx.dm.FileKey,
		"sourceReady",
	)
	_, err = writer.Write([]byte(map[bool]string{true: "1", false: "0"}[isReady]))

	litter.D("source ready updated!")
	return err
}

func (ctx *DQItem) ExtractVideo() error {
	litter.D("extracting video...")

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
	<-exitCode
	close(exitCode)

	litter.D("video extracted!")
	return nil
}

func (ctx *DQItem) ExtractDubs() error {
	litter.D("extracting dub...")

	var err error

	for _, dub := range ctx.dm.FileMeta.Dubs {
		fmt.Println("extract audio: ", dub)

		err = func(exitCode chan int) error {
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
			<-exitCode
			close(exitCode)

			return nil
		}(make(chan int))
	}

	litter.D("dub extracted!")
	return err
}

func (ctx *DQItem) ExtractSubtitles() error {
	litter.D("extracting subtitles...")

	var err error

	for _, sub := range ctx.dm.FileMeta.Subtitles {
		fmt.Println("extract subtitle:", sub)

		err = func(exitCode chan int) error {
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
			<-exitCode
			close(exitCode)

			return nil
		}(make(chan int))
	}

	litter.D("subtitle extracted!")
	return err
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
