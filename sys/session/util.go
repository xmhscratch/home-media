package session

import (
	"crypto/sha1"
	"fmt"
	"home-media/sys"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
)

var FILE_PATH_REGEX *regexp.Regexp = regexp.MustCompile(`(([\\/]+[\d\w]{32}[\\/]+[12]{1})|[\\/]{0,})([^\\/]+)(\.[\w]{2,6})$`)

func GetFilePath(filePath string) string {
	var (
		err           error
		unescapedPath string
	)

	escapedPath := FILE_PATH_REGEX.ReplaceAllString(filePath, "$3$4")
	if unescapedPath, err = url.PathUnescape(escapedPath); err != nil {
		return escapedPath
	}
	return unescapedPath
}

func IsMatchingFilePath(filePath string) bool {
	return FILE_PATH_REGEX.MatchString(filePath)
}

func GetKeyName(sessionId string, parts ...string) string {
	result := sys.BuildString("hms_session:", sessionId)
	for _, val := range parts {
		result = sys.BuildString(result, val)
	}
	return result
}

func FormatDuration(length float64) string {
	d := time.Duration(length) * time.Second

	hour := int(d.Hours())
	minute := int(d.Minutes()) % 60
	second := int(d.Seconds()) % 60
	milisecond := int(d.Milliseconds()) % 1000

	return fmt.Sprintf("%02d:%02d:%02d.%04d", hour, minute, second, milisecond)
}

func BuildDQMessage(
	nodeId string,
	sessionId string,
	sourceType FileSourceType,
	filePath string,
) *DQMessage {
	var (
		savePath string = strings.Join(
			[]string{
				nodeId,
				filePath,
			},
			string(os.PathSeparator),
		)
		downloadUrl string = strings.Join(
			[]string{
				sessionId,
				sourceType.String(),
				filePath,
			},
			string(os.PathSeparator),
		)
	)

	return &DQMessage{
		SessionId:   sessionId,
		FileType:    sourceType,
		SavePath:    savePath,
		DownloadURL: downloadUrl,
	}
}

func GetFileKeyName(savePath string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(savePath)))
}

func GetMimeType(filePath string) string {
	re, err := regexp.Compile(`[\/]{0,1}([\w\W]+)+([\.]{1}[a-zA-Z0-9]+?)$`)
	matches := re.FindStringSubmatch(filePath)
	if err != nil {
		fmt.Println(err)
	}

	if len(matches) == 3 {
		return matches[2] // mime.TypeByExtension(matches[2])
	} else {
		return "" // "application/octet-stream"
	}
}

func minToSec(min float64) float64 {
	return min * 60
}

func torrentConfig() (config *torrent.ClientConfig) {
	config = torrent.NewDefaultClientConfig()

	config.DataDir = os.TempDir()
	ci := storage.NewFile(config.DataDir)
	defer ci.Close()

	config.DefaultStorage = ci

	// fmt.Println(config.DataDir)
	config.DisableTrackers = false
	config.NoDHT = false
	config.PeriodicallyAnnounceTorrentsToDht = true
	config.MaxUnverifiedBytes = 0

	config.ListenPort = 42104
	// config.NoDefaultPortForwarding = true
	// config.NoUpload = true
	config.Seed = false
	config.DisableUTP = true
	// config.DisableIPv6 = false
	// config.DisableIPv4 = false
	// config.DisableIPv4Peers = false
	// config.DisableTCP = false
	// config.Debug = true
	config.DisableWebtorrent = true
	config.DisableWebseeds = true

	return config
}
