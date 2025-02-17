package segment

import (
	"home-media/sys"
	"home-media/sys/session"
)

type SQItem struct {
	sys.QItem[SQItem]
	Config      *sys.Config        `json:"-"`
	sm          *session.SQMessage `json:"-"`
	ConcatPaths []string
	KeyId       string
}
