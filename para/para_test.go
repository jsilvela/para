package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func Test_wrapcol(t *testing.T) {
	wrapped := wrapcol(10, "12 12 12 233 23 3 345 34")
	lines := strings.Split(wrapped, "\n")
	if len(lines) != 3 {
		t.Errorf("Missing newlines in: %s", wrapped)
	}
}

func wrap_wraptext(col int, in string) string {
	reader := strings.NewReader(in)
	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	wraptext(col, scanner, writer)
	return buffer.String()
}

func Test_wraptext_adds_newline_on_flush(t *testing.T) {
	test_s := "My name is Wile E. Coyote"
	result := wrap_wraptext(15, test_s)
	if len(strings.Split(result, "\n")) != 3 {
		t.Errorf("Bad wrap:\n%s", result)
	}
}

func Test_wraptext_respects_paragraphs(t *testing.T) {
	test_s := "My name is Wile E.\nCoyote"
	result := wrap_wraptext(15, test_s)
	if len(strings.Split(result, "\n")) != 4 {
		t.Errorf("Didn't respect period-break:\n%s", result)
	}
}
