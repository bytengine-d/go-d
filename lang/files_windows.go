//go:build windows

package lang

import (
	"os"
	"syscall"
	"time"
)

func filetimeToTime(fi *syscall.Timespec) time.Time {
	return time.Unix(fi.Nanoseconds()/1e9, 0)
}

func FileCreateTime(fi os.FileInfo) time.Time {
	fileAttr := fi.Sys().(*syscall.Win32FileAttributeData)
	return filetimeToTime(fileAttr.CreationTime)
}

func FileLastModifiedTime(fi os.FileInfo) time.Time {
	fileAttr := fi.Sys().(*syscall.Win32FileAttributeData)
	return filetimeToTime(fileAttr.LastWriteTime)
}

func FileLastAccessTime(fi os.FileInfo) time.Time {
	fileAttr := fi.Sys().(*syscall.Win32FileAttributeData)
	return filetimeToTime(fileAttr.LastAccessTime)
}
