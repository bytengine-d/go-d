//go:build windows

package lang

import (
	"os"
	"syscall"
	"time"
)

func FileCreateTime(fi os.FileInfo) time.Time {
	fileAttr := fi.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, fileAttr.CreationTime.Nanoseconds())
}

func FileLastModifiedTime(fi os.FileInfo) time.Time {
	fileAttr := fi.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, fileAttr.LastWriteTime.Nanoseconds())
}

func FileLastAccessTime(fi os.FileInfo) time.Time {
	fileAttr := fi.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, fileAttr.LastAccessTime.Nanoseconds())
}
