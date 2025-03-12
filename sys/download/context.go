package download

import (
	"home-media/sys"
	"home-media/sys/session"

	"github.com/redis/go-redis/v9"
)

type DQItem struct {
	sys.QItem[DQItem]
	HasOriginSource bool               `json:"hasOriginSource"`
	HasUpdtDuration bool               `json:"hasUpdtDuration"`
	HasExtrSubtitle bool               `json:"hasExtrSubtitle"`
	HasExtrAudio    bool               `json:"hasExtrAudio"`
	HasExtrVideo    bool               `json:"hasExtrVideo"`
	cfg             *sys.Config        `json:"-"`
	rds             *redis.Client      `json:"-"`
	dm              *session.DQMessage `json:"-"`
}
