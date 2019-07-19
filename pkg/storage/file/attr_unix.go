// +build !windows

package file

import (
	"os"
	"syscall"
	"time"
)

func timeFromFile(path string) (time.Time, time.Time, error) {
	info, _ := os.Lstat(path)
	sys := info.Sys().(*syscall.Stat_t)
	sec1, nano1 := sys.Ctim.Unix()
	sec2, nano2 := sys.Mtim.Unix()
	return time.Unix(sec1, nano1), time.Unix(sec2, nano2), nil
}
