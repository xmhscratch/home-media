package segment

import (
	"home-media/sys/command"
	"home-media/sys/runtime"
	"os"
	"path/filepath"
	"time"
)

func (ctx SQItem) Index() int {
	return int(time.Now().Unix())
}

func (ctx SQItem) Key() string {
	return ctx.sm.KeyId
}

func (ctx *SQItem) ReEncode() error {
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

			Main: SegmentShell,
		}

		stdin.WriteVar("ExecBin", filepath.Join(ctx.Config.BinPath, "./segment.sh"))
		stdin.WriteVar("Input", ctx.sm.Source)
		stdin.WriteVar("Start", ctx.sm.Start)       //"00:00:00.0000"
		stdin.WriteVar("Duration", ctx.sm.Duration) //"00:00:03.0000"
		stdin.WriteVar("Output", ctx.sm.Output)     //"./test_000.mp4"

		exitCode <- shell.Run()
	}()
	<-exitCode

	return nil
}

// func (ctx *SQItem) asd(keyId string) error {
// 	var (
// 		err         error
// 		concatPaths chan []string = make(chan []string)
// 		foundKeyId  chan string   = make(chan string)
// 	)

// 	go func() error {
// 		var (
// 			err     error
// 			listKey map[string]string
// 		)
// 		if listKey, err = rds.HGetAll(
// 			context.TODO(),
// 			session.GetKeyName("segment", ":count"),
// 		).Result(); err != nil {
// 			return err
// 		}
// 		for keyId, c := range listKey {
// 			var (
// 				diffCount  int64
// 				totalCount int64
// 			)

// 			if diffCount, err = rds.SInterCard(
// 				context.TODO(), 0,
// 				session.GetKeyName("concat:queue", ":", keyId),
// 				session.GetKeyName("concat:ready", ":", keyId),
// 			).Result(); err != nil {
// 				return err
// 			}

// 			if totalCount, err = strconv.ParseInt(c, 5<<1, 0); err != nil {
// 				return err
// 			}

// 			if diffCount != totalCount {
// 				continue
// 			}

// 			if result, err := rds.SPopN(
// 				context.TODO(),
// 				session.GetKeyName("concat:queue", ":", keyId),
// 				totalCount,
// 			).Result(); err != nil {
// 				return err
// 			} else {
// 				concatPaths <- result
// 			}
// 			foundKeyId <- keyId
// 			return err
// 		}
// 		return err
// 	}()

// 	return nil
// }

// func (ctx *SQItem) JoinSegments(keyId string, concatPaths []string) error {
// 	var exitCode chan int = make(chan int)

// 	go func() {
// 		stdin := command.NewCommandReader()
// 		stdout := command.NewNullWriter()
// 		stderr := command.NewNullWriter()

// 		shell := runtime.Shell{
// 			PID: os.Getpid(),

// 			Stdin:  stdin,
// 			Stdout: stdout,
// 			Stderr: stderr,

// 			Args: os.Args,

// 			Main: ConcatShell,
// 		}

// 		if err := ctx.writeSegmentFile(keyId, concatPaths); err != nil {
// 			fmt.Println(err)
// 			exitCode <- 9
// 		}

// 		stdin.WriteVar("ExecBin", filepath.Join(ctx.Config.BinPath, "./concat.sh"))
// 		stdin.WriteVar("Input", filepath.Join(filepath.Dir(concatPaths[0]), "segments.txt"))
// 		stdin.WriteVar("Output", filepath.Join(
// 			filepath.Dir(concatPaths[0]),
// 			sys.BuildString(keyId, filepath.Ext(concatPaths[0])),
// 		))

// 		exitCode <- shell.Run()
// 	}()
// 	<-exitCode

// 	return nil
// }

// func (ctx *SQItem) writeSegmentFile(keyId string, concatPaths []string) error {
// 	var (
// 		err  error
// 		file *os.File
// 	)

// 	if file, err = os.Create(filepath.Join(filepath.Dir(concatPaths[0]), "segments.txt")); err != nil {
// 		return err
// 	}
// 	defer file.Close()

// errorLoop:
// 	for _, filePath := range concatPaths {
// 		if _, err = file.WriteString(sys.BuildString("file ", `'`+filePath+`'`, "\n")); err != nil {
// 			break errorLoop
// 		}
// 	}

// 	return err
// }

// if err := rds.SAdd(
// 	sys.SessionContext,
// 	session.GetKeyName("concat:ready", ":", item.sm.KeyId),
// 	item.sm.Output,
// ).Err(); err != nil {
// 	fmt.Println(err)
// 	return
// }
