package sample

import "home-media/sys"

var SampleMagnetSource = "magnet:?xt=urn:btih:b5dac11259fe56c1e9fbf7f487b943b769b4dc9d&dn=%5BErai-raws%5D%20Re%3AZero%20kara%20Hajimeru%20Isekai%20Seikatsu%203rd%20Season%20-%2010%20%5B1080p%20CR%20WEBRip%20HEVC%20EAC3%5D%5BMultiSub%5D%5B93760310%5D&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce"
var SampleRootID = "678bb472f4420b064e4ab471"
var SampleNodeID = "678bb5a27e785308b9e937a3"
var SampleSessionID = sys.GenerateID(sys.UUIDNamespace, SampleNodeID) //"2df240facfbf57d3afa6870490b905f5"
var SampleFilePath = "/[Erai-raws] Re Zero kara Hajimeru Isekai Seikatsu 3rd Season - 10 [1080p CR WEBRip HEVC EAC3][MultiSub][93760310].mkv"

var SampleConfig *sys.Config

var Sample_ListInput = `
Raspberry Pi’s			I have ’em all over my house	Nutella			It's good on toast
Bitter melon			It cools you down				Nice socks		And by that I mean socks without holes
Eight hours of sleep	I had this once					Cats			Usually
`

func init() {
	var err error

	if SampleConfig, err = sys.NewConfig("/home/web/repos/home-media"); err != nil {
		panic(err)
	}
}
