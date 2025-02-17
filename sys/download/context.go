package download

import (
	"home-media/sys"
	"home-media/sys/session"

	"github.com/redis/go-redis/v9"
)

type DQItem struct {
	sys.QItem[DQItem]
	cfg *sys.Config        `json:"-"`
	rds *redis.Client      `json:"-"`
	dm  *session.DQMessage `json:"-"`
}
