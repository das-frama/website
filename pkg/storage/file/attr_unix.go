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
	sec, nano := sys.Ctim.Unix()
	return time.Unix(sec, nano)
}
