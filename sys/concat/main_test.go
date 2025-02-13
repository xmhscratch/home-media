package concat

import (
	"home-media/sys/sample"
	"testing"
)

func TestJoinSegment(t *testing.T) {
	Join(sample.SampleConfig, "ad288d149b79c2c409db528d4bd734e7fbb2a204", []string{
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/9e4095203ae99f1efb34ddaf399671e78c485cd1_000.mp4",
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/9e4095203ae99f1efb34ddaf399671e78c485cd1_001.mp4",
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/9e4095203ae99f1efb34ddaf399671e78c485cd1_002.mp4",
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/9e4095203ae99f1efb34ddaf399671e78c485cd1_003.mp4",
		"/home/web/repos/home-media/public/678bb5a27e785308b9e937a3/9e4095203ae99f1efb34ddaf399671e78c485cd1_004.mp4",
	})
}
