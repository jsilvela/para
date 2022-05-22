package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func wraptext(col int, in string) string {
	reader := strings.NewReader(in)
	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	Wrapper{maxCols: col}.Wraptext(scanner, writer)
	return buffer.String()
}

func sameText(in, wr string) bool {
	despacer := func(x rune) rune {
		if x == '\n' || x == ' ' {
			return 'S'
		}
		return x
	}
	return strings.Map(despacer, wr) == strings.Map(despacer, in)
}

func TestWrapper_Compact(t *testing.T) {
	for _, spec := range []struct {
		input         string
		wrapAt        int
		expectedLines int
	}{
		{
			input:         "My name is Wile\nE Coyote\nGenius",
			wrapAt:        40,
			expectedLines: 1,
		},
	} {
		t.Run("wrapping"+spec.input, func(t *testing.T) {
			result := wraptext(spec.wrapAt, spec.input)
			if !sameText(result, spec.input) {
				t.Errorf("text was altered: %s", result)
			}
			if lines := len(strings.Split(result, "\n")); lines != spec.expectedLines {
				t.Errorf("expected %d lines, got %d", spec.expectedLines, lines)
			}
		})
	}
}

func TestWrapper_Break(t *testing.T) {
	for _, spec := range []struct {
		input         string
		wrapAt        int
		expectedLines int
	}{
		{
			input:         "12 12 12 2334 233",
			wrapAt:        8,
			expectedLines: 2,
		},
		{
			input:         "12 12 12 233 23 3 345 34",
			wrapAt:        10,
			expectedLines: 3,
		},
		{
			input:         "My name is Wile E. Coyote",
			wrapAt:        15,
			expectedLines: 2,
		},
		{
			input:         "My name is Wile E.\nCoyote.\n\nGenius",
			wrapAt:        15,
			expectedLines: 5,
		},
		{
			input:         "My name is Wile E.\nCoyote",
			wrapAt:        15,
			expectedLines: 3,
		},
	} {
		t.Run("wrapping"+spec.input, func(t *testing.T) {
			result := wraptext(spec.wrapAt, spec.input)
			if !sameText(result, spec.input) {
				t.Errorf("text was altered: %s", result)
			}
			if lines := len(strings.Split(result, "\n")); lines != spec.expectedLines {
				t.Errorf("expected %d lines, got %d", spec.expectedLines, lines)
			}
		})
	}
}

func TestWrapper_Passthrough(t *testing.T) {
	for _, spec := range []struct {
		name   string
		input  string
		wrapAt int
	}{
		{
			name:   "utf-8 text",
			input:  "Qué tal está señora?",
			wrapAt: 40,
		},
		{
			name:   "blank lines",
			input:  "My name is E.\n\n\nCoyote",
			wrapAt: 15,
		},
		{
			name:   "incompressible",
			input:  "12 456\n1012 is ok",
			wrapAt: 10,
		},
		{
			name: "markdown list compression",
			input: `# title
This is an enumeration:
## and a title right away

* one foo bar
baz quux
* two
`,
			wrapAt: 40,
		},
		{
			name: "markdown",
			input: `# title
This is an enumeration:

* one
* two

- three
- four

And then, that too`,
			wrapAt: 40,
		},
	} {
		t.Run(spec.name, func(t *testing.T) {
			result := wraptext(spec.wrapAt, spec.input)
			if result != spec.input {
				t.Errorf("testing '%s', result altered:\n%s", spec.name, result)
			}
		})
	}
}
