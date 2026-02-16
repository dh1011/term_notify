package process

import "time"

// WaitForPID blocks until the process with the given PID exits.
// Returns the time spent waiting.
func WaitForPID(pid int) (time.Duration, error) {
	return waitForPIDPlatform(pid)
}
