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

		// fmt.Println("queue added:", dm)
		if err != nil {
			return nil, err
		}
		return (&DQItem{dm: dm, cfg: cfg, rds: rds}).VerifyInfo()
	}
}

func ConsumeHandler(cfg *sys.Config, rds *redis.Client) sys.ConsumeFunc[DQItem] {
	return func(queue *sys.QueueStack[DQItem], item *DQItem) error {
		var (
			err error
			wg  *sync.WaitGroup = &sync.WaitGroup{}
		)

	retryDownload:
		for range []int{1, 2, 3, 4, 5} {
			if !item.HasOriginSource {
				wg.Add(1)
				go func() {
					if err = (&session.FSMessage{
						Stage:   1,
						Message: "Downloading content...",
					}).SendToSocket(rds, item.dm.FileKey); err != nil {
						fmt.Println(err)
						return
					}
					time.Sleep(time.Duration(1) * time.Second)

					if err = item.StartDownload(); err != nil {
						return
					}
					wg.Done()
				}()
				wg.Wait()
			}

			if !item.HasUpdtDuration {
				wg.Add(1)
				go func() {
					if err = (&session.FSMessage{
						Stage:   3,
						Message: "Fetching file metadata...",
					}).SendToSocket(rds, item.dm.FileKey); err != nil {
						return
					}
					time.Sleep(time.Duration(1) * time.Second)

					err = item.UpdateDuration()
					wg.Done()
				}()
				wg.Wait()
			}

			if !item.HasExtrSubtitle {
				wg.Add(1)
				go func() {
					if err = (&session.FSMessage{
						Stage:   3,
						Message: "Extracting subtitles...",
					}).SendToSocket(rds, item.dm.FileKey); err != nil {
						return
					}
					time.Sleep(time.Duration(1) * time.Second)

					err = map[int]error{
						0: item.UpdateSubtitles(),
						1: item.ExtractSubtitles(),
					}[0]
					wg.Done()
				}()
				wg.Wait()
			}

			if !item.HasExtrAudio {
				wg.Add(1)
				go func() {
					if err = (&session.FSMessage{
						Stage:   3,
						Message: "Extracting audio...",
					}).SendToSocket(rds, item.dm.FileKey); err != nil {
						return
					}
					time.Sleep(time.Duration(1) * time.Second)

					err = map[int]error{
						0: item.UpdateDubs(),
						1: item.ExtractDubs(),
					}[0]
					wg.Done()
				}()
				wg.Wait()
			}

			if !item.HasExtrVideo {
				wg.Add(1)
				go func() {
					if err = (&session.FSMessage{
						Stage:   3,
						Message: "Extracting video...",
					}).SendToSocket(rds, item.dm.FileKey); err != nil {
						return
					}
					time.Sleep(time.Duration(1) * time.Second)

					err = item.ExtractVideo()
					wg.Done()
				}()
				wg.Wait()
			}

			item, err = item.VerifyInfo()

			p := map[int]bool{
				0: item.HasOriginSource,
				1: item.HasUpdtDuration,
				2: item.HasExtrSubtitle,
				3: item.HasExtrAudio,
				4: item.HasExtrVideo,
				5: err == nil,
			}
			// fmt.Println(p)
			if p[0] && p[1] && p[2] && p[3] && p[4] {
				item.UpdateSourceReady(true)
				break retryDownload
			}
		}

		return err
	}
}

func OnConsumedHandler(cfg *sys.Config, rds *redis.Client) sys.OnConsumedFunc[DQItem] {
	return func(item *DQItem) {
		if err := (&session.FSMessage{
			Stage:   4,
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
			Stage:   4,
			Message: "Completed!",
		}).SendToSocket(rds, item.dm.FileKey); err != nil {
			fmt.Println(err)
		}
		time.Sleep(time.Duration(1) * time.Second)

		if bMeta, err := json.Marshal(item.dm.FileMeta); err != nil {
			fmt.Println(err)
		} else {
			if err = (&session.FSMessage{
				Stage:   5,
				Message: string(bMeta),
			}).SendToSocket(rds, item.dm.FileKey); err != nil {
				fmt.Println(err)
			}
		}
	}
}
