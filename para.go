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
	r := Wrapper{maxCols: wrap}
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
	maxCols int
}

// isMarkdownStart determines whether the text opens with
// a Markdown section or list indicator
func isMarkdownStart(text string) bool {
	return strings.HasPrefix(text, "#") ||
		strings.HasPrefix(text, "-") ||
		strings.HasPrefix(text, "*")
}

// closesParagraph checks if the text ends in a full stop
func closesParagraph(text string) bool {
	return strings.HasSuffix(text, ".")
}

// Wraptext wraps text to column length, compacting paragraphs along the way.
// It respects lines that end in a period, as well as Markdown lists and sections
func (wr Wrapper) Wraptext(scanner *bufio.Scanner, writer *bufio.Writer) error {
	var carry int
	// flushRunningText puts a line break on the open end of the text
	flushRunningText := func() {
		if carry > 0 {
			writer.WriteString("\n")
			carry = 0
		}
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			// Respect paragraphs
			flushRunningText()
			writer.WriteString("\n")
			continue
		}
		if isMarkdownStart(line) {
			flushRunningText()
		}
		var wrapped string
		wrapped, carry = wr.wrapLine(line, carry)
		writer.WriteString(wrapped)
		if closesParagraph(line) || isMarkdownStart(line) {
			// Respect full stops, markdown
			flushRunningText()
		}
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	return writer.Flush()
}

// wrapLine wraps a single line to a max number of columns, possibly breaking it
// into more lines by adding newlines in-place between words.
// There may be carry from a previous line wrapping.
//
// returns the wrapped text, and the number of carry-over characters
//
// NOTE: the input lines come from a line-scanner and don't have the terminating `\n`
func (r Wrapper) wrapLine(line string, carry int) (string, int) {
	lastWhite := -1
	lastNewline := -carry - 1

	out := make([]byte, len(line))
	copy(out, line)
	var startWithBreak bool
	for j := 0; j < len(line); j++ {
		if unicode.IsSpace(rune(line[j])) {
			lastWhite = j
		}
		if j-lastWhite > r.maxCols {
			log.Fatal("Word exceeds maxCols, line cannot be wrapped: " + line)
		}
		if j-lastNewline > r.maxCols && lastWhite > -1 {
			// we exceeded maxcols and can put a linebreak between previous words
			out[lastWhite] = '\n'
			lastNewline = lastWhite
		} else if j-lastNewline > r.maxCols && lastNewline < -1 {
			// counting the carry, we would exceed maxCols.
			// We should break from the previous fragment
			startWithBreak = true
			lastNewline = -1
		} else if j-lastNewline > r.maxCols {
			panic("Should never get here")
		}
	}

	var wrapped string
	switch {
	case startWithBreak:
		// we were unable to compress by appending to the previous carry
		wrapped = "\n" + string(out)
	case carry > 0:
		// we appended text to previous carry - we add a word break
		wrapped = " " + string(out)
	default:
		wrapped = string(out)
	}

	lastBreak := strings.LastIndex(wrapped, "\n")
	newCarry := len(wrapped) - lastBreak

	return wrapped, newCarry
}
