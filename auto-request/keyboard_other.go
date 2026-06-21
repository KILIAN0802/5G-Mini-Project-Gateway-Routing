//go:build !windows
package main

// IsEscPressed is a fallback stub for non-Windows platforms.
func IsEscPressed() bool {
	return false
}
