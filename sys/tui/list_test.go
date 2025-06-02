package tui_test

import (
	"home-media/sys/tui"
	"log"
	"os"
	"testing"
)

var Sample_ListInput = `
Raspberry Pi’s			I have ’em all over my house	Nutella			It's good on toast
Bitter melon			It cools you down				Nice socks		And by that I mean socks without holes
Eight hours of sleep	I had this once					Cats			Usually
`

func TestFromRawPipe(t *testing.T) {
	rawInput := Sample_ListInput
	log.Println(tui.FromRawPipe(rawInput))
}

func TestList(t *testing.T) {
	// rawInput := Sample_ListInput
	rawInput, err := tui.ReadPipe()
	log.Println(rawInput)
	if err != nil {
		os.Exit(1)
		// return err
	}
	data, err := tui.FromRawPipe(rawInput)
	m := &tui.ListModel{TuiModel: &tui.TuiModel{PipeData: data}}
	if err != nil {
		os.Exit(1)
		// return err
	}
	m.MarshalData()
	log.Println()
}
