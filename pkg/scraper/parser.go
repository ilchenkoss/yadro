package scraper

import (
	"encoding/json"
	"fmt"
	"myapp/pkg/database"
	"myapp/pkg/words"
)

type ResponseData struct {
	Alt string `json:"alt"`
	Img string `json:"img"`
}

func decodeResponse(data []byte) (ResponseData, error) {

	var result ResponseData
	if err := json.Unmarshal(data, &result); err != nil {
		return ResponseData{}, err
	}
	return result, nil
}

func DataToPrint(data map[int]database.ParsedData) string {
	bytes, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		fmt.Println("Ошибка при форматировании JSON:", err)
	}
	return string(bytes)
}
func responseParser(data []byte) (database.ParsedData, error) {

	result := database.ParsedData{}
	if data == nil {
		return database.ParsedData{}, nil
	}

	jsonData, err := decodeResponse(data)
	result.Keywords = words.StringNormalization(jsonData.Alt)
	result.Url = jsonData.Img

	return result, err
}
