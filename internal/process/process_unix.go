//go:build !windows

package process

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// waitForPIDPlatform polls by sending signal 0 to the process on Unix systems.
func waitForPIDPlatform(pid int) (time.Duration, error) {
	start := time.Now()

	proc, err := os.FindProcess(pid)
	if err != nil {
		return 0, fmt.Errorf("cannot find process %d: %w", pid, err)
	}

	for {
		err := proc.Signal(syscall.Signal(0))
		if err != nil {
			return time.Since(start), nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}
