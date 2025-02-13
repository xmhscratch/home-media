package session

import (
	"context"
	"encoding/json"
	"home-media/sys"
	"os"
	"path/filepath"
	"strings"
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
				SourceType:  FILE_SOURCE_TYPE_DIRECT,
				SourceReady: 0,
				FileDirect:  &FileDirect{},
			},
		}
	case FILE_SOURCE_TYPE_TORRENT.String():
		ctx = &Session[FileSourceType]{
			ID:      sessionId,
			KeyName: keyName,
			Config:  cfg,
			File: &File[FileSourceType]{
				SourceType:  FILE_SOURCE_TYPE_TORRENT,
				SourceReady: 0,
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

	ctx.File.notify = func(filePath string) error {
		var msg *DQMessage = BuildDQMessage(ctx.NodeID, sessionId, ctx.File.SourceType, filePath)
		return ctx.NotifyDownloaded(msg)
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
		savePaths []string
		rootId    string
		nodeId    string
		setMode   bool = len(opts) >= 2
		isCreated bool = false
	)
	if setMode {
		sourceURL = opts[0]
		savePaths = opts[1:]
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
			"rootId":      rootId,
			"nodeId":      nodeId,
			"savePath":    strings.Join(savePaths, string(os.PathSeparator)),
			"sourceURL":   ctx.File.SourceURL,
			"sourceType":  ctx.File.SourceType.String(),
			"sourceReady": ctx.File.SourceReady,
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
		savePaths []string
		rootId    string
		nodeId    string
		setMode   bool = len(opts) >= 2
		isCreated bool = false
	)
	if setMode {
		magnetURI = opts[0]
		savePaths = opts[1:]
		rootId = opts[1:2][0]
		nodeId = opts[2:3][0]
		sessionId = sys.GenerateID(sys.UUIDNamespace, nodeId)
	}

	rds := sys.NewClient(cfg)
	defer rds.Close()

	if ctx, isCreated, err = FILE_SOURCE_TYPE_TORRENT.InitSession(cfg, sessionId); err != nil {
		return nil, err
	}

	if !setMode && isCreated {
		var (
			err      error
			b        []byte
			rawFiles map[string]string
			files    map[string]*FileMetaInfo
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
	}

	if setMode || !isCreated {
		ctx.RootID = rootId
		ctx.NodeID = nodeId
		ctx.File.NodeID = nodeId
		ctx.File.SourceURL = magnetURI
	}

	if err = os.MkdirAll(filepath.Join(cfg.RootPath, cfg.DataDir, ctx.NodeID), 0755); err != nil {
		return nil, err
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
			"savePath":    strings.Join(savePaths, string(os.PathSeparator)),
			"sourceUrl":   ctx.File.SourceURL,
			"sourceType":  ctx.File.SourceType.String(),
			"sourceReady": ctx.File.SourceReady,
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

	return ctx, err
}

func GetFiles(
	cfg *sys.Config,
	sessionId string,
) (map[string]string, error) {
	var (
		err     error
		files   map[string]string
		keyName string = GetKeyName(sessionId)
	)

	rds := sys.NewClient(cfg)
	defer rds.Close()

	if files, err = rds.HGetAll(
		sys.SessionContext,
		sys.BuildString(keyName, ":files"),
	).Result(); err != nil {
		return nil, err
	}

	return files, err
}

func CreateDownload(
	cfg *sys.Config,
	sessionId string,
	filePath string,
) error {
	var (
		err         error
		keyName     string = GetKeyName(sessionId)
		sourceReady bool   = false
		sourceType  FileSourceType
		ctx         *Session[FileSourceType]
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

	if ctx, _, err = sourceType.InitSession(cfg, sessionId); err != nil {
		return err
	}

	sourceReady = (ctx.File.SourceReady == 1)
	if sourceReady {
		return err
	}

	if filePath == "" || filePath == "/" {
		// get first file founded
		for key := range ctx.Files {
			meta := ctx.Files.GetValue(key)
			filePath = meta.Path
			break
		}
	}

	var msg *DQMessage = BuildDQMessage(ctx.NodeID, sessionId, sourceType, filePath)

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err = rds.SAdd(
		sys.SessionContext,
		GetKeyName("download", ":queue"),
		string(msgJSON),
	).Err(); err != nil {
		return err
	}

	return err

}

func (ctx *Session[I]) IsDownloadable() bool {
	return true
}

func (ctx *Session[I]) NotifyDownloaded(msg *DQMessage) error {
	rds := sys.NewClient(ctx.Config)
	defer rds.Close()

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// litter.D("item downloaded: ", msg)

	return rds.SAdd(
		sys.SessionContext,
		GetKeyName("download", ":done"),
		string(msgJSON),
	).Err()
}
