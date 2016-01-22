package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func ExampleWraptext() {
	text := "My name is Wile E. Coyote, genius.  I'm not here selling"
	reader := strings.NewReader(text)
	scanner := bufio.NewScanner(reader)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	Rapper{maxcols: 20}.Wraptext(scanner, writer)
	fmt.Println(buffer.String())
	// Output: My name is Wile E.
	//Coyote, genius.  I'm
	//not here selling
}
