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

	"github.com/redis/go-redis/v9"
	"github.com/sanity-io/litter"
)

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

		stdin.WriteVar("ExecBin", "/bin/home-media/download.sh")
		stdin.WriteVar("DownloadURL", msg.DownloadURL)
		stdin.WriteVar("Output", msg.SavePath)
		stdin.WriteVar("BaseURL", cfg.StreamApiURL)
		stdin.WriteVar("RootDir", filepath.Join(cfg.DataPath))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func PeriodicPushHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicPushFunc {
	return func(queue map[string]interface{}) (interface{}, string, error) {
		var (
			err   error
			qItem string
			dm    *session.DQMessage
		)

		if qItem, err = rds.SPop(
			sys.SessionContext,
			session.GetKeyName("download", ":queue"),
		).Result(); err != nil {
			return nil, "", err
		} else {
			if err = json.Unmarshal([]byte(qItem), &dm); err != nil {
				return nil, "", err
			}
		}

		// litter.D(dm, session.GetFileKeyName(dm.SavePath))
		return dm, session.GetFileKeyName(dm.SavePath), err
	}
}

func OnPushedHandler(cfg *sys.Config, rds *redis.Client) sys.OnPushedFunc {
	return func(item interface{}, keyId string) {
		var dm *session.DQMessage = item.(*session.DQMessage)

		Start(cfg, dm)

		metadata.UpdateDuration(cfg, dm)
		metadata.UpdateSubtitles(cfg, dm)
		metadata.UpdateDubs(cfg, dm)

		litter.D("item pushed", keyId)
	}
}

func PeriodicRemoveHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicRemoveFunc {
	return func(queue map[string]interface{}) (string, error) {
		var (
			err   error
			qItem string
			dm    *session.DQMessage
		)

		if qItem, err = rds.SPop(
			sys.SessionContext,
			session.GetKeyName("download", ":done"),
		).Result(); err != nil {
			return "", err
		} else {
			if err = json.Unmarshal([]byte(qItem), &dm); err != nil {
				return "", err
			}
		}

		// litter.D("item removing", session.GetFileKeyName(dm.SavePath))
		return session.GetFileKeyName(dm.SavePath), err
	}
}

func OnRemovedHandler(cfg *sys.Config, rds *redis.Client) sys.OnRemovedFunc {
	return func(item interface{}, keyId string) {
		litter.D("item removed", keyId)

		dm := item.(*session.DQMessage)

		sm := session.SQSegmentInfo{DQMessage: dm}
		sm.Init(cfg)
		// litter.D(dm, sm)

		if err := sm.PushQueue(); err != nil {
			litter.D(err)
		}

		if err := extract.ExtractSubtitles(cfg, dm); err != nil {
			litter.D(err)
		}

		if err := extract.ExtractDubs(cfg, dm); err != nil {
			litter.D(err)
		}
	}
}
