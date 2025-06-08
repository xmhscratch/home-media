package tui

import (
	"fmt"
	"home-media/sys/sample"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

func TestInstaller(t *testing.T) {
	pipeData1, err := ParseInput(sample.Sample_InstallPackages)
	if err != nil {
		os.Exit(1)
	}
	// fmt.Printf("%v\n", pipeData1)
	// parsePkgsData(pipeData1)

	// var (
	// 	myPackages []string       = []string{}
	// 	myNotes    map[int]string = map[int]string{}
	// 	myTotal    int            = 0
	// )

	// var c int = 0
	// var _cursor *int = &c
	// var _notes map[int]string = map[int]string{}

	for i := range len(pipeData1) {
		time.Sleep(time.Millisecond * time.Duration(50))

		line := pipeData1[i][0]

		// pipeData, err := ParseInput(line)
		// if err != nil {
		// 	// log.Fatal(err)
		// 	continue
		// }
		// _, notes, _ := insertPackagesData(pipeData, _cursor)
		// newNotes := lo.MapValues(notes, func(v string, k int) string {
		// 	if _, ok := _notes[k]; !ok {
		// 		return v
		// 	} else {
		// 		return strings.Join([]string{_notes[k], v}, "\n")
		// 	}
		// })
		// _notes = lo.Assign(_notes, newNotes)
		// fmt.Printf("%v\n", _notes)

		// if myTotal == 0 {
		// 	myTotal = total
		// }

		// src := strings.NewReader(fmt.Sprintf("%s|%s", OUTPUT_VIEW_INSTALLER, line))
		src := strings.NewReader(fmt.Sprintf("%s|%s", OUTPUT_VIEW_TEXT, line))
		buf := make([]byte, 1)

		go func() {
			log.Println(line)
			// log.Println(parsePkgsData(pipeData))

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

	// log.Println(myPackages, myNotes, myTotal)
}
