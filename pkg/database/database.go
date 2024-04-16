package database

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CreateTempFolder(tempDirPath string, tempFolderPattern string) string {
	tempDirPathResult, err := os.MkdirTemp(tempDirPath, tempFolderPattern)
	if err != nil {
		fmt.Println("Error from create temp dir:", err)
		return ""
	}
	return tempDirPathResult
}

func SaveTemp(data []byte, tempFolderPath string, tempFilePattern string, ID int) {
	os.CreateTemp(tempFolderPath, fmt.Sprintf("%s%d-", tempFilePattern, ID))
}

type Temp struct {
	TempPaths map[string][]string
	TempIDs   []int
}

func FoundTemp(tempDirPath string, tempFolderPattern string, tempFilePattern string) Temp {
	var tempIDs []int
	tempPaths := make(map[string][]string)
	temp := Temp{
		TempIDs:   tempIDs,
		TempPaths: tempPaths,
	}

	folders, err := os.ReadDir(tempDirPath)
	if err != nil {
		fmt.Println("Error of reading temp dir:", err)
		return temp
	}
	for _, folder := range folders {
		if folder.IsDir() && strings.HasPrefix(folder.Name(), tempFolderPattern) {
			tempFolderPath := fmt.Sprintf("%s%s", tempDirPath, folder.Name())
			tempFiles := FoundTempFiles(tempFolderPath)
			temp.TempPaths[tempFolderPath] = tempFiles

			for _, tempFileName := range tempFiles {
				fileNameWithoutPattern := strings.Split(tempFileName, tempFilePattern)[1]
				fileIDString := strings.Split(fileNameWithoutPattern, "-")[0]
				fileID, strToIntErr := strconv.Atoi(fileIDString)
				if strToIntErr == nil {
					temp.TempIDs = append(temp.TempIDs, fileID)
				}
			}
		}
	}
	return temp
}

func FoundTempFiles(tempDirPath string) []string {
	var tempFiles []string
	files, err := os.ReadDir(tempDirPath)
	if err != nil {
		fmt.Println("Error of reading temp files:", err)
		return tempFiles
	}
	for _, file := range files {
		tempFiles = append(tempFiles, file.Name())
	}
	return tempFiles
}

func ReadBytesFromFile(filePath string) []byte {

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return fileData
}

func WriteData(filePath string, eDBpath string, data []byte) error {

	//try to write data in main db file
	err := WriteToFile(filePath, data, false)

	if err != nil {

		//try to write data in emergency db file
		err = WriteToFile(eDBpath, data, true)

		if err != nil {
			fmt.Println("Все потеряно, данные не записаны:", err)
			return err
		}
	}
	return nil
}

func WriteToFile(filePath string, data []byte, emergency bool) error {

	// Open or create database file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write to database file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	if emergency {
		fmt.Printf("Данные записаны в аварийную базу данных: %s", filePath)
	}

	return nil
}
