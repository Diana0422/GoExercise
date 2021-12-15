package worker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (

)

func findPattern(file *os.File, pattern string) bool {
	fileReader := bufio.NewReader(file)
	lines := 0
	lineIdx := 0

	for {
		line, err := fileReader.ReadString('\n')
		lineIdx++
		if lineIdx == 1 {
			fmt.Println(DELIMETER)
		}
		if strings.Contains(line, pattern) {
			founded := fmt.Sprintf("%s[%s]: %s", color.GreenString(file.Name()), color.GreenString(strconv.Itoa(lineIdx)), strings.ReplaceAll(line, pattern, color.New(color.Underline).Add(color.FgGreen).Sprintf("%s", pattern)))
			fmt.Println(founded + DELIMETER)

			lines++
			continue
		}
		if err == io.EOF {
			return lines > 0
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println(color.CyanString("Usage: "), color.HiGreenString("ggrep [d] pattern filename1 filename2 filename3"))
		fmt.Println(color.CyanString("d"), color.HiGreenString(" - find in directories"))
		return
	}
	isDir := os.Args[1] == "d"
	if isDir {
		findInDir(os.Args[2], os.Args[3:])
		fmt.Println(isDir)
	} else {
		find(os.Args[1], os.Args[2:])
	}
}
