package sys

import (
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *Config) *redis.Client {
	var (
		err  error
		opts *redis.Options
	)

	rdbConnStr := BuildString("redis://", cfg.Redis.Username, ":", cfg.Redis.Password, "@", cfg.Redis.Host, ":", cfg.Redis.Port, "/", cfg.Redis.DB)
	if opts, err = redis.ParseURL(rdbConnStr); err != nil {
		panic(err)
	}

	opts.Protocol = 3 // specify 2 for RESP 2 or 3 for RESP 3
	opts.DialTimeout = time.Duration(5) * time.Second
	opts.ReadTimeout = time.Duration(5) * time.Second
	opts.WriteTimeout = time.Duration(5) * time.Second
	// opts.PoolSize = 1

	ctx := redis.NewClient(opts)
	if dbNum, err := strconv.Atoi(cfg.Redis.DB); err != nil {
		ctx.Conn().Select(SessionContext, dbNum)
	}

	return ctx
}
