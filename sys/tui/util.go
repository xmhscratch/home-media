package tui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

func ToRawData(rd io.Reader) (string, error) {
	var (
		err error
		b   strings.Builder
	)

	reader := bufio.NewReader(rd)

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

func ReadFromPipe() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return "", fmt.Errorf("try piping in some text")
	}
	return ToRawData(os.Stdin)
}

func (ctx T_OutputMode) String() string {
	return strconv.Itoa(int(ctx))
}

// func (ctx T_PipeData) DumpStruct(msg string) string {
// 	return map[int]map[int]string{0: map[int]string{0: msg.T_OutputMode.String(), 1: msg.string}}
// }

func ParseInput(rawInput string) (data T_PipeData, err error) {
	var (
		reader  *strings.Reader
		scanner *bufio.Scanner
		mu      *sync.Mutex = &sync.Mutex{}
	)

	mu.Lock()
	{
		inputScanner := bufio.NewScanner(strings.NewReader(rawInput))

		var sanitized string
		for inputScanner.Scan() {
			line := inputScanner.Text()
			cleanLine := normalizeTSVLine(line)
			sanitized += cleanLine + "\n"
		}
		// log.Println(sanitized)

		reader = strings.NewReader(sanitized)
		scanner = bufio.NewScanner(reader)

		i := 0
		data = T_PipeData{}
		for scanner.Scan() {
			line := scanner.Text()
			columns := strings.Split(line, "\t")
			data[i] = map[int]string{}
			for j, col := range columns {
				data[i][j] = col
			}
			i += 1
		}
	}
	mu.Unlock()

	if err := scanner.Err(); err != nil {
		return data, fmt.Errorf("error reading lines: %s", err)
	}

	return data, nil
}

func normalizeTSVLine(line string) string {
	re := regexp.MustCompile(`[\t]{1,}`)
	return re.ReplaceAllString(line, "\t")
}

// Generate a blend of colors.
func makeRampStyles(colorA, colorB string, steps float64) (s []lipgloss.Style) {
	cA, _ := colorful.Hex(colorA)
	cB, _ := colorful.Hex(colorB)

	for i := 0.0; i < steps; i++ {
		c := cA.BlendLuv(cB, i/steps)
		s = append(s, lipgloss.NewStyle().Foreground(lipgloss.Color(colorToHex(c))))
	}
	return
}

// Convert a colorful.Color to a hexadecimal format.
func colorToHex(c colorful.Color) string {
	return fmt.Sprintf("#%s%s%s", colorFloatToHex(c.R), colorFloatToHex(c.G), colorFloatToHex(c.B))
}

// Helper function for converting colors to hex. Assumes a value between 0 and
// 1.
func colorFloatToHex(f float64) (s string) {
	s = strconv.FormatInt(int64(f*255), 16)
	if len(s) == 1 {
		s = "0" + s
	}
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
