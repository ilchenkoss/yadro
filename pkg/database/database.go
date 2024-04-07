package database

import (
	"fmt"
	"os"
)

func ReadBytesFromFile(filePath string) []byte {

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return fileData
}

func WriteData(filePath string, eDBpath string, data []byte) {

	//try to write data in main db file
	err := writeToFile(filePath, data, false)

	if err != nil {

		//try to write data in emergency db file
		err = writeToFile(eDBpath, data, true)

		if err != nil {
			fmt.Println("Все потеряно, данные не записаны:", err)
			return
		}
	}
	return
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
