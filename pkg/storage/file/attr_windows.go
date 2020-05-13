// +build windows

package file

import (
	"os"
	"syscall"
	"time"
)

func timeFromFile(path string) (time.Time, time.Time, error) {
	info, _ := os.Lstat(path)
	sys := info.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, sys.CreationTime.Nanoseconds()), time.Unix(0, sys.LastWriteTime.Nanoseconds()), nil
}
