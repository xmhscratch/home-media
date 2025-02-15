package segment

import (
	"encoding/json"
	"fmt"
	"home-media/sys"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"home-media/sys/session"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sanity-io/litter"
)

type SQItem struct {
	sys.QItem
	sm *session.SQMessage
}

func (ctx SQItem) Index() int {
	now := time.Now()
	return int(now.Unix())
}

func (ctx SQItem) Key() string {
	return ctx.sm.KeyId
}

func Encode(cfg *sys.Config, sm *session.SQMessage) error {
	var exitCode chan int = make(chan int)

	stdin := command.NewCommandReader()
	stdout := command.NewNullWriter()
	stderr := command.NewNullWriter()

	go func() {
		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,

			Args: os.Args,

			Main: Main,
		}

		stdin.WriteVar("ExecBin", "/export/bin/segment.sh")
		stdin.WriteVar("Input", sm.Source)
		stdin.WriteVar("Start", sm.Start)       //"00:00:00.0000"
		stdin.WriteVar("Duration", sm.Duration) //"00:00:03.0000"
		stdin.WriteVar("Output", sm.Output)     //"./test_000.mp4"

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func PeriodicPushHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicPushFunc[SQItem] {
	return func(queue *sys.QueueStack[SQItem]) (SQItem, error) {
		var (
			err   error
			qItem *redis.ZWithKey
			sm    *session.SQMessage
		)

		if qItem, err = rds.BZPopMin(
			sys.SessionContext, 0,
			session.GetKeyName("segment", ":queue"),
		).Result(); err != nil {
			return SQItem{sm: sm}, err
		} else {
			err = json.Unmarshal([]byte(qItem.Member.(string)), &sm)
		}

		// litter.D(sm, sm.KeyId)
		return SQItem{sm: sm}, err
	}
}

func OnPushedHandler(cfg *sys.Config, rds *redis.Client) sys.OnPushedFunc[SQItem] {
	return func(item SQItem) {
		Encode(cfg, item.sm)

		// litter.D(rds.RPush(
		// 	sys.SessionContext,
		// 	session.GetKeyName("segment", ":done"),
		// 	[]string{item.sm.KeyId, item.sm.Output},
		// ).Err())

		litter.D("item pushed", item.sm)
	}
}

func OnRemovedHandler(cfg *sys.Config, rds *redis.Client) sys.OnRemovedFunc[SQItem] {
	return func(item SQItem) {
		var (
			err error
			// qItem []string = make([]string, 2)
			// keyId string
		)

		// if qItem, err = rds.BLPop(
		// 	sys.SessionContext, 0,
		// 	session.GetKeyName("segment", ":done"),
		// ).Result(); err != nil {
		// 	return
		// }
		// keyId = qItem[0]

		err = rds.SAdd(
			sys.SessionContext,
			session.GetKeyName("concat:ready", ":", item.sm.KeyId),
			item.sm.Output,
		).Err()

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
