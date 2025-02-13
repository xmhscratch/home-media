package session

import (
	"crypto/sha1"
	"fmt"
	"home-media/sys"
	"io"
	"os"
	"path/filepath"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/gin-gonic/gin"
)

var client *torrent.Client

func init() {
	var (
		err    error
		config *torrent.ClientConfig
	)

	config = torrentConfig()
	if client, err = torrent.NewClient(config); err != nil {
		panic(err)
	}
}

// GetTorrent comment
func (ctx *File[FileTorrent]) InitTorrent() (*torrent.Torrent, error) {
	var (
		err   error
		spec  *torrent.TorrentSpec
		mytor *torrent.Torrent
	)

	if spec, err = torrent.TorrentSpecFromMagnetUri(ctx.SourceURL); err != nil {
		return nil, err
	}
	spec.Storage = fileStorage()

	if mytor, _, err = client.AddTorrentSpec(spec); err != nil {
		return nil, err
	}
	<-mytor.GotInfo()

	ctx.TorrentName = mytor.Name()
	var files map[string]FileMetaInfo = map[string]FileMetaInfo{}
	for _, file := range mytor.Files() {
		files[sys.GenerateID(ctx.NodeID, file.Path())] = FileMetaInfo{
			Path: file.Path(),
			Size: file.Length(),
		}
	}
	ctx.Files = files

	return mytor, err
}

// GetTorrent comment
func (ctx *File[FileTorrent]) DownloadTorrent(ginCtx *gin.Context, filePath string) (err error) {
	var (
		// fileSize     int64
		selectedFile *torrent.File
		reader       torrent.Reader // = ctx.Torrent.NewReader()
		mytor        *torrent.Torrent
	)

	filePath = GetFilePath(filePath)
	if mytor, err = ctx.InitTorrent(); err != nil {
		return err
	}
	fp := fmt.Sprintf("%s/%s", ctx.TorrentName, filePath)
	defer mytor.Drop()

	for _, file := range mytor.Files() {
		if file.Path() != fp {
			continue
		}
		selectedFile = file
	}
	// fileSize = selectedFile.Length()

	stats := mytor.Stats()
	if stats.ChunksWritten.Int64() > 0 {
		mytor.VerifyData()
	}

	if ginCtx.Request.Header.Get("Range") == "" {
		selectedFile.Download()

		reader = selectedFile.NewReader()
		defer reader.Close()
		io.Copy(ginCtx.Writer, reader)

		defer ginCtx.Request.Body.Close()
	}

	if ok := client.WaitAll(); ok {
		ctx.notify(filePath)
	}

	// ginCtx.Header("Accept-Ranges", "bytes")
	// ginCtx.Header("Cache-Control", "no-cache")

	// re, err := regexp.Compile(`[\/]{0,1}([\w\W]+)+([\.]{1}[a-zA-Z0-9]+?)$`)
	// matches := re.FindStringSubmatch(filePath)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// if len(matches) == 3 {
	// 	ginCtx.Header("Content-Type", mime.TypeByExtension(matches[2]))
	// } else {
	// 	ginCtx.Header("Content-Type", "application/octet-stream")
	// }
	// ginCtx.Header("Content-Disposition", "attachment; filename="+filePath)

	// if ginCtx.Request.Header.Get("Range") != "" {
	// 	rp := RangeParser.Parse(fileSize, ginCtx.Request.Header.Get("Range"))[0]

	// 	startByte := rp.Start
	// 	endByte := rp.End

	// 	_ = missinggo.NewSectionReadSeeker(reader, selectedFile.Offset()+startByte, selectedFile.Offset()+endByte)

	// 	// ginCtx.Header("Content-Range", "bytes "+strconv.FormatInt(startByte, 10)+"-"+strconv.FormatInt(endByte, 10)+"/"+strconv.FormatInt(selectedFile.Length(), 10))
	// 	// ginCtx.Header("Content-Length", strconv.FormatInt(endByte-startByte+1, 10))

	// 	ginCtx.Writer.WriteHeader(http.StatusPartialContent)
	// 	io.Copy(ginCtx.Writer, reader)
	// } else {
	// 	ginCtx.Header("Content-Length", strconv.FormatInt(selectedFile.Length(), 10))

	// 	_ = missinggo.NewSectionReadSeeker(reader, selectedFile.Offset(), selectedFile.Length())

	// 	io.Copy(ginCtx.Writer, reader)
	// }

	// fmt.Println("start downloading:", filePath)
	// ctx.Torrent.DownloadAll()

	return err
}

func fileStorage() storage.ClientImplCloser {
	pc, err := storage.NewDefaultPieceCompletionForDir(os.TempDir())
	if err != nil {
		pc = storage.NewMapPieceCompletion()
	}
	clientImplCloser := NewFileOpts(NewFileClientOpts{
		ClientBaseDir: os.TempDir(),
		TorrentDirMaker: func(baseDir string, info *metainfo.Info, infoHash metainfo.Hash) string {
			return os.TempDir()
		},
		FilePathMaker: func(opts storage.FilePathMakerOpts) string {
			var parts []string
			if opts.Info.BestName() != metainfo.NoName {
				parts = append(parts, opts.Info.BestName())
			}
			savePath := filepath.Join(append(parts, opts.File.BestPath()...)...)
			return sys.BuildString(fmt.Sprintf("%x", sha1.Sum([]byte(savePath))), ".tmp")
		},
		PieceCompletion: pc,
	})
	return clientImplCloser
}
