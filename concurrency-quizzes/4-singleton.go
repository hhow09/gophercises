package main

import (
	"sync"
)

// double check singleton
// What is the result ?
// A: cannot compile
// B: can compile, correct singleton
// C: can compile, wrong singleton
// D: can compile, will panic

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if o.done == 1 {
		return
	}
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		o.done = 1
		f()
	}
}

//
// === ANS ===
//
// C) can compile, wrong singleton
// o.done = 1 might not be atomic operation
//
// [correct way] with atomic
// type Once struct {
// 	m    sync.Mutex
// 	done uint32
// }

// func (o *Once) Do(f func()) {
// 	if atomic.LoadUint32(&o.done) == 1 {
// 		return
// 	}
// 	o.m.Lock()
// 	defer o.m.Unlock()
// 	if o.done == 0 {
// 		atomic.StoreUint32(&o.done, 1)
// 		f()
// 	}
// }
