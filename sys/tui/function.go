package tui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReadPipe() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return "", fmt.Errorf("try piping in some text")
	}

	reader := bufio.NewReader(os.Stdin)
	var b strings.Builder

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return "", fmt.Errorf("error getting input: %s", err)
		}
	}
	return strings.TrimSpace(b.String()), err
}

func FromRawPipe(rawInput string) (data T_PipeData, err error) {
	// var mu *sync.Mutex = &sync.Mutex{}

	inputScanner := bufio.NewScanner(strings.NewReader(rawInput))

	var sanitized string
	for inputScanner.Scan() {
		line := inputScanner.Text()
		cleanLine := normalizeTSVLine(line)
		sanitized += cleanLine + "\n"
	}
	// log.Println(sanitized)

	reader := strings.NewReader(sanitized)
	scanner := bufio.NewScanner(reader)

	i := 0
	data = T_PipeData{}
	for scanner.Scan() {
		// mu.Lock()
		line := scanner.Text()
		columns := strings.Split(line, "\t")
		data[i] = map[int]string{}
		for j, col := range columns {
			data[i][j] = col
		}
		i += 1
		// mu.Unlock()
	}

	if err := scanner.Err(); err != nil {
		return data, fmt.Errorf("error reading lines: %s", err)
	}

	return data, nil
}

func NewTuiModel(args ...string) (*TuiModel, error) {
	rawInput, err := ReadPipe()
	if err != nil {
		return nil, err
	}
	data, err := FromRawPipe(rawInput)
	return &TuiModel{PipeData: data}, err
}

func NewComponent_List(header string, args ...string) error {
	mt, err := NewTuiModel(args...)
	m := &ListModel{TuiModel: mt, Header: header}
	if err := m.Render(); err != nil {
		return fmt.Errorf("error running program: %s", err)
	}
	return err
}

func NewComponent_Pipe(header string, args ...string) error {
	mt, err := NewTuiModel(args...)
	m := &PipeModel{TuiModel: mt}
	if err := m.Render(); err != nil {
		return fmt.Errorf("error running program: %s", err)
	}
	return err
}
