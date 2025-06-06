package tui

import (
	"home-media/sys/sample"
	"log"
	"testing"
)

func TestFromRawPipe(t *testing.T) {
	// pipeData, err := ParseInput(sample.Sample_ListInput)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println(pipeData)
	parseListData(sample.Sample_ListInput1)
	log.Println()
}
