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

func TestParsePkgsData(t *testing.T) {
	pipeData1, err := ParseInput(sample.Sample_InstallPackages1)
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

	for i := range len(pipeData1) {
		time.Sleep(time.Second * time.Duration(1))

		line := pipeData1[i][0]

		// pipeData, err := ParseInput(line)
		// if err != nil {
		// 	// log.Fatal(err)
		// 	continue
		// }
		// packages, notes, total := parsePkgsData(pipeData)
		// fmt.Printf("%v\n", notes)

		// if myTotal == 0 {
		// 	myTotal = total
		// }

		// myPackages = append(myPackages, packages...)
		// myNotes[i] = lo.Reduce(lo.Times(len(notes), func(i int) int {
		// 	return int(i)
		// }), func(agg string, index int, _ int) string {
		// 	return strings.Join([]string{agg, notes[index]}, "\n")
		// }, "")
		// if myNotes[i] == "" {
		// 	delete(myNotes, i)
		// }

		src := strings.NewReader(fmt.Sprintf("%s|%s", OUTPUT_VIEW_INSTALLER, line))
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
