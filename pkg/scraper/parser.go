package scraper

import (
	"encoding/json"
	"fmt"
	"myapp/pkg/words"
	"time"
)

type ResponseData struct {
	Alt string `json:"alt"`
	Img string `json:"img"`
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

func decodeResponse(data []byte) (ResponseData, error) {
	var result ResponseData
	if err := json.Unmarshal(data, &result); err != nil {
		return ResponseData{}, err
	}
	return result, nil
}

func DecodeFileData(fileData []byte) ScrapeResult {
	var data ScrapeResult
	if err := json.Unmarshal(fileData, &data); err != nil {
		return ScrapeResult{
			Data:   map[int]ParsedData{},
			BadIDs: map[int]int{},
		}
	}
	return data
}

func codeData(bytesData ScrapeResult) []byte {
	// Code to JSON
	data, err := json.MarshalIndent(bytesData, "", "\t")
	if err != nil {
		return nil
	}
	return data
}

func DataToPrint(data map[int]ParsedData) string {
	bytes, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		fmt.Println("Ошибка при форматировании JSON:", err)
	}
	return string(bytes)
}
func responseParser(data []byte) (ParsedData, error) {

	result := ParsedData{}
	if data == nil {
		return ParsedData{}, nil
	}

	jsonData, err := decodeResponse(data)
	result.Keywords = words.StringNormalization(jsonData.Alt)
	result.Url = jsonData.Img

	return result, err
}
