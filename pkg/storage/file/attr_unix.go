// +build !windows

package file

import (
	"os"
	"syscall"
	"time"
)

func timeCreation(path string) time.Time {
	info, _ := os.Lstat(path)
	sys := info.Sys().(*syscall.Stat_t)
	return time.Unix(0, sys.Ctimespec.Nsec)
}
