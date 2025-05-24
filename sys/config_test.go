package sys

import (
	"testing"
)

func TestParseHostPort(t *testing.T) {
	cfg, _ := NewConfig("/home/web/repos/home-media")
	t.Log(cfg.ParseHostPort("frontend.hms:4200"))
}
