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
	WeightRequestWordIndex     = 0.3 //чем ближе слово к началу, тем больше его вес
	WeightRequestWordDuplicate = 1.0 //чем чаще слово повторяется, тем больше его вес

	//ВЕСА ДЛЯ ДАННЫХ
	WeightComicsWordIndex     = 0.2   //чем ближе слово к началу, тем больше его вес
	WeightComicsWordDuplicate = 2.0   //чем чаще слово повторяется, тем больше его вес
	WeightComicsActual        = 0.001 //чем больше ID, тем актуальнее комикс

	//ВЕСА ДЛЯ ВЫДАЧИ
	WeightResponseRelevantCoverage = 10  //вес 100% покрытия комикса ключевыми словами
	WeightResponseIDIndex          = 0.2 //чем ближе слово к началу, тем больше его вес
)

type IDWeight struct {
	ID     int
	Weight float64
}

func createWeightData(dbData map[int]scraper.ScrapedData) map[string][]IDWeight {
	//Взвешиваем комиксы и создаем слайс ID, weight
	weightData := make(map[string][]IDWeight)

	for ID, data := range dbData {
		if len(data.Keywords) == 0 {
			continue
		}
		for word, wordInfo := range data.Keywords {

			var weight float64
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
	//Сортируем слайс по weight и преобразовываем в конечный вид
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
	//определяем вес запрашиваемых слов
	var resultWeightSlice []WordsWeight

	for word, wordInfo := range RequestWords {

		var weight float64
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

	//взвешиваем и сортируем комиксы

	var result []IDWeight
	competitiveIDs := make(map[int]float64)
	competitiveWords := make(map[string]bool, len(indexRequest))

	//создание веса ID на основе позиции в index.json
	for _, indexedRequestWord := range indexRequest {
		if RelevantFlag {
			competitiveWords[indexedRequestWord.Word] = true
		}

		idsWithWord := indexData[indexedRequestWord.Word]
		for index, ID := range idsWithWord {
			competitiveIDs[ID] += indexedRequestWord.Weight + WeightResponseIDIndex/math.Log(float64(index+2))
		}
	}

	//Преобразование результатов в slice для последующей сортировки и коррекция веса по заполнению ключевыми словами у комикса
	for ID, weight := range competitiveIDs {

		if RelevantFlag {
			//коррекция веса по заполнению ключевыми словами у комикса
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
