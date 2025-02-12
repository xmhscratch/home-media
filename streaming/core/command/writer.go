package command

import (
	"context"
	"home-media/streaming/core"
	"home-media/streaming/core/session"
	"strings"
)

func NewSessionWriter(cfg *core.Config, sessionId string, attrName string) *SessionWriter {
	ctx := &SessionWriter{
		SessionId: sessionId,
		AttrName:  attrName,
		Config:    cfg,
	}
	ctx.redis = core.NewClient(ctx.Config)
	return ctx
}

func (ctx SessionWriter) Read(p []byte) (int, error) {
	return len(p), nil
}

func (ctx SessionWriter) Write(p []byte) (int, error) {
	_, err := ctx.redis.HSet(
		context.TODO(),
		session.GetKeyName(ctx.SessionId, ":info"),
		[]string{ctx.AttrName, strings.TrimSpace(string(p))},
	).Result()

	return len(p), err
}

func (ctx SessionWriter) Close() error {
	return ctx.redis.Close()
}
