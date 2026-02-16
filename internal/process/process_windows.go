package process

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows"
)

// waitForPIDPlatform uses the Windows API to wait for process exit.
func waitForPIDPlatform(pid int) (time.Duration, error) {
	start := time.Now()

	handle, err := windows.OpenProcess(
		windows.SYNCHRONIZE|windows.PROCESS_QUERY_LIMITED_INFORMATION,
		false,
		uint32(pid),
	)
	if err != nil {
		return 0, fmt.Errorf("cannot open process %d: %w", pid, err)
	}
	defer windows.CloseHandle(handle)

	event, err := windows.WaitForSingleObject(handle, windows.INFINITE)
	if err != nil {
		return 0, fmt.Errorf("waiting for process %d: %w", pid, err)
	}
	if event != windows.WAIT_OBJECT_0 {
		return 0, fmt.Errorf("unexpected wait result for process %d: %d", pid, event)
	}

	return time.Since(start), nil
}
