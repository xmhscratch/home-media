package session

import (
	"home-media/sys/sample"
	"path/filepath"
)

var SampleDQMessage *DQMessage = &DQMessage{
	SessionId:   sample.SampleSessionID,
	SavePath:    filepath.Join(sample.SampleNodeID, sample.SampleFilePath),
	DownloadURL: filepath.Join(sample.SampleSessionID, "2", sample.SampleFilePath),
}
