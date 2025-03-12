package command

import (
	"context"
	"encoding/json"
	"home-media/sys"
	"home-media/sys/session"
	"strings"
)

func NewSessionFileWriter(
	cfg *sys.Config,
	fileMeta *session.FileMetaInfo,
	sessionId string,
	fileKey string,
	attrName string,
) *SessionFileWriter {
	return &SessionFileWriter{
		sessionWriterAbstract: &sessionWriterAbstract{
			SessionId: sessionId,
			AttrName:  attrName,
			Config:    cfg,
		},
		FileKey:  fileKey,
		FileMeta: fileMeta,
	}
}

func (ctx SessionFileWriter) Read(p []byte) (int, error) {
	return len(p), nil
}

func (ctx SessionFileWriter) Write(p []byte) (int, error) {
	var (
		err error
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
		if err := json.Unmarshal([]byte(fInfStr), &ctx.FileMeta); err != nil {
			return -1, err
		}
	}

	var attrValue any
	switch ctx.AttrName {
	case "dubs":
		attrValue = &ctx.FileMeta.Dubs
	case "subtitles":
		attrValue = &ctx.FileMeta.Subtitles
	case "duration":
		attrValue = &ctx.FileMeta.Duration
	case "sourceReady":
		attrValue = &ctx.FileMeta.SourceReady
	default:
		break
	}
	// fmt.Println(strings.TrimSpace(string(p)))
	if err = json.Unmarshal([]byte(strings.TrimSpace(string(p))), attrValue); err != nil {
		return -1, err
	}

	if fInfByt, err := json.Marshal(ctx.FileMeta); err != nil {
		return -1, nil
	} else {
		// fmt.Println(ctx.FileMeta, string(fInfByt))
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
