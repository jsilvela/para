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

func same_text(in, wr string) bool {
	despacer := func(x rune) rune {
		if x == '\n' || x == ' ' {
			return 'S'
		} else {
			return x
		}
	}
	return strings.Map(despacer, wr) == strings.Map(despacer, in)
}

func Test_utf8_works(t *testing.T) {
	test_s := "Qué tal está señora?"
	result := wraptext(40, test_s)
	if !same_text(test_s, result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
}

func Test_limit_length_respected(t *testing.T) {
	test_s := `## Purity and the Mob

It happens in some places that live off tourism: at the restaurant where you're
having dinner, your friendly waiter suggests that if you're looking to party,
his brother in law has the best club in town. Ah, not in the mood for partying?`
	result := wraptext(80, test_s)
	if !same_text(test_s, result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
}

func Test_lines_at_limit_work(t *testing.T) {
	test_s := "12 12 12 2334 233"
	result := wraptext(8, test_s)
	lines := strings.Split(result, "\n")
	if !same_text(test_s, result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(lines) != 2 {
		t.Errorf("Missing newlines in: %s", result)
	}

}

func Test_wrap_single_line(t *testing.T) {
	test_s := "12 12 12 233 23 3 345 34"
	result := wraptext(10, test_s)
	lines := strings.Split(result, "\n")
	if !same_text(test_s, result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(lines) != 3 {
		t.Errorf("Missing newlines in: %s", result)
	}
}

func Test_wrap_lines(t *testing.T) {
	test_s := "My name is Wile E. Coyote"
	result := wraptext(15, test_s)
	if !same_text(test_s, result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(strings.Split(result, "\n")) != 2 {
		t.Errorf("Bad wrap:\n%s", result)
	}
}

func Test_wrap_respects_one_blank_line(t *testing.T) {
	test_s := "My name is Wile E.\nCoyote.\n\nGenius"
	result := wraptext(15, test_s)
	if !same_text(test_s, result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(strings.Split(result, "\n")) != 5 {
		t.Errorf("Didn't respect blank line:\n%s", result)
	}
}

func Test_wrap_respects_two_blank_lines(t *testing.T) {
	test_s := "My name is E.\n\n\nCoyote"
	result := wraptext(15, test_s)
	if !same_text(test_s, result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(strings.Split(result, "\n")) != 4 {
		t.Errorf("Didn't respect two blank lines:\n%s", result)
	}
}

func Test_wrap_respects_full_stops(t *testing.T) {
	test_s := "My name is Wile E.\nCoyote"
	result := wraptext(15, test_s)
	if !same_text(test_s, result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(strings.Split(result, "\n")) != 3 {
		t.Errorf("Didn't respect full stop:\n%s", result)
	}
}

func Test_wrap_does_not_add_extra_breaks(t *testing.T) {
	test_s := "My name is."
	result := wraptext(15, test_s)
	if !same_text(test_s, result) {
		t.Errorf("Content was altered:\n%s\n%s", test_s, result)
	}
	if len(strings.Split(result, "\n")) != 1 {
		t.Errorf("Didn't respect period-break:\n%s", result)
	}
}
