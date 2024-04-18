package scraper

import (
	"encoding/json"
	"myapp/pkg/words"
	"sync"
)

func decodeFileData(fileData []byte) map[int]ScrapedData {
	data := map[int]ScrapedData{}
	if err := json.Unmarshal(fileData, &data); err != nil {
		return data
	}
	return data
}

func codeFileData(bytesData map[int]ScrapedData) []byte {
	data, err := json.MarshalIndent(bytesData, "", "\t")
	if err != nil {
		return nil
	}
	return data
}

type ScrapedData struct {
	Keywords []string `json:"keywords"`
	Url      string   `json:"url"`
}
type ParsedData struct {
	ID       int      `json:"id"`
	Keywords []string `json:"keywords"`
	Url      string   `json:"url"`
}
type ResponseData struct {
	Alt string `json:"alt"`
	Img string `json:"img"`
	ID  int    `json:"num"`
}

func responseParser(data []byte) (ParsedData, error) {

	var response = ResponseData{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		return ParsedData{}, err
	}

	result := ParsedData{
		ID:       response.ID,
		Keywords: words.StringNormalization(response.Alt),
		Url:      response.Img,
	}

	return result, nil
}

func parserWorker(dbData map[int]ScrapedData, goodScrapesCh chan []byte, pwg *sync.WaitGroup, resultCh chan map[int]ScrapedData) {

	for scrape := range goodScrapesCh {
		data, err := responseParser(scrape)
		if err == nil {
			dbData[data.ID] = ScrapedData{
				Keywords: data.Keywords,
				Url:      data.Url,
			}
		}
		pwg.Done()
	}
	resultCh <- dbData
	return
}
