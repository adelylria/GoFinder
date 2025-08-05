// hotkey/hotkey.go
package hotkey

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lhotkey
#include "hotkey.h"
#include <stdlib.h>
*/
import "C"

import (
	"log"
	"sync"
)

var wg sync.WaitGroup

func StartListening() {
	wg.Add(1)
	go func() {
		defer wg.Done()
		C.ListenHotkey()
	}()
	log.Println("Started hotkey listener")
}

func StopListening() {
	C.StopHotkey()
	wg.Wait()
	log.Println("Stopped hotkey listener")
}
