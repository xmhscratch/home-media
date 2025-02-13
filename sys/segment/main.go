package segment

import (
	"encoding/json"
	"fmt"
	"home-media/sys"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"home-media/sys/session"
	"os"
	"path/filepath"

	"github.com/redis/go-redis/v9"
	"github.com/sanity-io/litter"
)

func Encode(cfg *sys.Config, sm *session.SQMessage) error {
	var exitCode chan int = make(chan int)

	reader := command.NewCommandReader()

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  reader,
			Stdout: os.Stdout,
			Stderr: os.Stderr,

			Args: os.Args,

			Main: Main,
		}

		reader.WriteVar("ExecBin", filepath.Join(cfg.RootPath, "./ci/segment.sh"))
		reader.WriteVar("Input", sm.Source)
		reader.WriteVar("Start", sm.Start)       //"00:00:00.0000"
		reader.WriteVar("Duration", sm.Duration) //"00:00:03.0000"
		reader.WriteVar("Output", sm.Output)     //"./test_000.mp4"

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func PeriodicPushHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicPushFunc {
	return func(queue map[string]interface{}) (interface{}, string, error) {
		var (
			err   error
			qItem *redis.ZWithKey
			sm    *session.SQMessage
		)

		if qItem, err = rds.BZPopMin(
			sys.SessionContext, 0,
			session.GetKeyName("segment", ":queue"),
		).Result(); err != nil {
			return nil, "", err
		} else {
			err = json.Unmarshal([]byte(qItem.Member.(string)), &sm)
		}

		return sm, sm.KeyId, err
	}
}

func OnPushedHandler(cfg *sys.Config, rds *redis.Client) sys.OnPushedFunc {
	return func(item interface{}, key string) {
		var sm *session.SQMessage = item.(*session.SQMessage)

		// go func() error {
		// 	Encode(cfg, sm)

		// 	return rds.RPush(
		// 		sys.SessionContext,
		// 		session.GetKeyName("segment", ":done"),
		// 		[]string{sm.KeyId, sm.Output},
		// 	).Err()
		// }()

		litter.D("item pushed", sm)
	}
}

func PeriodicRemoveHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicRemoveFunc {
	return func(queue map[string]interface{}) (string, error) {
		var (
			err   error
			qItem []string = make([]string, 2)
			keyId string
		)

		if qItem, err = rds.BLPop(
			sys.SessionContext, 0,
			session.GetKeyName("segment", ":done"),
		).Result(); err != nil {
			return "", err
		}
		keyId = qItem[0]

		return keyId, err
	}
}

func OnRemovedHandler(cfg *sys.Config, rds *redis.Client) sys.OnRemovedFunc {
	return func(item interface{}, keyId string) {
		var sm *session.SQMessage = item.(*session.SQMessage)

		err := rds.SAdd(
			sys.SessionContext,
			session.GetKeyName("concat:ready", ":", keyId),
			sm.Output,
		).Err()

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
