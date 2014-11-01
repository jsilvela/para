package main

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"unicode"
)

const wrap_col = 10

func main() {
	var out []byte
	s := bufio.NewReader(os.Stdin)
	for line, err := s.ReadBytes('\n'); line!=nil; line, err = s.ReadBytes('\n') {
		j, last_white, line_chars := 0, 0, 0
		out = make([]byte, 2*len(line))
		for _, c := range(line) {
			out[j] = c
			if unicode.IsSpace(rune(c)) {
				last_white = j
			}
			if line_chars > wrap_col {
				out[last_white] = '\n'
				line_chars = j - last_white
			}
			j++
			line_chars++
		}
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(string(out))
		}
	}
}