package duration

import (
	"home-media/sys/sample"
	"home-media/sys/session"
	"path/filepath"
	"testing"
)

var SampleDQMessage *session.DQMessage = &session.DQMessage{
	SessionId:   sample.SampleSessionID,
	SavePath:    filepath.Join(sample.SampleNodeID, sample.SampleFilePath),
	DownloadURL: filepath.Join(sample.SampleSessionID, "2", sample.SampleFilePath),
}

func TestDuration(t *testing.T) {
	Update(sample.SampleConfig, SampleDQMessage)
}
