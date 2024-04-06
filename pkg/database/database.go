package database

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func DataToPrint(data map[int]ParsedData) string {
	bytes, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		fmt.Println("Ошибка при форматировании JSON:", err)
	}
	return string(bytes)
}

type ParsedData struct {
	Keywords []string `json:"keywords"`
	Url      string   `json:"url"`
}
type ScrapeResult struct {
	Data      map[int]ParsedData `json:"data"`
	BadIDs    map[int]int        `json:"badIDs"`
	Timestamp time.Time          `json:"timestamp"`
}

func ReadData(filePath string) ([]byte, error) {
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
func WriteData(filePath string, data ScrapeResult) error {
	// create file or exist
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
	}
	defer file.Close()

	//code to json
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println("Ошибка при кодировании в JSON:", err)
		return err
	}

	//write to file
	err = os.WriteFile("./pkg/database/database.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Ошибка при записи данных в файл:", err)
		return err
	}
	return nil
}
