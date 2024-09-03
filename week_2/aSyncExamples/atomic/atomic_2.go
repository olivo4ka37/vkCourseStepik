package main

import (
	"fmt"
	"sync/atomic" // atomic_2.go
	"time"
)

var totalOperations2 int32

func inc2() {
	atomic.AddInt32(&totalOperations2, 1) // автомарно
}

func main() {
	for i := 0; i < 1000; i++ {
		go inc2()
	}
	time.Sleep(2 * time.Millisecond)
	fmt.Println("total operation = ", totalOperations2)
}
