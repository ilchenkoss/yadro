package indexing

import (
	"encoding/json"
	"fmt"
	"math"
	"myapp/pkg/database"
	"myapp/pkg/scraper"
	"myapp/pkg/words"
	"os"
	"sort"
)

const (
	RelevantFlag = true

	WeightStandard = 1.0

	//ВЕСА ДЛЯ ЗАПРОСА
	WeightRequestWordIndex     = 0.3 //чем слово ближе к началу, тем больше вес
	WeightRequestWordDuplicate = 1.0 //слово повторяется чаще в запросе

	//ВЕСА ДЛЯ ДАННЫХ
	WeightComicsWordIndex     = 0.2   //слово используется в начале
	WeightComicsWordDuplicate = 2.0   //слово повторяется чаще в комиксе
	WeightComicsActual        = 0.001 //актуальность комикса измеряется по ID

	//ВЕСА ДЛЯ ВЫДАЧИ
	WeightResponseRelevantCoverage = 10
	WeightResponseIDIndex          = 0.2
)

type IDWeight struct {
	ID     int
	Weight float64
}

func createWeightData(dbData map[int]scraper.ScrapedData) map[string][]IDWeight {
	weightData := make(map[string][]IDWeight)

	for ID, data := range dbData {
		if len(data.Keywords) == 0 {
			continue
		}
		for word, wordInfo := range data.Keywords {

			weight := WeightStandard
			weight += WeightComicsWordIndex / math.Log(float64(wordInfo.EntryIndex+2))
			if wordInfo.Repeat > 1 {
				weight += float64(wordInfo.Repeat) * WeightComicsWordDuplicate
			}
			weight += float64(ID) * WeightComicsActual

			weightData[word] = append(weightData[word], IDWeight{Weight: weight, ID: ID})
		}
	}
	return weightData
}

func createIndexData(weightData map[string][]IDWeight) map[string][]int {

	result := make(map[string][]int)

	for word, wordInfo := range weightData {

		//sort slice by weight
		sort.Slice(wordInfo, func(i, j int) bool {
			return wordInfo[i].Weight > wordInfo[j].Weight
		})

		//append sorted result
		var weightSlice []int
		for _, item := range wordInfo {
			weightSlice = append(weightSlice, item.ID)
		}

		result[word] = weightSlice
	}

	return result
}

type WordsWeight struct {
	Word   string
	Weight float64
}

func createWeightRequest(RequestWords map[string]words.KeywordsInfo) []WordsWeight {

	var resultWeightSlice []WordsWeight

	for word, wordInfo := range RequestWords {

		weight := WeightStandard
		weight += WeightRequestWordIndex / math.Log(float64(wordInfo.EntryIndex+2))
		weight += float64(wordInfo.Repeat) * WeightRequestWordDuplicate

		resultWeightSlice = append(resultWeightSlice, WordsWeight{Word: word, Weight: weight})
	}

	//sort slice by weight
	sort.Slice(resultWeightSlice, func(i, j int) bool {
		return resultWeightSlice[i].Weight > resultWeightSlice[j].Weight
	})
	fmt.Println("response weight: ", resultWeightSlice)

	return resultWeightSlice
}

func createIndexingDB(dbData map[int]scraper.ScrapedData) map[string][]int {
	weightData := createWeightData(dbData)
	indexData := createIndexData(weightData)
	return indexData
}

func createWeightComics(indexData map[string][]int, indexRequest []WordsWeight, dbData map[int]scraper.ScrapedData) []IDWeight {

	//проверяем, какие комиксы содержат как можно больше запрашиваемых слов и умножаем их на вес релевантного слова
	//вычетаем вес нерелевантных слов

	//определяем результат
	var result []IDWeight

	//определяем веса ID, которые используются в словах
	competitiveIDs := make(map[int]float64)

	//определяем все Words, которые участвуют в конкурсе :)
	competitiveWords := make(map[string]bool, len(indexRequest))

	//по каждому слову в запросе
	for _, indexedRequestWord := range indexRequest {
		if RelevantFlag {
			competitiveWords[indexedRequestWord.Word] = true
		}
		idsWithWord := indexData[indexedRequestWord.Word]
		//определяем вес ID
		for index, ID := range idsWithWord {
			//вес слова + коррекция на индекс в отсортированном по релевантности слайсе
			competitiveIDs[ID] += indexedRequestWord.Weight + WeightResponseIDIndex/math.Log(float64(index+2))
		}
	}

	for ID, weight := range competitiveIDs {

		if RelevantFlag {
			//корректировка по наполненности выдачи искомыми словами
			relevantWord := 0
			for word := range dbData[ID].Keywords {
				_, ok := competitiveWords[word]
				if ok {
					relevantWord++
				}
			}
			weight += float64((relevantWord / len(dbData[ID].Keywords)) * WeightResponseRelevantCoverage)
		}

		result = append(result, IDWeight{ID: ID, Weight: weight})
	}

	//sort slice by weight
	sort.Slice(result, func(i, j int) bool {
		return result[i].Weight > result[j].Weight
	})

	return result
}

func MainIndexing(requestString string, dbPath string, indexPath string) {

	//load db
	dbDataBytes := database.ReadBytesFromFile(dbPath)
	dbData := scraper.DecodeFileData(dbDataBytes)

	//create indexedDB
	indexDB := createIndexingDB(dbData)

	//save indexedDB
	data, _ := json.MarshalIndent(indexDB, "", "\t")
	os.WriteFile(indexPath, data, 0644)

	//create indexedRequest
	requestWords := words.StringNormalization(requestString)
	indexRequest := createWeightRequest(requestWords)

	//get response
	response := createWeightComics(indexDB, indexRequest, dbData)

	fmt.Printf("\n\n\n")

	limitedResponse := response
	if len(response) > 10 {
		limitedResponse = response[:10]
	}
	for _, responseData := range limitedResponse {
		fmt.Printf("https://xkcd.com/%d/\n", responseData.ID)

	}

}
