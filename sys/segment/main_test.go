package segment

import (
	"home-media/sys/sample"
	"home-media/sys/session"
	"testing"
)

func TestSegmentEncode(t *testing.T) {
	sm := &session.SQMessage{
		KeyId:    "3467902dcd01be2a54f1201ed23bae03ceff8510",
		Index:    2,
		Source:   "/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/[SubsPlease] Thunderbolt Fantasy S4 - 12 (1080p) [B25E935B].mkv",
		Start:    "00:10:00.0000",
		Duration: "00:05:00.0000",
		Output:   "/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/3ff5ab6fa16862e4c4fbf060f50ce69f288a4557_002",
	}
	Encode(sample.SampleConfig, sm)
}
