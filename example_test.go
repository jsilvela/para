package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func ExampleWrapper_Wraptext_splits() {
	text := "My name is Wile E. Coyote, genius.  I'm not here selling"
	reader := strings.NewReader(text)
	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	Wrapper{maxCols: 20}.Wraptext(scanner, writer)
	fmt.Println(buffer.String())
	// Output:
	// My name is Wile E.
	// Coyote, genius.  I'm
	// not here selling
}

func ExampleWrapper_Wraptext_compacts() {
	text := "My name is Wile E\nCoyote, genius. I'm\nnot here selling"
	reader := strings.NewReader(text)
	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	Wrapper{maxCols: 80}.Wraptext(scanner, writer)
	fmt.Println(buffer.String())
	// Output:
	// My name is Wile E Coyote, genius. I'm not here selling
}
