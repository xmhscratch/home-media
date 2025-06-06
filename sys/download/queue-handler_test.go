package download

import (
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

// func TestFileExists(t *testing.T) {
// 	fmt.Println(
// 		sys.CheckFileExists("/home/data/storage/c041112c23b3a441957ee47a/fe2359ea7d5a53f88574d59461325f26.eng2.vtt"),
// 		sys.CheckFileExists("/home/data/storage/c041112c23b3a441957ee47a/fe2359ea7d5a53f88574d59461325f26.jpn1.mp4"),
// 		sys.CheckFileExists("/home/data/storage/c041112c23b3a441957ee47a/fe2359ea7d5a53f88574d59461325f26.mp4"),
// 	)
// }

// func TestDownload(t *testing.T) {
// 	Start(sample.SampleConfig, SampleDQMessage)
// }

// func TestDuration(t *testing.T) {
// 	UpdateDuration(sample.SampleConfig, SampleDQMessage)
// }

// func TestSubtitle(t *testing.T) {
// 	UpdateSubtitles(sample.SampleConfig, SampleDQMessage)
// }

// func TestDub(t *testing.T) {
// 	UpdateDubs(sample.SampleConfig, SampleDQMessage)
// }

// func TestExtractVideo(t *testing.T) {
// 	fmt.Println(ExtractVideo(sample.SampleConfig, SampleDQMessage))
// }

// func TestExtractDubs(t *testing.T) {
// 	fmt.Println(ExtractDubs(sample.SampleConfig, SampleDQMessage))
// }

// func TestExtractSubtitles(t *testing.T) {
// 	fmt.Println(ExtractSubtitles(sample.SampleConfig, SampleDQMessage))
// }

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
