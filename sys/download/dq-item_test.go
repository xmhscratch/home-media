package download

import (
	"fmt"
	"home-media/sys/command"
	"home-media/sys/sample"
	"home-media/sys/session"
	"testing"
)

func TestUpdateSubtitles(t *testing.T) {
	stdout := command.NewSessionFileWriter(
		sample.SampleConfig,
		&session.FileMetaInfo{},
		sample.SampleSessionID,
		"95883c907fd158fcb0a46fb7b97dab33",
		"subtitles",
	)
	stdout.Write([]byte("[{\"stream_index\":1,\"codec_name\":\"eac3\",\"lang_code\":\"jpn\",\"lang_title\":\"Japanese (E-AC-3) (2.0) [AMZN]\"}]"))
	fmt.Println(stdout)
}

func TestUpdateDuration(t *testing.T) {
	stdout := command.NewSessionFileWriter(
		sample.SampleConfig,
		&session.FileMetaInfo{},
		sample.SampleSessionID,
		"95883c907fd158fcb0a46fb7b97dab33",
		"duration",
	)
	stdout.Write([]byte("1450.43"))
	fmt.Println(stdout)
}

func TestUpdateSourceReady(t *testing.T) {
	stdout := command.NewSessionFileWriter(
		sample.SampleConfig,
		&session.FileMetaInfo{},
		sample.SampleSessionID,
		"95883c907fd158fcb0a46fb7b97dab33",
		"sourceReady",
	)
	stdout.Write([]byte("1"))
	fmt.Println(stdout)
}
