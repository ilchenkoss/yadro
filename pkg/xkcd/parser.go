package xkcd

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

	fmt.Println(string(data))
	var result ResponseData
	if err := json.Unmarshal(data, &result); err != nil {
		return ResponseData{}, err
	}
	return result, nil
}

func responseParser(data []byte) (database.ParsedData, error) {

	result := database.ParsedData{}

	jsonData, err := decodeResponse(data)
	result.Keywords = words.StringNormalization(jsonData.Alt)
	result.Url = jsonData.Img

	return result, err
}
