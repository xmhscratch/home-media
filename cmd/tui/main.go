package main

import (
	"fmt"
	"home-media/sys/tui"
	"os"
)

// echo "test	test" | tui --component list --header "My Fave Things 1" arg1 arg2
func main() {
	// var (
	// 	cn     string
	// 	header string
	// )

	// flag.StringVar(&cn, "c", "pipe", "Name of the component")
	// flag.StringVar(&cn, "component", "pipe", "Name of the component")

	// flag.StringVar(&header, "h", "Untitled", "Display header")
	// flag.StringVar(&header, "header", "Untitled", "Display header")

	// flag.Parse()

	// fmt.Println("Component name:", cn)
	// args := flag.Args()
	// fmt.Println("Non-flag arguments:", args)

	_, err := tui.NewTuiManager()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// go tm.ListenToSocket()

	// _ = tui.NewComponent_Text(mt, header, args...)
	// _ = tui.NewComponent_List(mt, header, args...)

	// for {
	// 	time.Sleep(time.Duration(2) * time.Second)

	// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 	slice := []string{"list", "text"}
	// 	cpn := slice[r.Intn(len(slice))]

	// 	switch strings.ToLower(cpn) {
	// 	case "list":
	// 		_ = tui.NewComponent_List(mt, header, args...)
	// 	case "text":
	// 		_ = tui.NewComponent_Text(mt, header, args...)
	// 	}
	// }
}

// echo "hello" | socat - UNIX-CONNECT:/run/tuid.sock
