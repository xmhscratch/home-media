package command

import (
	"context"
	"home-media/sys"
	"home-media/sys/session"
	"strings"
)

func NewSessionInfoWriter(cfg *sys.Config, sessionId string, attrName string) *SessionInfoWriter {
	return &SessionInfoWriter{
		&sessionWriterAbstract{
			SessionId: sessionId,
			AttrName:  attrName,
			Config:    cfg,
		},
	}
}

func (ctx SessionInfoWriter) Read(p []byte) (int, error) {
	return len(p), nil
}

func (ctx SessionInfoWriter) Write(p []byte) (int, error) {
	rds := sys.NewClient(ctx.Config)
	defer rds.Close()

	_, err := rds.HSet(
		context.TODO(),
		session.GetKeyName(ctx.SessionId, ":info"),
		[]string{ctx.AttrName, strings.TrimSpace(string(p))},
	).Result()

	return len(p), err
}

func (ctx SessionInfoWriter) Close() error {
	return nil
}
