package singleinstance

import (
	"log"
	"net"
	"sync"
	"time"
)

const activateAddr = "127.0.0.1:47891"

var (
	mu               sync.Mutex
	instanceListener net.Listener
	isOwner          bool
	activateFn       func()
	pendingShow      bool
)

// EnsureFirstInstance binds the local activation port. A second process cannot
// bind it and will ask the running instance to show its window instead.
func EnsureFirstInstance() bool {
	ln, err := net.Listen("tcp", activateAddr)
	if err != nil {
		requestActivation()
		return false
	}

	mu.Lock()
	instanceListener = ln
	isOwner = true
	mu.Unlock()

	go serveActivationRequests(ln)
	return true
}

// IsOwner reports whether this process holds the single-instance lock.
func IsOwner() bool {
	mu.Lock()
	defer mu.Unlock()
	return isOwner
}

// Release closes the activation listener for this process.
func Release() {
	mu.Lock()
	defer mu.Unlock()
	if instanceListener != nil {
		_ = instanceListener.Close()
		instanceListener = nil
	}
	isOwner = false
}

// SetActivationHandler registers the callback used when another launch requests
// the running instance to show its window (for example from the systray path).
func SetActivationHandler(onActivate func()) {
	mu.Lock()
	activateFn = onActivate
	pending := pendingShow
	pendingShow = false
	mu.Unlock()

	if pending && onActivate != nil {
		onActivate()
	}
}

func serveActivationRequests(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		buf := make([]byte, 16)
		_, _ = conn.Read(buf)
		_ = conn.Close()
		notifyActivate()
	}
}

func notifyActivate() {
	mu.Lock()
	fn := activateFn
	if fn == nil {
		pendingShow = true
		mu.Unlock()
		return
	}
	mu.Unlock()
	fn()
}

func requestActivation() {
	payload := []byte("show")
	for attempt := 0; attempt < 40; attempt++ {
		conn, err := net.DialTimeout("tcp", activateAddr, 100*time.Millisecond)
		if err == nil {
			_, _ = conn.Write(payload)
			_ = conn.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	log.Printf("singleinstance: could not reach running instance")
}
