package session

import (
	"encoding/json"
	"testing"

	"github.com/sanity-io/litter"
)

var testUnmarshalJSON string = `{"28b60db601115a54ae3c3e9d6f89e62d":"{\"path\":\"[SubsPlease] Thunderbolt Fantasy S4 (01-12) (1080p) [Batch]/[SubsPlease] Thunderbolt Fantasy S4 - 02 (1080p) [F79FE9B7].mkv\",\"size\":1458633202}"}`

func TestMarshalJSON(t *testing.T) {
	var (
		err      error
		b        []byte
		rawFiles map[string]FileMetaInfo = map[string]FileMetaInfo{
			"28b60db601115a54ae3c3e9d6f89e62d": FileMetaInfo{
				Path: "[SubsPlease] Thunderbolt Fantasy S4 (01-12) (1080p) [Batch]/[SubsPlease] Thunderbolt Fantasy S4 - 02 (1080p) [F79FE9B7].mkv",
				Size: 1458633202,
			},
		}
	)
	b, err = json.Marshal(rawFiles)
	litter.D(err, string(b))
}

func TestToArray(t *testing.T) {
	var rawFiles FileMetaInfoList = map[string]FileMetaInfo{
		"28b60db601115a54ae3c3e9d6f89e62d": FileMetaInfo{
			Path: "[SubsPlease] Thunderbolt Fantasy S4 (01-12) (1080p) [Batch]/[SubsPlease] Thunderbolt Fantasy S4 - 02 (1080p) [F79FE9B7].mkv",
			Size: 1458633202,
		},
	}
	litter.D(rawFiles.ToArray())
}

func TestMarshalBinary(t *testing.T) {
	var (
		err      error
		b        []byte
		rawFiles FileMetaInfoList = map[string]FileMetaInfo{
			"28b60db601115a54ae3c3e9d6f89e62d": FileMetaInfo{
				Path: "[SubsPlease] Thunderbolt Fantasy S4 (01-12) (1080p) [Batch]/[SubsPlease] Thunderbolt Fantasy S4 - 02 (1080p) [F79FE9B7].mkv",
				Size: 1458633202,
			},
		}
	)
	// _, _ = rawFiles.MarshalBinary()
	b, err = rawFiles.MarshalBinary()
	litter.D(err, string(b))
}

func TestUnmarshalJSON(t *testing.T) {
	var (
		err   error
		files map[string]FileMetaInfo
	)
	err = json.Unmarshal([]byte(testUnmarshalJSON), &files)
	litter.D(err, files)
}

func TestUnmarshalBinary(t *testing.T) {
	var (
		err      error
		rawFiles FileMetaInfoList
	)
	err = rawFiles.UnmarshalBinary([]byte(testUnmarshalJSON))
	litter.D(err, rawFiles)
}
