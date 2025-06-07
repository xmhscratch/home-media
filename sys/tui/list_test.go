package tui

import (
	"home-media/sys/sample"
	"log"
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
)

// import (
// 	"home-media/sys/sample"
// 	"log"
// 	"testing"
// )

func TestList(t *testing.T) {
	re := regexp2.MustCompile(`^((\d+(?=\|))((?=\|)..[^\|\n]*|)((?=\|).*)|.*)$`, regexp2.RE2)
	testcases := map[int]string{
		0: "1||hello||world",
		1: "1|2hello||world|1|432",
		2: "324234|2hello|world|",
		3: "2|2|233||2hello|world|",
		4: "2|2|233|2hello|world|",
		5: "1||hello",
		6: "1|hello",
		7: "1|||hello|||",
		8: "hello",
	}
	for _, v := range testcases {
		matches, _ := re.FindStringMatch(v)
		log.Println(matches.GroupByNumber(0))
		log.Println(matches.GroupByNumber(1).String())
		log.Println(matches.GroupByNumber(2).String())
		log.Println(strings.Trim(matches.GroupByNumber(3).String(), "|"))
		log.Println(strings.Trim(matches.GroupByNumber(4).String(), "|"))
		log.Println("========================")
	}
}

func TestParseListData(t *testing.T) {
	parseListData(sample.Sample_ListInput1)
	// rawInput := parseListData(sample.Sample_ListInput1)
	// log.Println(rawInput)
}

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
