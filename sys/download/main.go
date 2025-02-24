package download

import (
	"context"
	"encoding/json"
	"fmt"
	"home-media/sys"
	"home-media/sys/session"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func PeriodicHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicFunc[DQItem] {
	return func(queue *sys.QueueStack[DQItem]) (*DQItem, error) {
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
		// litter.D("queue adding...")
		if qItem, err = rds.SPop(
			context.TODO(),
			rdsKeyName,
		).Result(); err != nil {
			// fmt.Println(qItem, err)
			return nil, err
		} else {
			if err = json.Unmarshal([]byte(qItem), &dm); err != nil {
				return nil, err
			}
		}

		fmt.Println("queue added:", dm)
		return &DQItem{dm: dm, cfg: cfg, rds: rds}, err
	}
}

func ConsumeHandler(cfg *sys.Config, rds *redis.Client) sys.ConsumeFunc[DQItem] {
	return func(queue *sys.QueueStack[DQItem], item *DQItem) error {
		var (
			err error
			wg  *sync.WaitGroup = &sync.WaitGroup{}
		)

		if err = (&session.FSMessage{
			Stage:   1,
			Message: "Downloading content...",
		}).SendToSocket(rds, item.dm.FileKey); err != nil {
			fmt.Println(err)
			return err
		}
		time.Sleep(time.Duration(1) * time.Second)
		if err = item.StartDownload(); err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			err = item.UpdateDuration()
			wg.Done()
			if err = (&session.FSMessage{
				Stage:   2,
				Message: "Fetching file metadata...",
			}).SendToSocket(rds, item.dm.FileKey); err != nil {
				return
			}
			time.Sleep(time.Duration(1) * time.Second)
		}()

		wg.Add(1)
		go func() {
			err = map[int]error{
				0: item.UpdateSubtitles(),
				1: item.ExtractSubtitles(),
			}[0]
			wg.Done()
			if err = (&session.FSMessage{
				Stage:   2,
				Message: "Extracting subtitles...",
			}).SendToSocket(rds, item.dm.FileKey); err != nil {
				return
			}
			time.Sleep(time.Duration(1) * time.Second)
		}()

		wg.Add(1)
		go func() {
			err = map[int]error{
				0: item.UpdateDubs(),
				1: item.ExtractDubs(),
			}[0]
			wg.Done()
			if err = (&session.FSMessage{
				Stage:   2,
				Message: "Extracting audio...",
			}).SendToSocket(rds, item.dm.FileKey); err != nil {
				return
			}
			time.Sleep(time.Duration(1) * time.Second)
		}()

		wg.Add(1)
		go func() {
			err = item.ExtractVideo()
			wg.Done()
			if err = (&session.FSMessage{
				Stage:   2,
				Message: "Extracting video...",
			}).SendToSocket(rds, item.dm.FileKey); err != nil {
				return
			}
			time.Sleep(time.Duration(1) * time.Second)
		}()

		wg.Wait()
		if err == nil {
			item.UpdateSourceReady(true)
		}

		return err
	}
}

func OnConsumedHandler(cfg *sys.Config, rds *redis.Client) sys.OnConsumedFunc[DQItem] {
	return func(item *DQItem) {
		if err := (&session.FSMessage{
			Stage:   3,
			Message: "Re-encoding...",
		}).SendToSocket(rds, item.dm.FileKey); err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(1) * time.Second)

		sm := &session.SQSegmentInfo{DQMessage: item.dm}
		sm.Init(cfg)

		if err := (&session.FSMessage{
			Stage:   4,
			Message: "Segments concatenation...",
		}).SendToSocket(rds, item.dm.FileKey); err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(1) * time.Second)

		// litter.D("queue completed:", item.dm, sm)
		if err := sm.PushQueue(); err != nil {
			fmt.Println(err)
		}

		if err := (&session.FSMessage{
			Stage:   5,
			Message: "Completed!",
		}).SendToSocket(rds, item.dm.FileKey); err != nil {
			fmt.Println(err)
		}
	}
}
