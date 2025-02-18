package session

import (
	"context"
	"encoding/json"
	"home-media/sys"
	"os"
	"path/filepath"
	// "github.com/sanity-io/litter"
	// "mime"
)

func (T FileSourceType) InitSession(
	cfg *sys.Config,
	sessionId string,
) (ctx *Session[FileSourceType], isCreated bool, err error) {
	keyName := GetKeyName(sessionId)

	rds := sys.NewClient(cfg)
	defer rds.Close()

	switch T.String() {
	case FILE_SOURCE_TYPE_DIRECT.String():
		ctx = &Session[FileSourceType]{
			ID:      sessionId,
			KeyName: keyName,
			Config:  cfg,
			File: &File[FileSourceType]{
				SourceType: FILE_SOURCE_TYPE_DIRECT,
				FileDirect: &FileDirect{},
			},
		}
	case FILE_SOURCE_TYPE_TORRENT.String():
		ctx = &Session[FileSourceType]{
			ID:      sessionId,
			KeyName: keyName,
			Config:  cfg,
			File: &File[FileSourceType]{
				SourceType:  FILE_SOURCE_TYPE_TORRENT,
				FileTorrent: &FileTorrent{},
			},
		}
	default:
		break
	}

	var ssCreated int64
	if ssCreated, err = rds.Exists(sys.SessionContext, sys.BuildString(keyName, ":info")).Result(); err != nil {
		return nil, false, err
	}

	if ssCreated == 0 {
		return ctx, false, err
	}

	var (
		b    []byte
		info map[string]string
	)
	if info, err = rds.HGetAll(sys.SessionContext, sys.BuildString(keyName, ":info")).Result(); err != nil {
		return nil, true, err
	}
	if b, err = json.Marshal(info); err != nil {
		return nil, true, err
	}

	if err = json.Unmarshal(b, ctx); err != nil {
		return nil, true, err
	}
	if err = json.Unmarshal(b, ctx.File); err != nil {
		return nil, true, err
	}

	if err = os.MkdirAll(filepath.Join(cfg.DataPath, ctx.NodeID), 0755); err != nil {
		return nil, true, err
	}

	ctx.File.notify = func(filePath string) error {
		return ctx.NotifyDownloaded()
	}

	return ctx, true, err
}

func InitDirect(
	cfg *sys.Config,
	sessionId string,
	opts ...string,
) (ctx *Session[FileSourceType], err error) {
	var (
		sourceURL string
		rootId    string
		nodeId    string
		setMode   bool = len(opts) >= 2
		isCreated bool = false
	)
	if setMode {
		sourceURL = opts[0]
		rootId = opts[1:2][0]
		nodeId = opts[2:3][0]
		sessionId = sys.GenerateID(sys.UUIDNamespace, nodeId)
	}

	rds := sys.NewClient(cfg)
	defer rds.Close()

	if ctx, isCreated, err = FILE_SOURCE_TYPE_DIRECT.InitSession(cfg, sessionId); err != nil {
		return nil, err
	}

	if setMode || !isCreated {
		ctx.RootID = rootId
		ctx.NodeID = nodeId
		ctx.File.NodeID = nodeId
		ctx.File.SourceURL = sourceURL
	}

	if extraInfo, err := ctx.File.InitDirect(); err != nil {
		return nil, err
	} else {
		if err := rds.HSet(
			sys.SessionContext,
			sys.BuildString(ctx.KeyName, ":info"),
			extraInfo,
		).Err(); err != nil {
			return nil, err
		}
	}

	if setMode || !isCreated {
		if err := rds.HSet(sys.SessionContext, sys.BuildString(ctx.KeyName, ":info"), map[string]interface{}{
			"rootId":     rootId,
			"nodeId":     nodeId,
			"savePath":   filepath.Join(ctx.NodeID),
			"sourceURL":  ctx.File.SourceURL,
			"sourceType": ctx.File.SourceType.String(),
			// "duration":    ctx.Duration,
		}).Err(); err != nil {
			return nil, err
		}
	}

	return ctx, err
}

func InitTorrent(
	cfg *sys.Config,
	sessionId string,
	opts ...string,
) (ctx *Session[FileSourceType], err error) {
	var (
		magnetURI string
		rootId    string
		nodeId    string
		setMode   bool = len(opts) >= 2
		isCreated bool = false
	)
	if setMode {
		magnetURI = opts[0]
		rootId = opts[1:2][0]
		nodeId = opts[2:3][0]
		sessionId = sys.GenerateID(sys.UUIDNamespace, nodeId)
	}

	rds := sys.NewClient(cfg)
	defer rds.Close()

	if ctx, isCreated, err = FILE_SOURCE_TYPE_TORRENT.InitSession(cfg, sessionId); err != nil {
		return nil, err
	}

	if setMode || !isCreated {
		ctx.RootID = rootId
		ctx.NodeID = nodeId
		ctx.File.NodeID = nodeId
		ctx.File.SourceURL = magnetURI
	}

	if _, err = ctx.File.InitTorrent(); err != nil {
		return nil, err
	}
	// defer ctx.File.Client.Close()

	if setMode || !isCreated {
		if err := rds.HSet(context.TODO(), sys.BuildString(ctx.KeyName, ":info"), map[string]interface{}{
			"torrentName": ctx.File.TorrentName,
			"rootId":      rootId,
			"nodeId":      nodeId,
			"savePath":    filepath.Join(ctx.NodeID),
			"sourceUrl":   ctx.File.SourceURL,
			"sourceType":  ctx.File.SourceType.String(),
			// "duration":    ctx.Duration,
		}).Err(); err != nil {
			// litter.D(err)
			return nil, err
		}

		if err := rds.HSet(
			sys.SessionContext,
			sys.BuildString(ctx.KeyName, ":files"),
			ctx.File.Files.ToArray()...,
		).Err(); err != nil {
			return nil, err
		}
	}

	if ctx.Files, err = func() (map[string]FileMetaInfo, error) {
		var (
			err      error
			b        []byte
			rawFiles map[string]string
			files    map[string]FileMetaInfo
		)

		if rawFiles, err = rds.HGetAll(sys.SessionContext, sys.BuildString(ctx.KeyName, ":files")).Result(); err != nil {
			return nil, err
		}
		if b, err = json.Marshal(rawFiles); err != nil {
			return nil, err
		}
		if err = json.Unmarshal(b, &files); err != nil {
			return nil, err
		}
		return files, err
	}(); err != nil {
		return nil, err
	}

	return ctx, err
}

func CreateDownload(
	cfg *sys.Config,
	sessionId string,
	filePath string,
) error {
	var (
		err        error
		keyName    string = GetKeyName(sessionId)
		sourceType FileSourceType
		ctx        *Session[FileSourceType]
		// isCreated   bool
	)

	rds := sys.NewClient(cfg)
	defer rds.Close()

	if result, err := rds.HGet(sys.SessionContext, sys.BuildString(keyName, ":info"), "sourceType").Result(); err != nil {
		return err
	} else {
		if sourceType, err = ParseSourceType(result); err != nil {
			return err
		}
	}

	switch sourceType {
	case FILE_SOURCE_TYPE_DIRECT:
		{
			if ctx, err = InitDirect(cfg, sessionId); err != nil {
				return err
			}
			break
		}
	case FILE_SOURCE_TYPE_TORRENT:
		{
			if ctx, err = InitTorrent(cfg, sessionId); err != nil {
				return err
			}
			break
		}
	}

	var dm *DQMessage
	for fileKey, v := range ctx.Files {
		fileMeta := FileMetaInfo(v)
		if filepath.Join("./", filePath) != filepath.Join("./", fileMeta.Path) {
			continue
		}
		if fileMeta.SourceReady == 1 {
			continue
		}
		dm = BuildDQMessage(ctx.NodeID, sessionId, sourceType, fileKey, &fileMeta)
	}

	if dm == nil {
		return err
	}

	if dmJSON, err := json.Marshal(dm); err != nil {
		return err
	} else {
		if err = rds.SAdd(
			sys.SessionContext,
			GetKeyName("download", ":queue"),
			string(dmJSON),
		).Err(); err != nil {
			return err
		}
	}

	return err

}

func (ctx *Session[I]) IsDownloadable() bool {
	return true
}

func (ctx *Session[I]) NotifyDownloaded() error {
	return nil
}
