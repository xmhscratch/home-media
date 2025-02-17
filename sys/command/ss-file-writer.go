package command

import (
	"context"
	"encoding/json"
	"home-media/sys"
	"home-media/sys/session"
	"strings"
)

func NewSessionFileWriter(cfg *sys.Config, sessionId string, fileKey string, attrName string) *SessionFileWriter {
	return &SessionFileWriter{
		sessionWriterAbstract: &sessionWriterAbstract{
			SessionId: sessionId,
			AttrName:  attrName,
			Config:    cfg,
		},
		FileKey: fileKey,
	}
}

func (ctx SessionFileWriter) Read(p []byte) (int, error) {
	return len(p), nil
}

func (ctx SessionFileWriter) Write(p []byte) (int, error) {
	var (
		err  error
		meta *session.FileMetaInfo
	)

	rds := sys.NewClient(ctx.Config)
	defer rds.Close()

	if fInfStr, err := rds.HGet(
		context.TODO(),
		session.GetKeyName(ctx.SessionId, ":files"),
		ctx.FileKey,
	).Result(); err != nil {
		return -1, err
	} else {
		if err := json.Unmarshal([]byte(fInfStr), &meta); err != nil {
			return -1, err
		}
	}

	switch strings.ToLower(ctx.AttrName) {
	case "dubs":
		{
			if err = json.Unmarshal([]byte(strings.TrimSpace(string(p))), &meta.Dubs); err != nil {
				return -1, err
			}
			break
		}
	case "subtitles":
		{
			if err = json.Unmarshal([]byte(strings.TrimSpace(string(p))), &meta.Subtitles); err != nil {
				return -1, err
			}
			break
		}
	default:
		break
	}

	if fInfByt, err := json.Marshal(meta); err != nil {
		return -1, nil
	} else {
		if err = rds.HSet(
			context.TODO(),
			session.GetKeyName(ctx.SessionId, ":files"),
			[]string{ctx.FileKey, string(fInfByt)},
		).Err(); err != nil {
			return -1, err
		}
	}

	return len(p), err
}

func (ctx SessionFileWriter) Close() error {
	return nil
}
