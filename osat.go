package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	sort.Strings(lines)
	lastTable := ""
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "ALTER TABLE ") &&
		   strings.HasSuffix(lines[i], ";") {
			space := strings.Index(lines[i][12:], " ")
			if space > 0 {
				table := lines[i][12 : space+12]
				if table != "" && table == lastTable {
					lines[i-1] = strings.Replace(lines[i-1], ";", ",", -1)
					lines[i] = "\t\t" + lines[i][space+13:]
				}
				lastTable = table
				continue;
			}
		}
		lastTable = ""
	}
	for i := range lines {
		fmt.Printf("%s\n", lines[i])
	}
}
