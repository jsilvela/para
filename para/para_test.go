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
	Rapper{maxcols: col}.wraptext(scanner, writer)
	return buffer.String()
}

func Test_wrap_single_line(t *testing.T) {
	test_s := "12 12 12 233 23 3 345 34"
	result := wraptext(10, test_s)
	lines := strings.Split(result, "\n")
	if len(test_s) != len(result) {
		t.Errorf("Content was altered %d, %d:\n%s\n%s", len(test_s), 
			len(result), test_s, result)
	}
	if len(lines) != 3 {
		t.Errorf("Missing newlines in: %s", result)
	}
}

func Test_wrap_lines(t *testing.T) {
	test_s := "My name is Wile E. Coyote"
	result := wraptext(15, test_s)
	ressy := strings.Map(func (x rune) rune {
								if x == '\n' {
									return ' '
								} else {
									return x
								} }, result)
	if test_s !=  ressy {
		t.Errorf("Content was altered %d, %d:\n%s\n%s", len(test_s), 
			len(result), test_s, ressy)
	}
	if len(strings.Split(result, "\n")) != 2 {
		t.Errorf("Bad wrap:\n%s", result)
	}
}

func Test_wrap_respects_one_blank_line(t *testing.T) {
	test_s := "My name is Wile E.\nCoyote.\n\nGenius"
	result := wraptext(15, test_s)
	if len(test_s) != len(result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(strings.Split(result, "\n")) != 5 {
		t.Errorf("Didn't respect blank line:\n%s", result)
	}
}

func Test_wrap_respects_two_blank_lines(t *testing.T) {
	test_s := "My name is E.\n\n\nCoyote"
	result := wraptext(15, test_s)
	if len(test_s) != len(result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(strings.Split(result, "\n")) != 4 {
		t.Errorf("Didn't respect two blank lines:\n%s", result)
	}
}

func Test_wrap_respects_full_stops(t *testing.T) {
	test_s := "My name is Wile E.\nCoyote"
	result := wraptext(15, test_s)
	if len(test_s) != len(result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(strings.Split(result, "\n")) != 3 {
		t.Errorf("Didn't respect full stop:\n%s", result)
	}
}

func Test_wrap_does_not_add_extra_breaks(t *testing.T) {
	test_s := "My name is."
	result := wraptext(15, test_s)
	if len(test_s) != len(result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(strings.Split(result, "\n")) != 1 {
		t.Errorf("Didn't respect period-break:\n%s", result)
	}
}
