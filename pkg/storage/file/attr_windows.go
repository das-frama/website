// +build windows

package file

import (
	"os"
	"syscall"
	"time"
)

func timeCreation(path string) time.Time {
	info, _ := os.Lstat(path)
	sys := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, sys.CreationTime.Nanoseconds())
}
