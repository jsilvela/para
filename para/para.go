// Wrap text to a given column number. Acts as a filter from stdin to stdout.
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

const default_wrap_col = 80

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
	r := Rapper{maxcols: wrap}
	err := r.wraptext(scanner, writer)
	if err != nil {
		log.Fatal(err)
	}
}

type Rapper struct {
	maxcols       int
	carry         int
	pending_break bool
}

// wrap text to column length, compact paragraph along the way
// respect lines that end in a period
func (r Rapper) wraptext(scanner *bufio.Scanner, writer *bufio.Writer) error {
	for scanner.Scan() {
		if r.pending_break {
			writer.WriteString("\n")
			r.pending_break = false
		}
		line := scanner.Text()
		wrapped := r.wrapline(line)
		writer.WriteString(wrapped)
		if strings.HasSuffix(line, ".") || len(line) == 0 {
			// Respect paragraphs and full stops
			r.pending_break = true
			r.carry = 0
		} else if strings.HasSuffix(wrapped, ".") {
			r.carry = 0
		} else {
			lastbrk := strings.LastIndex(wrapped, "\n")
			r.carry = len(line) - 1 - lastbrk
			if r.carry > 0 {
				r.carry = r.carry + 1
			}
		}
	}
	return writer.Flush()
}

// wrap a single line to a colum length, possibly breaking it
func (r Rapper) wrapline(line string) string {
	last_white, last_newline := -1, -r.carry-1
	out := make([]byte, len(line))
	var initbreak bool
	if len(line) == 0 && r.carry > 0 {
		initbreak = true
	}
	for j := 0; j < len(line); j++ {
		out[j] = line[j]
		if unicode.IsSpace(rune(line[j])) {
			last_white = j
		}
		if j-last_newline > r.maxcols && last_white > -1 {
			out[last_white] = '\n'
			last_newline = last_white
		} else if j-last_newline > r.maxcols {
			initbreak = true
			last_newline = last_white
		}
	}
	if initbreak {
		return "\n" + string(out)
	} else if r.carry > 0 {
		return " " + string(out)
	} else {
		return string(out)
	}
}
