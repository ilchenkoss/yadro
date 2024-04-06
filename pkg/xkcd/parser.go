package xkcd

import (
	"encoding/json"
	"myapp/pkg/database"
	"myapp/pkg/words"
)

type responseData struct {
	Alt string `json:"alt"`
	Img string `json:"img"`
}

func decodeResponse(data []byte) (responseData, bool) {

	var result responseData
	if err := json.Unmarshal(data, &result); err != nil {
		return responseData{}, true
	}
	return result, false
}

func responseParser(data []byte) (database.ParsedData, bool) {

	result := database.ParsedData{}
	err := false
	jsonData := responseData{}

	jsonData, err = decodeResponse(data)
	result.Keywords = words.StringNormalization(jsonData.Alt)
	result.Url = jsonData.Img

	return result, err
}
