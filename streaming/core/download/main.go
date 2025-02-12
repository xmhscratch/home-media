package download

import (
	"encoding/json"
	"fmt"
	"home-media/streaming/core"
	"home-media/streaming/core/command"
	"home-media/streaming/core/duration"
	"home-media/streaming/core/runtime"
	"home-media/streaming/core/session"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
)

func Start(cfg *core.Config, msg *session.DQMessage) error {
	var exitCode chan int = make(chan int)

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

		reader.WriteVar("ExecBin", filepath.Join(cfg.RootPath, "./ci/download.sh"))
		reader.WriteVar("DownloadURL", msg.DownloadURL)
		reader.WriteVar("Output", msg.SavePath)
		reader.WriteVar("BaseURL", cfg.StreamApiURL)
		reader.WriteVar("RootDir", filepath.Join(cfg.RootPath, cfg.DataDir))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func PeriodicPushHandler(cfg *core.Config, rds *redis.Client) core.PeriodicPushFunc {
	return func(queue map[string]interface{}) (interface{}, string, error) {
		var (
			err   error
			qItem string
			dm    *session.DQMessage
		)

		if qItem, err = rds.SPop(
			core.SessionContext,
			session.GetKeyName("download", ":queue"),
		).Result(); err != nil {
			return nil, "", err
		} else {
			if err = json.Unmarshal([]byte(qItem), &dm); err != nil {
				return nil, "", err
			}
		}
		return dm, dm.SavePath, err
	}
}

func OnPushedHandler(cfg *core.Config, rds *redis.Client) core.OnPushedFunc {
	return func(item interface{}, key string) {
		var dm *session.DQMessage = item.(*session.DQMessage)
		// none blocking download
		go func() {
			Start(cfg, dm)
			duration.Update(cfg, dm)
		}()
		// litter.D("item pushed", dm)
	}
}

func PeriodicRemoveHandler(cfg *core.Config, rds *redis.Client) core.PeriodicRemoveFunc {
	return func(queue map[string]interface{}) (string, error) {
		var (
			err   error
			qItem string
			dm    *session.DQMessage
		)

		if qItem, err = rds.SPop(
			core.SessionContext,
			session.GetKeyName("download", ":done"),
		).Result(); err != nil {
			return "", err
		} else {
			if err = json.Unmarshal([]byte(qItem), &dm); err != nil {
				return "", err
			}
		}

		return dm.SavePath, err
	}
}

func OnRemovedHandler(cfg *core.Config, rds *redis.Client) core.OnRemovedFunc {
	return func(item interface{}, keyId string) {
		dm := item.(*session.DQMessage)
		go func() {
			sm := session.SQSegmentInfo{DQMessage: dm}
			sm.Init(cfg)

			if err := sm.PushQueue(); err != nil {
				fmt.Println(err)
			}
		}()
	}
}
