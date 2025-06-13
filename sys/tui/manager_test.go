package tui

import (
	"fmt"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"home-media/sys/sample"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
)

func TestQuote(t *testing.T) {
	cmdArgs := []any{"asd1", "asd2", "asd3"}
	fmt.Printf("%s %s %s", cmdArgs...)

	// str := `./exec.sh "echo asd1 | tee -a ./test.txt"`
	// // print string with double quotes
	// fmt.Println(strconv.Quote(str))
}

func TestPayloadParsing(t *testing.T) {
	re := regexp2.MustCompile(RGXP_MESSAGE_PAYLOAD, regexp2.RE2|regexp2.Singleline)
	testPayload := fmt.Sprintf("6%cecho %%s | tee -a ./test.txt%c%s", ASCII_RS, ASCII_RS, sample.Sample_ListInput)
	matches, err := re.FindStringMatch(testPayload)
	// fmt.Println(RGXP_MESSAGE_PAYLOAD)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// fmt.Printf("%v", matches.GroupByNumber(4).String())

	var (
		messagePayload string = ""
		modePayload    string = ""
		optsPayload    string = ""
	)

	if matches.GroupByNumber(4).String() != "" {
		messagePayload = trimRS(matches.GroupByNumber(4).String())
	} else {
		messagePayload = matches.GroupByNumber(1).String()
	}

	if matches.GroupByNumber(2).String() != "" {
		modePayload = matches.GroupByNumber(2).String()
	} else {
		modePayload = OUTPUT_VIEW_TEXT.String()
	}

	if matches.GroupByNumber(3).String() != "" {
		optsPayload = trimRS(matches.GroupByNumber(3).String())
	}

	outputMode, err := strconv.Atoi(modePayload)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	log.Println(T_OutputMode(outputMode))
	log.Println(optsPayload)
	fmt.Printf("%s", messagePayload)
}

func TestExecCmd(t *testing.T) {
	var exitCode chan int = make(chan int)

	go func() {
		stdin := command.NewCommandReader()
		stdout := os.Stdout
		stderr := os.Stderr

		shell := runtime.Shell{
			PID: os.Getpid(),

			Stdin:  stdin,
			Stdout: stdout,
			Stderr: stderr,

			Args: os.Args,

			Main: CmdExecShell,
		}
		cmdArgs := []any{"asd1", "asd2", "asd3"}
		cmdStr := fmt.Sprintf(`echo %s | tee -a ./test.txt`, cmdArgs...)

		re := regexp2.MustCompile(`\%\!\(EXTRA[\ ]{0,}([a-z0-9]+(?:\=[a-z0-9]+[\, ]{0,})+)*\)$`, regexp2.RE2|regexp2.Singleline)
		genCmd, _ := re.Replace(cmdStr, "", -1, 1)
		log.Println(genCmd)

		stdin.WriteVar("ExecBin", filepath.Join("./webos/bin/exec.sh"))
		stdin.WriteVar("ExecArgs", strconv.Quote(genCmd))

		exitCode <- shell.Run()
	}()
	<-exitCode
	close(exitCode)
}

func TestSendExCmd(t *testing.T) {
	cmdStr := fmt.Sprintf("%s%csh%c%s", OUTPUT_SOCKET, ASCII_RS, ASCII_RS, `kubectl get pods -n kube-system -o jsonpath="{range .items[*]}{@.metadata.name}:{range @.status.conditions[*]}{@.type}={@.status};{end}{'\n'}{end}"`)
	// cmdStr := fmt.Sprintf("%s%csh%c%s", OUTPUT_SOCKET, ASCII_RS, ASCII_RS, `echo asdasd | tee -a /home/web/repos/home-media/test.txt`)
	fmt.Printf("%v\n", cmdStr)
	src := strings.NewReader(cmdStr)
	buf := make([]byte, 1)

	conn, err := net.Dial("unix", UNIX_EX_SOCKET_PATH)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = io.CopyBuffer(conn, src, buf)
	if err != nil {
		log.Println(err)
	}
}

func TestSendData(t *testing.T) {
	// pipeData, err := ParseInput(sample.Sample_ListInput)
	// if err != nil {
	// 	os.Exit(1)
	// }
	// fmt.Printf("%v\n", pipeData)

	// pipeData, err := ParseInput(line)
	// if err != nil {
	// 	// log.Fatal(err)
	// 	continue
	// }
	// fmt.Printf("%v\n", _notes)

	// src := strings.NewReader(fmt.Sprintf("%s|%s", OUTPUT_VIEW_TEXT, sample.Sample_ListInput))
	cmd := fmt.Sprintf("%s%cecho %%s | tee -a /home/web/repos/home-media/test.txt%c%s", OUTPUT_VIEW_LIST, ASCII_RS, ASCII_RS, sample.Sample_ListInput)
	fmt.Printf("%v\n", cmd)
	src := strings.NewReader(cmd)
	buf := make([]byte, 1)

	conn, err := net.Dial("unix", "/run/tuid.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = io.CopyBuffer(conn, src, buf)
	if err != nil {
		log.Println(err)
	}
}
