//go:build amd64

package lang

import (
	"os"
	"syscall"
	"time"
)

func timespecToTime(timespec syscall.Timespec) time.Time {
	return time.Unix(timespec.Sec, timespec.Nsec)
}

func FileCreateTime(fi os.FileInfo) time.Time {
	fileAttr := fi.Sys().(*syscall.Stat_t)
	return timespecToTime(fileAttr.Ctimespec)
}

func FileLastModifiedTime(fi os.FileInfo) time.Time {
	fileAttr := fi.Sys().(*syscall.Stat_t)
	return timespecToTime(fileAttr.Mtimespec)
}

func FileLastAccessTime(fi os.FileInfo) time.Time {
	fileAttr := fi.Sys().(*syscall.Stat_t)
	return timespecToTime(fileAttr.Atimespec)
}
