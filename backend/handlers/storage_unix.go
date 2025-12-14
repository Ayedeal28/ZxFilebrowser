//go:build !windows

package handlers

import (
	"syscall"
)

// getDiskStats returns total and free disk space for Unix-like systems
func getDiskStats(path string) (total, free uint64, err error) {
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return 0, 0, err
	}

	// Calculate sizes
	total = stat.Blocks * uint64(stat.Bsize)
	free = stat.Bfree * uint64(stat.Bsize)

	return total, free, nil
}
