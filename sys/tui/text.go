package tui

import (
	"fmt"
	"strings"
)

func parseTextData(pipeData T_PipeData) string {
	var sb strings.Builder
	for i := range len(pipeData) {
		line := pipeData[i]
		for j := range len(line) {
			col := line[j]
			sb.WriteString(fmt.Sprintf("%s\t", col))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
