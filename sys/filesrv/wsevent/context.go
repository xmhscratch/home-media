package wsevent

import (
	"time"

	"github.com/redis/go-redis/v9"
)

const WSEVENT_NAMESPACE = "3df0bd87e067452aa952d3915e319461"
const TIMEOUT_READ_MESSAGE = time.Duration(3000) * time.Millisecond

type SocketConnection struct {
	UUID       string
	SessionId  string
	FileKey    string
	Subscriber *redis.PubSub
}
