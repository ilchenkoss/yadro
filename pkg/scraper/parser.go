package scraper

import (
	"encoding/json"
	"fmt"
	"myapp/pkg/words"
	"sync"
)

func DecodeFileData(fileData []byte) map[int]ScrapedData {
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
	Keywords map[string]words.KeywordsInfo `json:"keywords"`
	Url      string                        `json:"url"`
}
type ParsedData struct {
	ID       int                           `json:"id"`
	Keywords map[string]words.KeywordsInfo `json:"keywords"`
	Url      string                        `json:"url"`
}
type ResponseData struct {
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
	Title      string `json:"title"`
	Img        string `json:"img"`
	ID         int    `json:"num"`
}

func responseParser(data []byte) (ParsedData, error) {

	var response = ResponseData{}
	err := json.Unmarshal(data, &response)
	if err != nil {
		return ParsedData{}, err
	}

	responseWords := fmt.Sprintf("%s %s %s", response.Title, response.Transcript, response.Alt)

	result := ParsedData{
		ID:       response.ID,
		Keywords: words.StringNormalization(responseWords),
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
