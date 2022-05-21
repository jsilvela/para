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
			log.Fatalf("could not read column width: %v", err)
		} else {
			wrap = int(num)
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	r := Wrapper{maxcols: wrap}
	err := r.Wraptext(scanner, writer)
	if err != nil {
		log.Fatalf("could not wrap text: %v", err)
	}
}

// Wrapper contains methods for compressive text wrapping.
//
// It operates on line-oriented streams, as that is how it will be used from
// the command line.
type Wrapper struct {
	maxcols int
}

// isMarkdownStart determines whether the text opens with
// a Markdown section or list indicator
func isMarkdownStart(wrapped string) bool {
	return strings.HasPrefix(wrapped, "#") ||
		strings.HasPrefix(wrapped, "-") ||
		strings.HasPrefix(wrapped, "*")
}

// closesParagraph checks if the text ends in a full stop
func closesParagraph(wrapped string) bool {
	return strings.HasSuffix(wrapped, ".")
}

// Wraptext wraps text to column length, compacting paragraphs along the way.
// It respects lines that end in a period
func (wr Wrapper) Wraptext(scanner *bufio.Scanner, writer *bufio.Writer) error {
	var carry int
	closeText := func() {
		writer.WriteString("\n")
		carry = 0
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			// Respect paragraphs
			if carry > 0 {
				writer.WriteString("\n")
			}
			closeText()
			continue
		}
		if isMarkdownStart(line) && carry > 0 {
			// begin fresh, flush previous text
			closeText()
		}
		wrapped := wr.wrapLine(line, carry)
		writer.WriteString(wrapped)
		if closesParagraph(wrapped) || isMarkdownStart(wrapped) {
			// Respect full stops, markdown
			closeText()
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

// wrapLine wraps a single line to a max number of columns, possibly breaking it.
// There may be carry from a previous line wrapping
func (r Wrapper) wrapLine(line string, carry int) string {
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
			// we've exceeded maxcols and have seen a whitespace previously
			out[lastWhite] = '\n'
			lastNewline = lastWhite
		} else if j-lastNewline > r.maxcols && lastNewline < -1 {
			// counting the carry, we've exceeded maxcols without a space
			// we should break from the previous fragment
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
