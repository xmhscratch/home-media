package segment

import (
	"home-media/sys/sample"
	"home-media/sys/session"
	"testing"
)

func TestSegmentEncode(t *testing.T) {
	sm := &session.SQMessage{
		KeyId:    "ad288d149b79c2c409db528d4bd734e7fbb2a204",
		Index:    0,
		Source:   "/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/[Erai-raws] Re Zero kara Hajimeru Isekai Seikatsu 3rd Season - 10 [1080p CR WEBRip HEVC EAC3][MultiSub][93760310].mkv",
		Start:    "00:00:00.0000",
		Duration: "00:05:00.0000",
		Output:   "/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/9e4095203ae99f1efb34ddaf399671e78c485cd1_000",
	}
	Encode(sample.SampleConfig, sm)
}
