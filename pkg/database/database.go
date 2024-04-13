package database

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func CreateTempFolder(tempDirPath string, tempFolderPattern string, goroutineID int) string {
	//не уверен что темп не перезапишет самого себя, поэтому добавил айди горутины
	tempFolderGoroutinePattern := tempFolderPattern + "-" + strconv.Itoa(goroutineID) + "-"
	tempDirPathResult, err := os.MkdirTemp(tempDirPath, tempFolderGoroutinePattern)
	if err != nil {
		fmt.Println("Error from create temp dir:", err)
		return ""
	}

	return tempDirPathResult
}

func FoundTemp(tempDirPath string, tempFolderPattern string) map[string][]string {
	//есть вариант хранить имена папок в отдельном файле, но не хочется хламить

	tempFolders := map[string][]string{} //key tempfolder, value filenames

	folders, err := os.ReadDir(tempDirPath)
	if err != nil {
		fmt.Println("Error of reading temp dir:", err)
		return tempFolders
	}
	for _, folder := range folders {
		if folder.IsDir() && strings.HasPrefix(folder.Name(), tempFolderPattern) {
			tempFolders[folder.Name()] = FoundTempFiles(tempDirPath + "/" + folder.Name())
		}
	}
	fmt.Println(len(tempFolders))
	return tempFolders
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
	err := writeToFile(filePath, data, false)

	if err != nil {

		//try to write data in emergency db file
		err = writeToFile(eDBpath, data, true)

		if err != nil {
			fmt.Println("Все потеряно, данные не записаны:", err)
			return err
		}
	}
	return nil
}

func writeToFile(filePath string, data []byte, emergency bool) error {

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
