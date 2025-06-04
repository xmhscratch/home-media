package tui

import (
	"home-media/sys/sample"
	"log"
	"testing"
)

func TestFromRawPipe(t *testing.T) {
	pipeData, err := ParseInput(sample.Sample_ListInput)
	if err != nil {
		panic(err)
	}
	// log.Println(pipeData)
	data := ""
	lx := []string{""}
	for i, line := range pipeData {
		// log.Println(i, line)
		lx = append(lx, "")
		for _, col := range line {
			// log.Println(col)
			lx[i] += col + "\t"
			// log.Println(lx[i])
		}
		data += lx[i] + "\n"
	}
	log.Println(data)
}
