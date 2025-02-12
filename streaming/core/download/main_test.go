package download

import (
	"home-media/streaming/core/sample"
	"home-media/streaming/core/session"
	"path/filepath"
	"testing"
)

var SampleDQMessage *session.DQMessage = &session.DQMessage{
	SessionId:   sample.SampleSessionID,
	SavePath:    filepath.Join(sample.SampleNodeID, sample.SampleFilePath),
	DownloadURL: filepath.Join(sample.SampleSessionID, "2", sample.SampleFilePath),
}

func TestDownload(t *testing.T) {
	Start(sample.SampleConfig, SampleDQMessage)
}

// func TestTest(t *testing.T) {
// 	// var exitCode int = 0

// 	var waitExitCode chan int = make(chan int, 2)
// 	defer close(waitExitCode)

// 	// go func() {
// 	time.Sleep(time.Duration(2) * time.Second)
// 	waitExitCode <- 3
// 	// }()

// 	fmt.Println(<-waitExitCode)
// }
