package scraper

import (
	"encoding/json"
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

type KeywordInfo struct {
	Repeat     int
	EntryIndex int
	Position   string
}

type ScrapedData struct {
	Keywords map[string][]KeywordInfo `json:"keywords"`
	Url      string                   `json:"url"`
}

type ParsedData struct {
	ID                 int                           `json:"id"`
	KeywordsTitle      map[string]words.KeywordsInfo `json:"keywords_title"`
	KeywordsTranscript map[string]words.KeywordsInfo `json:"keywords_transcript"`
	KeywordsAlt        map[string]words.KeywordsInfo `json:"keywords_alt"`
	Url                string                        `json:"url"`
}

type ResponseData struct {
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
	Title      string `json:"title"`
	Img        string `json:"img"`
	ID         int    `json:"num"`
}

func responseParser(data []byte) (ParsedData, error) {

	var response ResponseData
	err := json.Unmarshal(data, &response)
	if err != nil {
		return ParsedData{}, err
	}

	result := ParsedData{
		ID:                 response.ID,
		KeywordsTitle:      words.StringNormalization(response.Title),
		KeywordsTranscript: words.StringNormalization(response.Transcript),
		KeywordsAlt:        words.StringNormalization(response.Alt),
		Url:                response.Img,
	}

	return result, nil
}

func mergeWords(keywordsTitle map[string]words.KeywordsInfo, keywordsTranscript map[string]words.KeywordsInfo, keywordsAlt map[string]words.KeywordsInfo) map[string][]KeywordInfo {
	mergedMap := make(map[string][]KeywordInfo)

	for word, wordInfo := range keywordsTitle {
		mergedMap[word] = append(mergedMap[word], KeywordInfo{Repeat: wordInfo.Repeat, EntryIndex: wordInfo.EntryIndex, Position: "title"})
	}

	for word, wordInfo := range keywordsTranscript {
		mergedMap[word] = append(mergedMap[word], KeywordInfo{Repeat: wordInfo.Repeat, EntryIndex: wordInfo.EntryIndex, Position: "transcript"})
	}

	for word, wordInfo := range keywordsAlt {
		mergedMap[word] = append(mergedMap[word], KeywordInfo{Repeat: wordInfo.Repeat, EntryIndex: wordInfo.EntryIndex, Position: "alt"})
	}

	return mergedMap
}

func parserWorker(dbData map[int]ScrapedData, goodScrapesCh chan []byte, pwg *sync.WaitGroup, resultCh chan map[int]ScrapedData, scrapeScore *int) {

	for scrape := range goodScrapesCh {

		data, err := responseParser(scrape)

		if err == nil {
			*scrapeScore++
			dbData[data.ID] = ScrapedData{
				Keywords: mergeWords(data.KeywordsTitle, data.KeywordsTranscript, data.KeywordsAlt),
				Url:      data.Url,
			}
		}
		pwg.Done()
	}
	resultCh <- dbData
	return
}
