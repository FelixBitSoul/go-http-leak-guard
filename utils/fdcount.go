package utils

import (
	"os"
	"runtime"
)

// GetFDCount returns the number of open file descriptors for the current process.
// Returns -1 if the count cannot be determined (e.g., on non-Linux systems or if access is denied).
func GetFDCount() int {
	if runtime.GOOS == "linux" {
		// Try to read the standard Linux process file descriptor directory
		files, err := os.ReadDir("/proc/self/fd")
		if err != nil {
			// Return -1 to distinguish "failed to get" from "count is 0"
			return -1
		}
		// Subtract 1 because ReadDir itself opens the directory and occupies one FD
		return len(files) - 1
	}

	// For non-Linux systems, we can't easily count handles
	return -1
}
