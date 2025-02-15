package download

import (
	"encoding/json"
	"home-media/sys"
	"home-media/sys/command"
	"home-media/sys/extract"
	"home-media/sys/metadata"
	"home-media/sys/runtime"
	"home-media/sys/session"
	"os"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sanity-io/litter"
)

type DQItem struct {
	sys.QItem
	dm *session.DQMessage
}

func (ctx DQItem) Index() int {
	now := time.Now()
	return int(now.Unix())
}

func (ctx DQItem) Key() string {
	return session.GetFileKeyName(ctx.dm.SavePath)
}

func Start(cfg *sys.Config, msg *session.DQMessage) error {
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

			Main: Main,
		}

		stdin.WriteVar("ExecBin", "/export/bin/download.sh")
		stdin.WriteVar("DownloadURL", msg.DownloadURL)
		stdin.WriteVar("Output", msg.SavePath)
		stdin.WriteVar("BaseURL", cfg.StreamApiURL)
		stdin.WriteVar("RootDir", filepath.Join(cfg.DataPath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func PeriodicPushHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicPushFunc[DQItem] {
	return func(queue *sys.QueueStack[DQItem]) (DQItem, error) {
		var (
			err   error
			qItem string
			dm    *session.DQMessage
		)

		if qItem, err = rds.SPop(
			sys.SessionContext,
			session.GetKeyName("download", ":queue"),
		).Result(); err != nil {
			return DQItem{dm: dm}, err
		} else {
			if err = json.Unmarshal([]byte(qItem), &dm); err != nil {
				return DQItem{dm: dm}, err
			}
		}

		// litter.D(dm, session.GetFileKeyName(dm.SavePath))
		return DQItem{dm: dm}, err
	}
}

func OnPushedHandler(cfg *sys.Config, rds *redis.Client) sys.OnPushedFunc[DQItem] {
	return func(item DQItem) {
		Start(cfg, item.dm)

		metadata.UpdateDuration(cfg, item.dm)
		metadata.UpdateSubtitles(cfg, item.dm)
		metadata.UpdateDubs(cfg, item.dm)

		litter.D("item pushed", item.Key())
	}
}

func OnRemovedHandler(cfg *sys.Config, rds *redis.Client) sys.OnRemovedFunc[DQItem] {
	return func(item DQItem) {
		// var (
		// 	err   error
		// 	qItem string
		// 	dm    *session.DQMessage
		// )

		// if qItem, err = rds.SPop(
		// 	sys.SessionContext,
		// 	session.GetKeyName("download", ":done"),
		// ).Result(); err != nil {
		// 	return
		// } else {
		// 	if err = json.Unmarshal([]byte(qItem), &dm); err != nil {
		// 		return
		// 	}
		// }

		litter.D("item removed", item.Key())

		sm := session.SQSegmentInfo{DQMessage: item.dm}
		sm.Init(cfg)
		// litter.D(dm, sm)

		if err := sm.PushQueue(); err != nil {
			litter.D(err)
		}

		if err := extract.ExtractSubtitles(cfg, item.dm); err != nil {
			litter.D(err)
		}

		if err := extract.ExtractDubs(cfg, item.dm); err != nil {
			litter.D(err)
		}
	}
}
