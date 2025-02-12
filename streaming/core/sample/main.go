package sample

import (
	"home-media/streaming/core"
)

var SampleMagnetSource = "magnet:?xt=urn:btih:86e2d9a70a955856d086e24fa88a13adcadf7d11&dn=%5BSubsPlease%5D%20Thunderbolt%20Fantasy%20S4%20%2801-12%29%20%281080p%29%20%5BBatch%5D&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce"
var SampleRootID = "678bb472f4420b064e4ab471"
var SampleNodeID = "678bb5a27e785308b9e937a3"
var SampleSessionID = core.GenerateID(core.UUIDNamespace, SampleNodeID) //"2df240facfbf57d3afa6870490b905f5"
var SampleFilePath = "[SubsPlease] Thunderbolt Fantasy S4 - 12 (1080p) [B25E935B].mkv"

var SampleConfig *core.Config

func init() {
	var err error

	if SampleConfig, err = core.NewConfig("/home/web/repos/home-media"); err != nil {
		panic(err)
	}
}
