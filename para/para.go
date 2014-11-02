package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode"
)

const default_wrap_col = 90

func main() {
	var wrap_col int
	if len(os.Args) < 2 {
		wrap_col = default_wrap_col
	} else {
		num, err := strconv.ParseInt(os.Args[1], 10, 0)
		if err != nil {
			log.Fatal(err)
		} else {
			wrap_col = int(num)
		}
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		j, last_white, line_chars := 0, 0, 0
		out := make([]rune, 2*len(line))
		for _, c := range line {
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
		fmt.Println(string(out))
	}
}
