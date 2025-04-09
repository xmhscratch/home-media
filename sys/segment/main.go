package segment

import (
	"encoding/json"
	"home-media/sys"
	"home-media/sys/session"

	"github.com/redis/go-redis/v9"
)

func PeriodicHandler(cfg *sys.Config, rds *redis.Client) sys.PeriodicFunc[SQItem] {
	return func(queue *sys.QueueStack[SQItem]) (*SQItem, error) {
		var (
			err   error
			qItem *redis.ZWithKey
			sm    *session.SQMessage
		)

		if qItem, err = rds.BZPopMin(
			sys.SessionContext, 0,
			session.GetKeyName("segment", ":queue"),
		).Result(); err != nil {
			return nil, err
		} else {
			err = json.Unmarshal([]byte(qItem.Member.(string)), &sm)
		}

		// litter.D(sm, sm.KeyId)
		return &SQItem{sm: sm}, err
	}
}

func ConsumeHandler(cfg *sys.Config, rds *redis.Client) sys.ConsumeFunc[SQItem] {
	return func(queue *sys.QueueStack[SQItem], item *SQItem) error {
		if err := item.ReEncode(); err != nil {
			return err
		}
		return nil
	}
}

func OnConsumedHandler(cfg *sys.Config, rds *redis.Client) sys.OnConsumedFunc[SQItem] {
	return func(item *SQItem) {}
}
