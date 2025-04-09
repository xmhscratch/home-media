package session

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"home-media/sys"
	"io"
	"math"
	"os"
	"path/filepath"
	"time"

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
func (ctx *File[FileTorrent]) InitTorrent(verifyData bool) (*torrent.Torrent, error) {
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
	if verifyData {
		mytor.VerifyData()
	}
	<-mytor.GotInfo()

	ctx.TorrentName = mytor.Name()
	var files map[string]FileMetaInfo = map[string]FileMetaInfo{}
	for _, file := range mytor.Files() {
		fileKey := sys.GenerateID(ctx.NodeID, file.Path())
		files[fileKey] = FileMetaInfo{
			Path:        file.Path(),
			Size:        file.Length(),
			SourceReady: 0,
		}
	}
	ctx.Files = files

	return mytor, err
}

// GetTorrent comment
func (ctx *File[FileTorrent]) DownloadTorrent(ginCtx *gin.Context, filePath string) (err error) {
	var (
		selFile *torrent.File
		reader  torrent.Reader // = ctx.Torrent.NewReader()
		mytor   *torrent.Torrent
	)

	filePath = GetFilePath(filePath)
	if mytor, err = ctx.InitTorrent(true); err != nil {
		return err
	}
	defer mytor.Drop()

	fmt.Println("download requested:", filePath)

breakFileSelected:
	for _, selFile = range mytor.Files() {
		var parts []string
		if mytor.Info().BestName() != metainfo.NoName {
			parts = append(parts, mytor.Info().BestName())
		}
		fullFilePath := filepath.Join(append(parts, selFile.FileInfo().PathUtf8...)...)

		if fullFilePath != filepath.Join("./", filePath) {
			continue
		}
		fmt.Println("selected for download:", filePath)
		break breakFileSelected
	}

	if selFile == nil {
		return errors.New("file does not exist")
	}

	if ginCtx.Request.Header.Get("Range") == "" {
		selFile.Download()

		fmt.Println("starting download..", filePath)

		go func() {
			var (
				bTotal  int64        = selFile.Length()
				fileKey string       = sys.GenerateID(ctx.NodeID, selFile.Path())
				ticker  *time.Ticker = time.NewTicker(time.Duration(500) * time.Millisecond)
			)
			defer ticker.Stop()

			progressCtx, cancel := context.WithTimeout(context.Background(), 6*time.Hour)
			defer cancel()

		breakCountPercentage:
			for range ticker.C {
				select {
				case <-ticker.C:
					var (
						stats             torrent.TorrentStats = mytor.Stats()
						bRead             int64                = stats.BytesReadUsefulData.Int64()
						precisePercentage float64              = (float64(bRead) / float64(bTotal)) * 100
						roundPercentage   float64              = math.Round(precisePercentage*100) / 100
					)

					// fmt.Println(fileKey, roundPercentage)
					ctx.notify(fileKey, roundPercentage)

					// fmt.Println(bRead, bTotal, precisePercentage, roundPercentage)
					// fmt.Println(bTotal, map[string]int64{
					// 	"BytesWritten":                stats.BytesWritten.Int64(),
					// 	"BytesWrittenData":            stats.BytesWrittenData.Int64(),
					// 	"BytesRead":                   stats.BytesRead.Int64(),
					// 	"BytesReadData":               stats.BytesReadData.Int64(),
					// 	"BytesReadUsefulData":         stats.BytesReadUsefulData.Int64(),
					// 	"BytesReadUsefulIntendedData": stats.BytesReadUsefulIntendedData.Int64(),
					// 	"ChunksWritten":               stats.ChunksWritten.Int64(),
					// 	"ChunksRead":                  stats.ChunksRead.Int64(),
					// 	"ChunksReadUseful":            stats.ChunksReadUseful.Int64(),
					// 	"ChunksReadWasted":            stats.ChunksReadWasted.Int64(),
					// 	"MetadataChunksRead":          stats.MetadataChunksRead.Int64(),
					// 	"PiecesDirtiedGood":           stats.PiecesDirtiedGood.Int64(),
					// 	"PiecesDirtiedBad":            stats.PiecesDirtiedBad.Int64(),
					// })

					if bRead >= bTotal {
						break breakCountPercentage
					}
				case <-progressCtx.Done():
					break breakCountPercentage
				case <-mytor.Complete().On():
					break breakCountPercentage
				}
			}
		}()

		reader = selFile.NewReader()
		defer reader.Close()
		io.Copy(ginCtx.Writer, reader)

		defer ginCtx.Request.Body.Close()
	}

	if ok := client.WaitAll(); ok {
		fmt.Println("ok")
	}

	// ===============================================
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

	// 	_ = missinggo.NewSectionReadSeeker(reader, selFile.Offset()+startByte, selFile.Offset()+endByte)

	// 	// ginCtx.Header("Content-Range", "bytes "+strconv.FormatInt(startByte, 10)+"-"+strconv.FormatInt(endByte, 10)+"/"+strconv.FormatInt(selFile.Length(), 10))
	// 	// ginCtx.Header("Content-Length", strconv.FormatInt(endByte-startByte+1, 10))

	// 	ginCtx.Writer.WriteHeader(http.StatusPartialContent)
	// 	io.Copy(ginCtx.Writer, reader)
	// } else {
	// 	ginCtx.Header("Content-Length", strconv.FormatInt(selFile.Length(), 10))

	// 	_ = missinggo.NewSectionReadSeeker(reader, selFile.Offset(), selFile.Length())

	// 	io.Copy(ginCtx.Writer, reader)
	// }

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
