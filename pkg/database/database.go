package database

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func ReadDatabase(dbpath string) ScrapeResult {

	dataBytes, databaseErr := readBytesFromFile(dbpath) //data

	if databaseErr != nil { //error reading database

		if os.IsNotExist(databaseErr) { //if database not exists
			return ScrapeResult{
				Data:   map[int]ParsedData{},
				BadIDs: map[int]int{},
			}
		}
		panic(databaseErr) //read error and file exists
	}

	dbData, decodeErr := DecodeData(dataBytes) //try decode db

	if decodeErr != nil { //decode err
		panic(decodeErr)
	}

	return dbData

}

type ScrapeResult struct {
	Data      map[int]ParsedData `json:"data"`
	BadIDs    map[int]int        `json:"badIDs"`
	Timestamp time.Time          `json:"timestamp"`
}
type ParsedData struct {
	Keywords []string `json:"keywords"`
	Url      string   `json:"url"`
}

func readBytesFromFile(filePath string) ([]byte, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileData, nil
}

func DecodeData(fileData []byte) (ScrapeResult, error) {
	var data ScrapeResult
	if err := json.Unmarshal(fileData, &data); err != nil {
		return ScrapeResult{}, err
	}
	return data, nil
}

func WriteData(filePath string, eDBpath string, data ScrapeResult) {

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

func writeToFile(filePath string, data ScrapeResult, emergency bool) error {

	// Open or create database file
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Code to JSON
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Write to database file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	if emergency {
		fmt.Printf("Данные записаны в аварийную базу данных: ", filePath)
	}

	return nil
}
