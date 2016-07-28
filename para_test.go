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
	Rapper{maxcols: col}.Wraptext(scanner, writer)
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

func Test_utf8_works(t *testing.T) {
	test := "Qué tal está señora?"
	result := wraptext(40, test)
	if test != result {
		t.Errorf("Content was altered:\n%s\n%s", test, result)
	}
}

func Test_limit_length_respected(t *testing.T) {
	test := `## Purity and the Mob

It happens in some places that live off tourism: at the restaurant where you're
having dinner, your friendly waiter suggests that if you're looking to party,
his brother in law has the best club in town. Ah, not in the mood for partying?`
	result := wraptext(80, test)
	if test != result {
		t.Errorf("Content was altered:\n%s\n%s", test, result)
	}
}

func Test_lines_at_limit_work(t *testing.T) {
	test := "12 12 12 2334 233"
	result := wraptext(8, test)
	lines := strings.Split(result, "\n")
	if !sameText(test, result) {
		t.Errorf("Content was altered:\n%s\n%s", test, result)
	}
	if len(lines) != 2 {
		t.Errorf("Missing newlines in: %s", result)
	}

}

func Test_wrap_single_line(t *testing.T) {
	test := "12 12 12 233 23 3 345 34"
	result := wraptext(10, test)
	lines := strings.Split(result, "\n")
	if !sameText(test, result) {
		t.Errorf("Content was altered:\n%s\n%s", test, result)
	}
	if len(lines) != 3 {
		t.Errorf("Missing newlines in: %s", result)
	}
}

func Test_wrap_lines(t *testing.T) {
	test := "My name is Wile E. Coyote"
	result := wraptext(15, test)
	if !sameText(test, result) {
		t.Errorf("Content was altered:\n%s\n%s", test, result)
	}
	if len(strings.Split(result, "\n")) != 2 {
		t.Errorf("Bad wrap:\n%s", result)
	}
}

func Test_wrap_respects_one_blank_line(t *testing.T) {
	test := "My name is Wile E.\nCoyote.\n\nGenius"
	result := wraptext(15, test)
	if !sameText(test, result) {
		t.Errorf("Content was altered:\n%s\n%s", test, result)
	}
	if len(strings.Split(result, "\n")) != 5 {
		t.Errorf("Didn't respect blank line:\n%s", result)
	}
}

func Test_wrap_respects_two_blank_lines(t *testing.T) {
	test := "My name is E.\n\n\nCoyote"
	result := wraptext(15, test)
	if test != result {
		t.Errorf("Content was altered:\n%s\n%s", test, result)
	}
}

func Test_wrap_respects_full_stops(t *testing.T) {
	test := "My name is Wile E.\nCoyote"
	result := wraptext(15, test)
	if !sameText(test, result) {
		t.Errorf("Content was altered:\n%s\n%s", test, result)
	}
	if len(strings.Split(result, "\n")) != 3 {
		t.Errorf("Didn't respect full stop:\n%s", result)
	}
}

func Test_lines_with_carry_wrap_to_limit(t *testing.T) {
	test := "12 456\n1012 is ok"
	result := wraptext(10, test)
	if result != test {
		t.Errorf("Unexpected result: %s", result)
	}
}
