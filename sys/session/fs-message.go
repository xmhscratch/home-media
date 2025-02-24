package session

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type FSMessage struct {
	Stage   int    `json:"stage"`
	Message string `json:"message"`
}

func (ctx *FSMessage) SendToSocket(rds *redis.Client, fileKey string) error {
	var (
		err error
		b   []byte
	)

	if b, err = json.Marshal(ctx); err != nil {
		return err
	}
	// fmt.Println(fileKey, string(b))
	if err := rds.Publish(context.TODO(), fileKey, string(b)).Err(); err != nil {
		return err
	}

	return err
}
