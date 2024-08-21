package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func uniq(input io.Reader, output io.Writer) error {
	in := bufio.NewScanner(input)
	var prev string
	for in.Scan() {
		txt := in.Text()
		if txt == prev && prev <= txt {
			continue
		} else if prev > txt {
			return fmt.Errorf("Unsorted data")
		} else {
			fmt.Fprintln(output, txt)
			prev = txt
		}
	}

	return nil
}

func main() {
	err := uniq(os.Stdin, os.Stdout)
	if err != nil {
		panic(err.Error())
	}
}
