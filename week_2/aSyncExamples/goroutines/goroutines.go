package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const (
	iterationsNum2 = 6
	goroutinesNum2 = 6
)

func doWork(th int) {
	for j := 0; j < iterationsNum2; j++ {
		fmt.Printf(formatWork2(th, j))
		// time.Sleep(time.Millisecond)
	}
}

func main() {
	for i := 0; i < goroutinesNum2; i++ {
		go doWork(i)
	}
	fmt.Scanln()
}

func formatWork2(in, j int) string {
	return fmt.Sprintln(strings.Repeat("  ", in), "█",
		strings.Repeat("  ", goroutinesNum2-in),
		"th", in,
		"iter", j, strings.Repeat("■", j))
}

func imports() {
	fmt.Println(time.Millisecond, runtime.NumCPU())
}
