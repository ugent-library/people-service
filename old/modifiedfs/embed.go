package modifiedfs

import (
	"io/fs"
	"net/http"
	"os"
	"time"
)

// add header Last-Modified for embedded files.
// stolen from https://github.com/golang/go/issues/44854#issuecomment-808906568
type staticFSWrapper struct {
	http.FileSystem
	FixedModTime time.Time
}

func (f *staticFSWrapper) Open(name string) (http.File, error) {
	file, err := f.FileSystem.Open(name)
	return &staticFileWrapper{File: file, fixedModTime: f.FixedModTime}, err
}

type staticFileWrapper struct {
	http.File
	fixedModTime time.Time
}

func (f *staticFileWrapper) Stat() (os.FileInfo, error) {
	fileInfo, err := f.File.Stat()
	return &staticFileInfoWrapper{FileInfo: fileInfo, fixedModTime: f.fixedModTime}, err
}

type staticFileInfoWrapper struct {
	os.FileInfo
	fixedModTime time.Time
}

func (f *staticFileInfoWrapper) ModTime() time.Time {
	return f.fixedModTime
}

// TODO: move this to separate repo?
func FSWithStatModified(f fs.FS, t time.Time) *staticFSWrapper {
	return &staticFSWrapper{
		FileSystem:   http.FS(f),
		FixedModTime: t,
	}
}
