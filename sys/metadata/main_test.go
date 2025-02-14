package metadata

import (
	"home-media/sys/sample"
	"home-media/sys/session"
	"path/filepath"
	"testing"
)

var SampleDQMessage *session.DQMessage = &session.DQMessage{
	SessionId:   sample.SampleSessionID,
	FileType:    session.FILE_SOURCE_TYPE_TORRENT,
	SavePath:    filepath.Join(sample.SampleNodeID, sample.SampleFilePath),
	DownloadURL: filepath.Join(sample.SampleSessionID, "2", sample.SampleFilePath),
}

func TestDuration(t *testing.T) {
	UpdateDuration(sample.SampleConfig, SampleDQMessage)
}

func TestSubtitle(t *testing.T) {
	UpdateSubtitles(sample.SampleConfig, SampleDQMessage)
}

func TestDub(t *testing.T) {
	UpdateDubs(sample.SampleConfig, SampleDQMessage)
}
