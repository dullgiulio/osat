package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// emitLines emits every line of the file on the channel
func emitLines(r io.Reader, c chan<- string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		c <- scanner.Text()
	}
	close(c)
}

// merge sorts the lines (each line one statement); then groups all statements by
// table they modify.  Returns the group of statements as one single multiline string.
func merge(lines []string) string {
	var lastTable string
	sort.Strings(lines)
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
	return strings.Join(lines, "\n")
}

// transform reads, merges and outputs to the out channel an group of consecutive ALTER statements
func transform(out chan<- string, in <-chan string) {
	lines := make([]string, 1)
	lines[0] = line
	for line := range in {
		// End of the alter table statements: merge the found ones
		// and remember to print out this normal line
		if !strings.HasPrefix(line, "ALTER TABLE ") {
			out <- merge(lines)
			out <- line
			return
		}
		// Read up until the end of the statement in one line
		if !strings.HasSuffix(line, ";") {
			var buf bytes.Buffer
			buf.WriteString(line)
			// Read until the end of this query
			for l := range in {
				buf.WriteString(l)
				if strings.HasSuffix(l, ";") {
					break
				}
			}
			line = buf.String()
		}
		lines = append(lines, line)
	}
	out <- merge(lines)
}

// parse reads lines from the in channel, outputs transformed or verbatim lines to the out channel.
func parse(out chan<- string, in <-chan string) {
	for line := range in {
		// Normal statement, just print it out
		if !strings.HasPrefix(line, "ALTER TABLE ") {
			out <- line
			continue
		}
		// Read one or more ALTER TABLE statements
		transform(out, in)
	}
	close(out)
}

func main() {
	in := make(chan string)
	out := make(chan string)

	go emitLines(os.Stdin, in)
	go parse(out, in)

	for line := range out {
		fmt.Printf("%s\n", line)
	}
}
