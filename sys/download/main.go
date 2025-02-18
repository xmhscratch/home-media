package download

import (
	"context"
	"encoding/json"
	"fmt"
	"home-media/sys"
	"home-media/sys/session"
	"sync"

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

		// litter.D("queue added:", dm)
		return &DQItem{dm: dm, cfg: cfg, rds: rds}, err
	}
}

func ConsumeHandler(cfg *sys.Config, rds *redis.Client) sys.ConsumeFunc[DQItem] {
	return func(queue *sys.QueueStack[DQItem], item *DQItem) error {
		// fmt.Println(item)
		var (
			err error
			wg  *sync.WaitGroup = &sync.WaitGroup{}
		)

		if err = item.StartDownload(); err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			err = item.UpdateDuration()
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			err = map[int]error{
				0: item.UpdateSubtitles(),
				1: item.ExtractSubtitles(),
			}[0]
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			err = map[int]error{
				0: item.UpdateDubs(),
				1: item.ExtractDubs(),
			}[0]
			wg.Done()
		}()

		wg.Add(1)
		go func() {
			err = item.ExtractVideo()
			wg.Done()
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
		sm := &session.SQSegmentInfo{DQMessage: item.dm}
		sm.Init(cfg)

		// litter.D("queue completed:", item.dm, sm)
		if err := sm.PushQueue(); err != nil {
			fmt.Println(err)
		}
	}
}
