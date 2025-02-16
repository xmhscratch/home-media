package extract

import (
	"fmt"
	"testing"

	"home-media/sys/sample"
	"home-media/sys/session"
	"path/filepath"
)

var SampleDQMessage *session.DQMessage = &session.DQMessage{
	SessionId:   sample.SampleSessionID,
	FileType:    session.FILE_SOURCE_TYPE_TORRENT,
	SavePath:    filepath.Join(sample.SampleNodeID, sample.SampleFilePath),
	DownloadURL: filepath.Join(sample.SampleSessionID, "2", sample.SampleFilePath),
}

func TestExtractVideo(t *testing.T) {
	fmt.Println(ExtractVideo(sample.SampleConfig, SampleDQMessage))
}

func TestExtractDubs(t *testing.T) {
	fmt.Println(ExtractDubs(sample.SampleConfig, SampleDQMessage))
}

func TestExtractSubtitles(t *testing.T) {
	fmt.Println(ExtractSubtitles(sample.SampleConfig, SampleDQMessage))
}
