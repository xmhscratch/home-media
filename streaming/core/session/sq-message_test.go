package session

import (
	"home-media/streaming/core/sample"
	"testing"

	"github.com/sanity-io/litter"
)

func TestSQMInit(t *testing.T) {
	msg := &SQSegmentInfo{DQMessage: SampleDQMessage}
	if err := msg.Init(sample.SampleConfig); err != nil {
		t.Fatal(err)
	}
	litter.D(msg)
}

func TestSQMPush(t *testing.T) {
	msg := &SQSegmentInfo{DQMessage: SampleDQMessage}
	if err := msg.Init(sample.SampleConfig); err != nil {
		t.Fatal(err)
	}

	// litter.D(msg)
	if err := msg.PushQueue(); err != nil {
		t.Fatal(err)
	}
}

func TestSQMFindBestSegmentValue(t *testing.T) {
	msg := &SQSegmentInfo{
		TotalLength: 967.061000, // 1440.25,
	}
	litter.D(msg.bestSegmentValue())
}
