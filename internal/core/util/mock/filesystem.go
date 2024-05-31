package mock

import (
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing/fstest"
	"time"
)

type MockOSFileSystem struct {
	Mfs fstest.MapFS
	Mu  sync.Mutex
}

func (mosfs *MockOSFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (mosfs *MockOSFileSystem) Mkdir(name string, perm fs.FileMode) error {
	mosfs.Mu.Lock()
	defer mosfs.Mu.Unlock()
	mosfs.Mfs[name] = &fstest.MapFile{Mode: fs.ModeDir}
	return nil
}

func (mosfs *MockOSFileSystem) CreateTemp(dir, pattern string) (*os.File, error) {
	mosfs.Mu.Lock()
	defer mosfs.Mu.Unlock()
	filename := fmt.Sprintf("%s/%s%d", dir, pattern, rand.Intn(100))
	mosfs.Mfs[filename] = &fstest.MapFile{}

	f, err := os.CreateTemp("", "mock")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (mosfs *MockOSFileSystem) RemoveAll(path string) error {
	mosfs.Mu.Lock()
	defer mosfs.Mu.Unlock()
	for p := range mosfs.Mfs {
		if strings.Contains(p, path) {
			delete(mosfs.Mfs, p)
		}
	}
	return nil
}

func (mosfs *MockOSFileSystem) Stat(name string) (fs.FileInfo, error) {
	mosfs.Mu.Lock()
	defer mosfs.Mu.Unlock()
	if file, ok := mosfs.Mfs[name]; ok {
		return mockFileInfo{name: name, MapFile: file}, nil
	}
	return nil, os.ErrNotExist
}

func (mosfs *MockOSFileSystem) ReadDir(dirname string) ([]fs.DirEntry, error) {
	mosfs.Mu.Lock()
	defer mosfs.Mu.Unlock()
	var entries []fs.DirEntry
	for name, file := range mosfs.Mfs {
		if name == dirname {
			continue
		}
		if len(name) > len(dirname) && name[:len(dirname)] == dirname && name[len(dirname)] == '/' {
			entries = append(entries, mockDirEntry{name: name[len(dirname)+1:], MapFile: file})
		}
	}
	return entries, nil
}

func (mosfs *MockOSFileSystem) ReadFile(filename string) ([]byte, error) {
	mosfs.Mu.Lock()
	defer mosfs.Mu.Unlock()
	if file, ok := mosfs.Mfs[filename]; ok {
		return file.Data, nil
	}
	return nil, os.ErrNotExist
}

type mockFileInfo struct {
	name string
	*fstest.MapFile
}

func (mfi mockFileInfo) Name() string       { return mfi.name }
func (mfi mockFileInfo) Size() int64        { return int64(len(mfi.Data)) }
func (mfi mockFileInfo) Mode() fs.FileMode  { return mfi.MapFile.Mode }
func (mfi mockFileInfo) ModTime() time.Time { return time.Now() }
func (mfi mockFileInfo) IsDir() bool        { return mfi.Mode().IsDir() }
func (mfi mockFileInfo) Sys() interface{}   { return nil }

type mockDirEntry struct {
	name string
	*fstest.MapFile
}

func (mde mockDirEntry) Name() string               { return mde.name }
func (mde mockDirEntry) IsDir() bool                { return mde.MapFile.Mode.IsDir() }
func (mde mockDirEntry) Type() fs.FileMode          { return mde.MapFile.Mode.Type() }
func (mde mockDirEntry) Info() (fs.FileInfo, error) { return mockFileInfo(mde), nil }
