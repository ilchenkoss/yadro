package util

import (
	"fmt"
	"myapp/internal/config"
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
}

func NewTemper(cfg *config.TempConfig) (*Temper, error) {

	// Check directory exists
	_, statErr := os.Stat(cfg.TempDir)
	if os.IsNotExist(statErr) {
		// Directory not exist
		mkDirErr := os.Mkdir(cfg.TempDir, 0777)
		if mkDirErr != nil {
			return nil, mkDirErr
		}
	}
	if statErr != nil && !os.IsNotExist(statErr) {
		return nil, statErr
	}

	tempFiles := make(map[string]bool)
	tempedIDs := make(map[int]bool)
	maxTempedID := 0

	existFiles, err := os.ReadDir(cfg.TempDir)
	if err != nil {
		return nil, err
	}
	for _, file := range existFiles {
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
		Config:      cfg}, nil
}

func (t *Temper) ReadTempFile(filePath string) []byte {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return fileData
}

func (t *Temper) SaveTempDataByID(data []byte, ID int) error {

	tempFile, createTempErr := os.CreateTemp(t.TempDir, fmt.Sprintf("%s%d-", t.Config.TempFilePattern, ID))
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
	err := os.RemoveAll(t.TempDir)
	if err != nil {
		return err
	}
	return nil
}
