// Wrap text to a given column number. Acts as a filter from stdin to stdout.
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"unicode"
)

const default_wrap_col = 90

func main() {
	var wrap int
	if len(os.Args) < 2 {
		wrap = default_wrap_col
	} else {
		num, err := strconv.ParseInt(os.Args[1], 10, 0)
		if err != nil {
			log.Fatal(err)
		} else {
			wrap = int(num)
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	err := wraptext(wrap, scanner, writer)
	if err != nil {
		log.Fatal(err)
	}
}

// wrap text to column length, compact paragraph along the way
// respect lines that end in a period
func wraptext(column int, scanner *bufio.Scanner, writer *bufio.Writer) error {
	for scanner.Scan() {
		line := scanner.Text()
		writer.WriteString(wrapcol(column, line) + "\n")
	}
	return writer.Flush()
}

// wrap a single line to a colum length, possibly breaking it
func wrapcol(colnum int, line string) string {
	last_white, last_newline := 0, 0
	out := make([]rune, 2*len(line)) // at most 1 LF per old char
	for j, c := range line {
		out[j] = c
		if unicode.IsSpace(rune(c)) {
			last_white = j
		}
		if j-last_newline > colnum {
			out[last_white] = '\n'
			last_newline = last_white
		}
	}
	return string(out)
}
