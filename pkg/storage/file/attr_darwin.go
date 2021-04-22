// +build darwin

package file

import (
	"os"
	"syscall"
	"time"
)

func timeFromFile(path string) (time.Time, time.Time, error) {
	info, _ := os.Stat(path)
	sys := info.Sys().(*syscall.Stat_t)
	sec1, nano1 := sys.Ctimespec.Unix()
	sec2, nano2 := sys.Mtimespec.Unix()
	return time.Unix(sec1, nano1), time.Unix(sec2, nano2), nil
}
