package main

import (
	"fmt"
	"sync"
)

// What is the result ?
// A: cannot compile
// B: output main --> A --> B --> C
// C: output main
// D: panic

var mu sync.Mutex
var chain string

func main() {
	chain = "main"
	A()
	fmt.Println(chain)
}
func A() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + " --> A"
	B()
}
func B() {
	chain = chain + " --> B"
	C()
}
func C() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + " --> C"
}

// === ANS ===
//
// D) panic
// by deadlock
