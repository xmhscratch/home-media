package tui

import (
	"fmt"
	"home-media/sys/sample"
	"io"
	"log"
	"net"
	"strings"
	"testing"
)

func TestSendData(t *testing.T) {
	// pipeData, err := ParseInput(sample.Sample_ListInput)
	// if err != nil {
	// 	os.Exit(1)
	// }
	// fmt.Printf("%v\n", pipeData1)

	// for i := range len(pipeData) {
	// 	time.Sleep(time.Millisecond * time.Duration(50))
	// }

	// pipeData, err := ParseInput(line)
	// if err != nil {
	// 	// log.Fatal(err)
	// 	continue
	// }
	// fmt.Printf("%v\n", _notes)

	src := strings.NewReader(fmt.Sprintf("%s|%s", OUTPUT_VIEW_LIST, sample.Sample_ListInput))
	buf := make([]byte, 1)

	go func() {
		conn, err := net.Dial("unix", "/run/tuid.sock")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		_, err = io.CopyBuffer(conn, src, buf)
		if err != nil {
			log.Println(err)
		}
	}()
}
