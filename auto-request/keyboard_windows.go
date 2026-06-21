//go:build windows
package main

import (
	"syscall"
)

var (
	msvcrt    = syscall.NewLazyDLL("msvcrt.dll")
	procKbhit = msvcrt.NewProc("_kbhit")
	procGetch = msvcrt.NewProc("_getch")
)

// IsEscPressed checks if the ESC key has been pressed without blocking the thread.
func IsEscPressed() bool {
	r, _, _ := procKbhit.Call()
	if r != 0 {
		ch, _, _ := procGetch.Call()
		return ch == 27 // 27 is ASCII value of ESC key
	}
	return false
}
