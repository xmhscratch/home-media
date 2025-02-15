package extract

import (
	"testing"

	"home-media/sys/sample"
	"home-media/sys/session"
	"path/filepath"

	"github.com/sanity-io/litter"
)

var SampleDQMessage *session.DQMessage = &session.DQMessage{
	SessionId:   sample.SampleSessionID,
	FileType:    session.FILE_SOURCE_TYPE_TORRENT,
	SavePath:    filepath.Join(sample.SampleNodeID, sample.SampleFilePath),
	DownloadURL: filepath.Join(sample.SampleSessionID, "2", sample.SampleFilePath),
}

func TestExtractVideo(t *testing.T) {
	litter.D(ExtractVideo(sample.SampleConfig, SampleDQMessage))
}

func TestExtractSubtitles(t *testing.T) {
	litter.D(ExtractSubtitles(sample.SampleConfig, SampleDQMessage))
}

func TestExtractDubs(t *testing.T) {
	litter.D(ExtractDubs(sample.SampleConfig, SampleDQMessage))
}
