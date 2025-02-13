package concat

import (
	"home-media/sys/sample"
	"testing"
)

func TestJoinSegment(t *testing.T) {
	Join(sample.SampleConfig, "3467902dcd01be2a54f1201ed23bae03ceff8510", []string{
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/3ff5ab6fa16862e4c4fbf060f50ce69f288a4557_000.mp4",
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/3ff5ab6fa16862e4c4fbf060f50ce69f288a4557_001.mp4",
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/3ff5ab6fa16862e4c4fbf060f50ce69f288a4557_002.mp4",
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/3ff5ab6fa16862e4c4fbf060f50ce69f288a4557_003.mp4",
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/3ff5ab6fa16862e4c4fbf060f50ce69f288a4557_004.mp4",
	})
}
