package download

import (
	"context"
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
	sys.QItem[DQItem]
	dm *session.DQMessage
}

func (ctx DQItem) Index() int {
	return int(time.Now().Unix())
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

		stdin.WriteVar("ExecBin", filepath.Join(cfg.BinPath, "./download.sh"))
		stdin.WriteVar("DownloadURL", msg.DownloadURL)
		stdin.WriteVar("Output", msg.SavePath)
		stdin.WriteVar("BaseURL", cfg.StreamApiURL)
		stdin.WriteVar("RootDir", filepath.Join(cfg.DataPath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func PeriodicHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicFunc[DQItem] {
	return func(queue *sys.QueueStack[DQItem]) (*DQItem, error) {
		// return &DQItem{dm: nil}, nil
		var (
			err    error
			qItem  string
			dm     *session.DQMessage
			hasKey int64 = 0
		)

		rdsKeyName := session.GetKeyName("download", ":queue")
		if hasKey, err = rds.Exists(context.TODO(), rdsKeyName).Result(); err != nil || hasKey == 0 {
			return nil, nil
		}
		// litter.D(rdsKeyName, hasKey)
		if qItem, err = rds.SPop(
			context.TODO(),
			rdsKeyName,
		).Result(); err != nil {
			// litter.D(qItem, err)
			return nil, nil
		} else {
			if err = json.Unmarshal([]byte(qItem), &dm); err != nil {
				return &DQItem{dm: dm}, err
			}
		}

		// litter.D(dm)
		return &DQItem{dm: dm}, err
	}
}

func ConsumeHandler(cfg *sys.Config, rds *redis.Client) sys.ConsumeFunc[DQItem] {
	return func(queue *sys.QueueStack[DQItem], item *DQItem) error {
		var err error
		// litter.D(item)
		if err = Start(cfg, item.dm); err != nil {
			return err
		}
		if err = metadata.UpdateDuration(cfg, item.dm); err != nil {
			return err
		}
		if err = metadata.UpdateSubtitles(cfg, item.dm); err != nil {
			return err
		}
		if err = metadata.UpdateDubs(cfg, item.dm); err != nil {
			return err
		}

		// litter.D("item pushed", item.Key())
		return err
	}
}

func OnConsumedHandler(cfg *sys.Config, rds *redis.Client) sys.OnConsumedFunc[DQItem] {
	return func(item *DQItem) {
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
		// =======================
		// litter.D("item removed", item.Key())

		sm := session.SQSegmentInfo{DQMessage: item.dm}
		sm.Init(cfg)
		// litter.D(item.dm, sm)

		if err := sm.PushQueue(); err != nil {
			litter.D(err)
		}

		if err := extract.ExtractVideo(cfg, item.dm); err != nil {
			litter.D(err)
		}

		if err := extract.ExtractDubs(cfg, item.dm); err != nil {
			litter.D(err)
		}

		if err := extract.ExtractSubtitles(cfg, item.dm); err != nil {
			litter.D(err)
		}
	}
}
