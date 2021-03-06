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

const defaultWrapCol = 80

func main() {
	var wrap int
	if len(os.Args) < 2 {
		wrap = defaultWrapCol
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
	err := r.Wraptext(scanner, writer)
	if err != nil {
		log.Fatal(err)
	}
}

// Rapper contains methods for compressive text wrapping
type Rapper struct {
	maxcols int
}

// wrap text to column length, compact paragraph along the way
// respect lines that end in a period
func (r Rapper) Wraptext(scanner *bufio.Scanner, writer *bufio.Writer) error {
	var carry int
	endLine := func() {
		writer.WriteString("\n")
		carry = 0
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			// Respect paragraphs
			if carry > 0 {
				endLine()
			}
			writer.WriteString("\n")
			continue
		}
		wrapped := r.wrapline(line, carry)
		writer.WriteString(wrapped)
		if strings.HasSuffix(wrapped, ".") ||
			strings.HasPrefix(wrapped, "#") ||
			strings.HasPrefix(wrapped, "*") {
			// Respect full stops, markdown
			endLine()
		} else {
			lastBrk := strings.LastIndex(wrapped, "\n")
			carry = len(wrapped) - lastBrk
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return writer.Flush()
}

// wrap a single line to a *colum* length, possibly breaking it
// there may be carry from a previous line wrapping
func (r Rapper) wrapline(line string, carry int) string {
	lastWhite := -1
	lastNewline := -carry - 1

	out := make([]byte, len(line))
	var startWithBreak bool
	for j := 0; j < len(line); j++ {
		out[j] = line[j]
		if unicode.IsSpace(rune(line[j])) {
			lastWhite = j
		}
		if j-lastWhite > r.maxcols {
			log.Fatal("Word exceeds maxcols, line cannot be wrapped: " + line)
		}
		if j-lastNewline > r.maxcols && lastWhite > -1 {
			out[lastWhite] = '\n'
			lastNewline = lastWhite
		} else if j-lastNewline > r.maxcols && lastNewline < -1 {
			startWithBreak = true
			lastNewline = -1
		} else if j-lastNewline > r.maxcols {
			panic("Should never get here")
		}
	}
	if startWithBreak {
		return "\n" + string(out)
	} else if carry > 0 {
		return " " + string(out)
	} else {
		return string(out)
	}
}
