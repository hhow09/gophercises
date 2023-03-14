package main

import (
	"fmt"
	"sync"
)

// What is the result ?
// A: cannot compile
// B: [output] 1, 1
// C: [output] 1, 2
// D: panic

type MyMutex struct {
	count int
	sync.Mutex
}

func main() {
	var mu MyMutex
	mu.Lock()
	var mu2 = mu
	mu.count++
	mu.Unlock()
	mu2.Lock()
	mu2.count++
	mu2.Unlock()
	fmt.Println(mu.count, mu2.count)
}

//
// === ANS ===
//
// D) panic
// we should not copy the mutex
