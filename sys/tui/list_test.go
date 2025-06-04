package tui

// import (
// 	"home-media/sys/sample"
// 	"log"
// 	"testing"
// )

// func TestFromRawPipe(t *testing.T) {
// 	rawInput := sample.Sample_ListInput
// 	log.Println(tui.FromRawPipe(rawInput))
// }

// func TestList(t *testing.T) {
// 	rawInput := sample.Sample_ListInput
// 	// rawInput, err := tui.ReadPipe()
// 	// log.Println(rawInput)
// 	// if err != nil {
// 	// 	os.Exit(1)
// 	// 	// return err
// 	// }
// 	data, err := tui.FromRawPipe(rawInput)
// 	m := &tui.ListModel{Root: &tui.TuiManager{PipeData: data}}
// 	if err != nil {
// 		os.Exit(1)
// 		// return err
// 	}
// 	log.Println(m.ParseData())
// }
