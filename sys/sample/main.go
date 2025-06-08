package sample

import (
	"home-media/sys"
)

var SampleMagnetSource = "magnet:?xt=urn:btih:b5dac11259fe56c1e9fbf7f487b943b769b4dc9d&dn=%5BErai-raws%5D%20Re%3AZero%20kara%20Hajimeru%20Isekai%20Seikatsu%203rd%20Season%20-%2010%20%5B1080p%20CR%20WEBRip%20HEVC%20EAC3%5D%5BMultiSub%5D%5B93760310%5D&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce"
var SampleRootID = "678bb472f4420b064e4ab471"
var SampleNodeID = "678bb5a27e785308b9e937a3"
var SampleSessionID = sys.GenerateID(sys.UUIDNamespace, SampleNodeID) //"2df240facfbf57d3afa6870490b905f5"
var SampleFilePath = "/[Erai-raws] Re Zero kara Hajimeru Isekai Seikatsu 3rd Season - 10 [1080p CR WEBRip HEVC EAC3][MultiSub][93760310].mkv"

var SampleConfig *sys.Config

var Sample_ListInput = `Raspberry Pi’s					I have ’em all over my house						Nutella							It's good on toast
Bitter melon					It cools you down
Nice socks						And by that I mean socks without holes
Eight hours of sleep			I had this once
Cats							Usually
Plantasia, the album			My plants love it too
Pour over coffee				It takes forever to make though
VR								Virtual reality...what is there to say?
Noguchi Lamps					Such pleasing organic forms
Linux							Pretty much the best OS
Business school					Just kidding
Pottery							Wet clay is a great feeling
Shampoo							Nothing like clean hair
Table tennis					It’s surprisingly exhausting
Milk crates						Great for packing in your extra stuff
Afternoon tea					Especially the tea sandwich part
Stickers						The thicker the vinyl the better
20° Weather						Celsius, not Fahrenheit
Warm light						Like around 2700 Kelvin
The vernal equinox				The autumnal equinox is pretty good too
Gaffer’s tape					Basically sticky fabric
Terrycloth						In other words, towel fabric`

var Sample_InstallPackages = `(1/21) Installing libcap-getcap (2.71-r0)
(2/21) Installing fakeroot (1.36-r0)
(3/21) Installing lzip (1.24.1-r1)
(4/21) Installing patch (2.7.6-r10)
(5/21) Installing pkgconf (2.3.0-r0)
(6/21) Installing acl-libs (2.3.2-r1)
(7/21) Installing tar (1.35-r2)
(8/21) Installing abuild (3.14.1-r4)
Executing abuild-3.14.1-r4.pre-install
(9/21) Installing abuild-sudo (3.14.1-r4)
(10/21) Installing libmagic (5.46-r2)
(11/21) Installing file (5.46-r2)
(12/21) Installing libstdc++-dev (14.2.0-r4)
(13/21) Installing g++ (14.2.0-r4)
(14/21) Installing make (4.4.1-r2)
(15/21) Installing fortify-headers (1.1-r5)
(16/21) Installing build-base (0.5-r3)
(17/21) Installing libexpat (2.7.0-r0)
(18/21) Installing git (2.47.2-r0)
(19/21) Installing git-init-template (2.47.2-r0)
(20/21) Installing alpine-sdk (1.1-r0)
(21/21) Installing git-zsh-completion (5.9-r4)
Executing busybox-1.37.0-r12.trigger
OK: 918 MiB in 139 packages`

func init() {
	var err error

	if SampleConfig, err = sys.NewConfig("/home/web/repos/home-media"); err != nil {
		panic(err)
	}
}
