package main

import (
	"fmt"
	"sync"
	"time"
)

// What is the result ?
// A: cannot compile
// B: [output] 1
// C: program hanging
// D: panic

var mu sync.RWMutex
var count int

func main() {
	go A()
	time.Sleep(2 * time.Second)
	mu.Lock()
	defer mu.Unlock()
	count++
	fmt.Println(count)
}
func A() {
	mu.RLock()
	defer mu.RUnlock()
	B()
}
func B() {
	time.Sleep(5 * time.Second)
	C()
}
func C() {
	mu.RLock()
	defer mu.RUnlock()
}

//
// === ANS ===
//
// D) panic
// line 21 will acquire the lock earlier than line 36 (func C)
// main goroutine and goroutine got deadlock
