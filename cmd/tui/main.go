package main

import (
	"flag"
	"fmt"
	"home-media/sys/tui"
	"os"
	"strings"
)

// echo "test	test" | tui --component list --header "My Fave Things 1" arg1 arg2
func main() {
	var (
		cn     string
		header string
	)

	flag.StringVar(&cn, "c", "pipe", "Name of the component")
	flag.StringVar(&cn, "component", "pipe", "Name of the component")

	flag.StringVar(&header, "h", "Untitled", "Display header")
	flag.StringVar(&header, "header", "Untitled", "Display header")

	flag.Parse()

	// fmt.Println("Component name:", cn)
	args := flag.Args()
	// fmt.Println("Non-flag arguments:", args)

	var err error
	switch strings.ToLower(cn) {
	case "list":
		err = tui.NewComponent_List(header, args...)
	case "pipe":
		err = tui.NewComponent_Pipe(header, args...)
	default:
		break
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// echo "hello" | socat - UNIX-CONNECT:/tmp/stdout.sock

// conn, err := net.Dial("unix", "/tmp/stdout.sock")
// if err != nil {
// 	panic(err)
// }
// defer conn.Close()

// conn.Write([]byte("hello world\n"))
