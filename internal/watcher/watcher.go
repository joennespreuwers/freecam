package watcher

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// KillResult holds information about a successful kill.
type KillResult struct {
	ProcessName string
	PID         int
	KilledAt    time.Time
}

// FindAndKill finds all running processes matching name and kills them.
// Returns one KillResult per killed process.
func FindAndKill(processName string) ([]KillResult, error) {
	out, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return nil, fmt.Errorf("ps aux: %w", err)
	}

	var results []KillResult
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if !strings.Contains(line, processName) {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
			continue
		}
		results = append(results, KillResult{
			ProcessName: processName,
			PID:         pid,
			KilledAt:    time.Now(),
		})
	}
	return results, nil
}
