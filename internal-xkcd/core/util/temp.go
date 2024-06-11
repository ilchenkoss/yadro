package util

import (
	"fmt"
	"myapp/internal-xkcd/config"
	"os"
	"strconv"
	"strings"
)

type Temper struct {
	TempDir     string
	TempFiles   map[string]bool
	TempedIDs   map[int]bool
	MaxTempedID int
	Config      *config.TempConfig
	Fs          FileSystem
}

func NewTemper(cfg *config.TempConfig, fs FileSystem) (*Temper, error) {

	// Check directory exists
	_, statErr := fs.Stat(cfg.TempDir)
	if fs.IsNotExist(statErr) {
		// Directory not exist
		mkDirErr := fs.Mkdir(cfg.TempDir, 0777)
		if mkDirErr != nil {
			return nil, mkDirErr
		}
	}
	if statErr != nil && !fs.IsNotExist(statErr) {
		return nil, statErr
	}

	tempFiles := make(map[string]bool)
	tempedIDs := make(map[int]bool)
	maxTempedID := 0

	existFiles, err := fs.ReadDir(cfg.TempDir)
	if err != nil {
		return nil, err
	}
	for _, file := range existFiles {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		tempFiles[fileName] = true
		fileNameWithoutPattern := strings.Split(fileName, cfg.TempFilePattern)[1]
		fileIDString := strings.Split(fileNameWithoutPattern, "-")[0]
		fileID, strToIntErr := strconv.Atoi(fileIDString)

		if strToIntErr != nil {
			continue
		}

		tempedIDs[fileID] = true

		if fileID > maxTempedID {
			maxTempedID = fileID
		}
	}

	return &Temper{
		TempDir:     cfg.TempDir,
		TempFiles:   tempFiles,
		TempedIDs:   tempedIDs,
		MaxTempedID: maxTempedID,
		Config:      cfg,
		Fs:          fs}, nil
}

func (t *Temper) ReadTempFile(filePath string) []byte {
	fileData, err := t.Fs.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return fileData
}
func (t *Temper) SaveTempDataByID(data []byte, ID int) error {
	tempFile, createTempErr := t.Fs.CreateTemp(t.TempDir, fmt.Sprintf("%s%d-", t.Config.TempFilePattern, ID))
	if createTempErr != nil {
		return fmt.Errorf("error creating temporary file: %v", createTempErr)
	}
	defer func(tempFile *os.File) {
		err := tempFile.Close()
		if err != nil {
			//nothing
			return
		}
	}(tempFile)

	if _, writeErr := tempFile.Write(data); writeErr != nil {
		return fmt.Errorf("error writing data to temporary file: %v", writeErr)
	}

	return nil
}

func (t *Temper) RemoveTemp() error {
	return t.Fs.RemoveAll(t.TempDir)
}

type FileSystem interface {
	Stat(name string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	ReadDir(name string) ([]os.DirEntry, error)
	ReadFile(name string) ([]byte, error)
	CreateTemp(dir, pattern string) (*os.File, error)
	RemoveAll(path string) error
	IsNotExist(err error) bool
}

type OSFileSystem struct{}

func (OSFileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OSFileSystem) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (OSFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

func (OSFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (OSFileSystem) CreateTemp(dir, pattern string) (*os.File, error) {
	return os.CreateTemp(dir, pattern)
}

func (OSFileSystem) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
