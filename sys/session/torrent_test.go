package session

import (
	"fmt"
	"home-media/sys/sample"
	"testing"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/sanity-io/litter"
)

func TestInitTorrent(t *testing.T) {
	var err error
	var client *torrent.Client

	var config = torrent.NewDefaultClientConfig()

	if client, err = torrent.NewClient(config); err != nil {
		t.Error(err)
		return
	}

	mytor, err := client.AddMagnet(sample.SampleMagnetSource)
	if err != nil {
		t.Error(err)
		return
	}
	<-mytor.GotInfo()

	var files []string = []string{}
	for i, file := range mytor.Files() {
		s := fmt.Sprintf("%c", i)
		files = append(files, s, file.Path())
	}
	litter.Dump(files)
}

func TestDownloadTorrent(t *testing.T) {
	t.Helper()

	var err error
	var client *torrent.Client

	var (
		// fileSize     int64
		selectedFile *torrent.File
	// reader       torrent.Reader = ctx.Torrent.NewReader()
	)

	var config = torrent.NewDefaultClientConfig()

	config.DisableTrackers = false
	config.NoDHT = false
	config.PeriodicallyAnnounceTorrentsToDht = true

	config.ListenPort = 42130
	// config.NoDefaultPortForwarding = true
	// config.NoUpload = true
	config.Seed = false
	config.DisableUTP = true
	// config.DisableTCP = false
	config.Debug = false
	config.DisableWebtorrent = true
	config.DisableWebseeds = true

	if client, err = torrent.NewClient(config); err != nil {
		t.Error(err)
		return
	}

	mytor, err := client.AddMagnet(sample.SampleMagnetSource)
	if err != nil {
		t.Error(err)
		return
	}
	<-mytor.GotInfo()

	defer client.Close()

	fp := fmt.Sprintf("%s/%s", mytor.Name(), sample.SampleFilePath)

out:
	for _, file := range mytor.Files() {
		if file.Path() != fp {
			continue
		}
		selectedFile = file
		break out
	}

	litter.D(selectedFile)
}

func TestTorrentTest(t *testing.T) {
	var (
		done chan struct{} = make(chan struct{})
	)

	go func() {
		var (
			// fileKey string       = sys.GenerateID(ctx.NodeID, selFile.Path())
			ticker *time.Ticker = time.NewTicker(time.Duration(500) * time.Millisecond)
		)
		defer ticker.Stop()

		// breakDownloadRate:
		for range ticker.C {
			select {
			case t := <-ticker.C:
				fmt.Println(t)
				// ctx.notify(fileKey, stats.BytesReadData.Int64())
				// default:
				// 	break breakDownloadRate
			}
		}
		done <- struct{}{}
	}()
	<-done
}
