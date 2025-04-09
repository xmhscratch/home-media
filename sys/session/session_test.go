package session

import (
	"encoding/json"
	"home-media/sys/sample"
	"testing"

	"github.com/sanity-io/litter"
)

func TestCreateSession(t *testing.T) {
	var (
		err error
		ss  *Session[FileSourceType]
	)

	if ss, err = InitTorrent(
		sample.SampleConfig,
		sample.SampleSessionID,
		sample.SampleMagnetSource,
		sample.SampleRootID, sample.SampleNodeID, sample.SampleFilePath,
	); err != nil {
		t.Fatal(err)
		return
	}

	if result, err := json.Marshal(ss); err != nil {
		t.Fatal(err)
		return
	} else {
		litter.Dump(string(result))
	}
}

func TestInitSession(t *testing.T) {
	var (
		err       error
		isCreated bool
		ss        *Session[FileSourceType]
	)

	if ss, isCreated, err = FILE_SOURCE_TYPE_TORRENT.InitSession(
		sample.SampleConfig,
		sample.SampleSessionID,
	); err != nil {
		t.Fatal(err)
		return
	}
	litter.Dump(ss, isCreated, err)
	// litter.Dump(ss.Files)
}

func TestGetSession(t *testing.T) {
	var (
		err error
		ss  *Session[FileSourceType]
	)

	if ss, err = InitTorrent(
		sample.SampleConfig,
		sample.SampleSessionID,
		// sample.SampleMagnetSource,
	); err != nil {
		t.Fatal(err)
		return
	}
	// ss = nil
	litter.Dump(ss)
}

func TestPrerequisite(t *testing.T) {
	var (
		err error
	)

	if err = Prerequisite(
		sample.SampleConfig,
		sample.SampleSessionID,
		sample.SampleFilePath,
	); err != nil {
		t.Fatal(err)
		return
	}
}

// func TestDownload(t *testing.T) {
// 	var (
// 		err error
// 		ss  *Session[FileSourceType]
// 	)

// 	if ss, err = InitTorrent(
// 		sample.SampleConfig,
// 		sample.SampleSessionID,
// 	); err != nil {
// 		t.Fatal(err)
// 		return
// 	}

// 	// ss.DownloadTorrent()

// 	// ss.NotifyDownload("", 0)
// 	// download.Start(sample.SampleConfig)
// }
