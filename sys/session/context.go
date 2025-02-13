package session

import (
	"home-media/sys"
	"net/http"
	"net/http/cookiejar"

	"github.com/juju/ratelimit"
)

// HTTPClientJar comment
var HTTPClientJar, _ = cookiejar.New(nil)
var HTTPClient = &http.Client{Jar: HTTPClientJar}

var DownloadBucketDefaultRate = int64(1024 * 1024 * 7) // 7Mbs
var DownloadBucket = ratelimit.NewBucketWithRate(float64(DownloadBucketDefaultRate), DownloadBucketDefaultRate)

type FileSourceType (int)

const (
	_                                       = iota
	FILE_SOURCE_TYPE_DIRECT  FileSourceType = 1
	FILE_SOURCE_TYPE_TORRENT FileSourceType = 2
)

const SEGMENT_CAPACITY int64 = 6

var fileSourceTypeName = map[FileSourceType]string{
	FILE_SOURCE_TYPE_DIRECT:  "1",
	FILE_SOURCE_TYPE_TORRENT: "2",
}

func (name FileSourceType) String() string {
	return fileSourceTypeName[name]
}

func (name FileSourceType) IsEqual(target string) bool {
	return target == name.String()
}

type File[T FileSourceType | FileDirect | FileTorrent] struct {
	*FileDirect  `json:"-"`
	*FileTorrent `json:"-"`
	NodeID       string             `json:"nodeId"`
	SourceURL    string             `json:"sourceUrl"`
	SourceType   FileSourceType     `json:"sourceType,string"`
	SourceReady  int                `json:"sourceReady,string"`
	notify       func(string) error `json:"-"`
}

type FileDirect struct {
	FileType string `json:"fileType"`
	FileSize int64  `json:"fileSize,string"`
}

type FileMetaInfoList map[string]FileMetaInfo

type FileMetaInfo struct {
	Path string `json:"path"`
	Size int64  `json:"size,string"`
}

type FileTorrent struct {
	TorrentName string           `json:"torrentName"`
	Files       FileMetaInfoList `json:"-"`
}

type Session[T FileSourceType] struct {
	*File[T] `json:"-"`
	ID       string           `json:"id"`
	KeyName  string           `json:"-"`
	RootID   string           `json:"rootId"`
	NodeID   string           `json:"nodeId"`
	Files    FileMetaInfoList `json:"files"`
	Config   *sys.Config      `json:"-"`
}

type DQMessage struct {
	SessionId   string `json:"sessionId"`
	SavePath    string `json:"savePath"`
	DownloadURL string `json:"downloadUrl"`
}

type SQMessage struct {
	KeyId    string `json:"keyId"`
	Index    int64  `json:"index"`
	Source   string `json:"source"`
	Start    string `json:"start"`
	Duration string `json:"duration"`
	Output   string `json:"output"`
}

type SQSegmentInfo struct {
	*DQMessage
	KeyId             string            `json:"keyId"`
	TotalLength       float64           `json:"totalLength"`
	SegmentLength     float64           `json:"segmentLength"`
	BestSegmentLength float64           `json:"bestSegmentLength"`
	BestSegmentCount  int64             `json:"bestSegmentCount"`
	Selected          *SQMessage        `json:"selected"`
	Segments          map[string]string `json:"-"`
	Config            *sys.Config       `json:"-"`
}
