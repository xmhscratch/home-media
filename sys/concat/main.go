package concat

import (
	"context"
	"fmt"
	"home-media/sys"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"home-media/sys/session"
	"os"
	"path/filepath"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func Join(cfg *sys.Config, keyId string, concatPaths []string) error {
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

		if err := WriteSegmentFile(cfg, keyId, concatPaths); err != nil {
			fmt.Println(err)
			exitCode <- 9
		}

		stdin.WriteVar("ExecBin", "/bin/home-media/concat.sh")
		stdin.WriteVar("Input", filepath.Join(filepath.Dir(concatPaths[0]), "segments.txt"))
		stdin.WriteVar("Output", filepath.Join(
			filepath.Dir(concatPaths[0]),
			sys.BuildString(keyId, filepath.Ext(concatPaths[0])),
		))

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

func WriteSegmentFile(cfg *sys.Config, keyId string, concatPaths []string) error {
	var (
		err  error
		file *os.File
	)

	if file, err = os.Create(filepath.Join(filepath.Dir(concatPaths[0]), "segments.txt")); err != nil {
		return err
	}
	defer file.Close()

errorLoop:
	for _, filePath := range concatPaths {
		if _, err = file.WriteString(sys.BuildString("file ", `'`+filePath+`'`, "\n")); err != nil {
			break errorLoop
		}
	}

	return err
}

func PeriodicPushHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicPushFunc {
	return func(queue map[string]interface{}) (interface{}, string, error) {
		var (
			err         error
			concatPaths chan []string = make(chan []string)
			foundKeyId  chan string   = make(chan string)
		)

		go func() error {
			var (
				err     error
				listKey map[string]string
			)
			if listKey, err = rds.HGetAll(
				context.TODO(),
				session.GetKeyName("segment", ":count"),
			).Result(); err != nil {
				return err
			}
			for keyId, c := range listKey {
				var (
					diffCount  int64
					totalCount int64
				)

				if diffCount, err = rds.SInterCard(
					context.TODO(), 0,
					session.GetKeyName("concat:queue", ":", keyId),
					session.GetKeyName("concat:ready", ":", keyId),
				).Result(); err != nil {
					return err
				}

				if totalCount, err = strconv.ParseInt(c, 5<<1, 0); err != nil {
					return err
				}

				if diffCount != totalCount {
					continue
				}

				if result, err := rds.SPopN(
					context.TODO(),
					session.GetKeyName("concat:queue", ":", keyId),
					totalCount,
				).Result(); err != nil {
					return err
				} else {
					concatPaths <- result
				}
				foundKeyId <- keyId
				return err
			}
			return err
		}()

		return <-concatPaths, <-foundKeyId, err
	}
}

func OnPushedHandler(cfg *sys.Config, rds *redis.Client) sys.OnPushedFunc {
	return func(item interface{}, keyId string) {
		var concatPaths []string = item.([]string)

		go func() error {
			return Join(cfg, keyId, concatPaths)
		}()
	}
}

func PeriodicRemoveHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicRemoveFunc {
	return func(queue map[string]interface{}) (string, error) {
		var (
			err   error
			keyId string
		)
		for keyId = range queue {
			break
		}
		return keyId, err
	}
}

func OnRemovedHandler(cfg *sys.Config, rds *redis.Client) sys.OnRemovedFunc {
	return func(item interface{}, keyId string) {
		var concatPaths []string = item.([]string)

		err := rds.SPopN(
			sys.SessionContext,
			session.GetKeyName("concat:ready", ":", keyId),
			int64(len(concatPaths)),
		).Err()

		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
