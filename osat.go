package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type statement struct {
	data []byte
}

func (s statement) split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if s.data == nil {
		s.data = make([]byte, len(data))
		copy(s.data, data)
	} else {
		s.data = append(s.data, data...)
	}
	if atEOF {
		return len(s.data), s.data, io.EOF
	}
	advance = strings.Index(string(s.data), ";\n")
	if advance >= 0 {
		advance += 2 // Include ;\n
		token = make([]byte, advance)
		copy(token, s.data[:advance-1]) // Next token
		copy(s.data, s.data[advance:])  // Remaining data
		return advance, token, nil
	}
	return 0, nil, nil
}

func (s statement) readAll(r io.Reader) []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(r)
	scanner.Split(s.split)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	sort.Strings(lines)
	return lines
}

func mergeLines(lines []string) []string {
	lastTable := ""
	for i := 0; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "ALTER TABLE ") ||
			!strings.HasSuffix(lines[i], ";") {
			lastTable = ""
			continue
		}
		space := strings.Index(lines[i][12:], " ")
		if space < 0 {
			continue
		}
		table := lines[i][12 : space+12]
		if table != "" && table == lastTable {
			lines[i-1] = strings.Replace(lines[i-1], ";", ",", -1)
			lines[i] = "\t\t" + lines[i][space+13:]
		}
		lastTable = table
	}
	return lines
}

func main() {
	stmt := &statement{}
	lines := mergeLines(stmt.readAll(os.Stdin))

	for i := range lines {
		fmt.Printf("%s\n", lines[i])
	}
}
