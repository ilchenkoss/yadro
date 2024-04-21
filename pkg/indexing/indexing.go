package indexing

import (
	"encoding/json"
	"fmt"
	"math"
	"myapp/pkg/scraper"
	"myapp/pkg/words"
	"os"
	"sort"
)

const (
	//ВЕСА ДЛЯ ЗАПРОСА
	WeightRequestWordIndex     = 0.3 //чем ближе слово к началу, тем больше его вес
	WeightRequestWordDuplicate = 1.0 //чем чаще слово повторяется, тем больше его вес

	//ВЕСА ДЛЯ ДАННЫХ
	WeightComicsWordIndex              = 0.2   //чем ближе слово к началу, тем больше его вес
	WeightComicsWordDuplicate          = 2.0   //чем чаще слово повторяется, тем больше его вес
	WeightComicsActual                 = 0.001 //чем больше ID, тем актуальнее комикс
	WeightComicsWordPositionTitle      = 6.0
	WeightComicsWordPositionTranscript = 0.9
	WeightComicsWordPositionAlt        = 0.8

	//ВЕСА ДЛЯ ВЫДАЧИ
	WeightResponseIDIndex                    = 0.2  //вес позиции комикса в index.json
	CoverageFlag                             = true // будем ли проверять покрытие комикса запрашиваемыми
	WeightResponseRelevantCoverage           = 1.5  //вес 100% покрытия комикса запрашиваемыми словами
	WeightResponseRelevantCoverageTitle      = 1.0
	WeightResponseRelevantCoverageTranscript = 0.8
	WeightResponseRelevantCoverageAlt        = 0.7
)

type IDWeight struct {
	ID     int
	Weight float64
}

func weightByKeyword(info scraper.KeywordInfo, positionWeight float64) float64 {

	var weight float64
	weight += WeightComicsWordIndex / math.Log(float64(info.EntryIndex+2))
	if info.Repeat > 1 {
		weight += float64(info.Repeat) * WeightComicsWordDuplicate
	}
	return weight * positionWeight
}

func createWeightData(dbData map[int]scraper.ScrapedData) map[string][]IDWeight {
	//Взвешиваем комиксы и создаем слайс ID, weight
	weightData := make(map[string][]IDWeight)
	for ID, data := range dbData {

		idWeight := IDWeight{ID: ID, Weight: 0}

		for word, sliceWordInfo := range data.Keywords {
			for _, wordInfo := range sliceWordInfo {
				if wordInfo.Position == "title" {
					idWeight.Weight += weightByKeyword(wordInfo, WeightComicsWordPositionTitle)
				}
				if wordInfo.Position == "transcript" {
					idWeight.Weight += weightByKeyword(wordInfo, WeightComicsWordPositionTranscript)
				}
				if wordInfo.Position == "alt" {
					idWeight.Weight += weightByKeyword(wordInfo, WeightComicsWordPositionAlt)
				}
			}
			idWeight.Weight += float64(ID) * WeightComicsActual
			weightData[word] = append(weightData[word], idWeight)
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

func CreateIndexingDB(dbData map[int]scraper.ScrapedData, indexPath string) map[string][]int {

	weightData := createWeightData(dbData)
	indexData := createIndexData(weightData)

	//save indexedDB
	data, errCode := json.MarshalIndent(indexData, "", "\t")
	if errCode != nil {
		fmt.Printf("Can't save index table: %s", errCode)
		return indexData
	}
	errWrite := os.WriteFile(indexPath, data, 0644)
	if errWrite != nil {
		fmt.Printf("Can't save index table: %s", errWrite)
		return indexData
	}
	return indexData
}

func createWordsWeightWithoutIndex(weightRequest []WordsWeight, dbData map[int]scraper.ScrapedData) map[string][]int {

	weightData := make(map[string][]IDWeight)

	for ID, data := range dbData {
		if len(data.Keywords) == 0 {
			continue
		}
		// Проверяем, содержатся ли ключевые слова из weightRequest в Keywords
		for _, wordRequestInfo := range weightRequest {
			if keyWordInfo, ok := data.Keywords[wordRequestInfo.Word]; ok {
				//ID содержит искоме слово
				//Создаем Вес для ID
				idWeight := IDWeight{ID: ID, Weight: 0}

				for _, wordInfo := range keyWordInfo {
					if wordInfo.Position == "title" {
						idWeight.Weight += weightByKeyword(wordInfo, WeightComicsWordPositionTitle)
					}
					if wordInfo.Position == "transcript" {
						idWeight.Weight += weightByKeyword(wordInfo, WeightComicsWordPositionTranscript)
					}
					if wordInfo.Position == "alt" {
						idWeight.Weight += weightByKeyword(wordInfo, WeightComicsWordPositionAlt)
					}
				}
				idWeight.Weight += float64(ID) * WeightComicsActual
				weightData[wordRequestInfo.Word] = append(weightData[wordRequestInfo.Word], idWeight)
			}

		}
	}

	wordsWeight := createIndexData(weightData)
	return wordsWeight
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

	return resultWeightSlice
}

func createWeightComics(indexData map[string][]int, indexRequest []WordsWeight, dbData map[int]scraper.ScrapedData) []IDWeight {

	//взвешиваем и сортируем комиксы

	var result []IDWeight
	competitiveIDs := make(map[int]float64)
	competitiveWords := make(map[string]bool, len(indexRequest))

	//создание веса ID на основе позиции в index.json
	for _, indexedRequestWord := range indexRequest {
		if CoverageFlag {
			competitiveWords[indexedRequestWord.Word] = true
		}

		idsWithWord := indexData[indexedRequestWord.Word]
		for index, ID := range idsWithWord {
			competitiveIDs[ID] += indexedRequestWord.Weight + WeightResponseIDIndex/math.Log(float64(index+2))
		}
	}

	//Преобразование результатов в slice для последующей сортировки и коррекция веса по заполнению ключевыми словами у комикса
	for ID, weight := range competitiveIDs {

		if CoverageFlag {

			//коррекция веса по заполнению ключевыми словами у комикса

			relevantWordTitle := 0
			countWordTitle := 0
			relevantWordTranscript := 0
			countWordTranscript := 0
			relevantWordAlt := 0
			countWordAlt := 0

			for word, sliceWordInfo := range dbData[ID].Keywords {

				_, ok := competitiveWords[word]

				for _, wordInfo := range sliceWordInfo {
					if wordInfo.Position == "title" {
						if ok {
							relevantWordTitle += wordInfo.Repeat
						}
						countWordTitle += wordInfo.Repeat
					}
					if wordInfo.Position == "transcript" {
						if ok {
							relevantWordTranscript += wordInfo.Repeat
						}
						countWordTranscript += wordInfo.Repeat
					}
					if wordInfo.Position == "alt" {
						if ok {
							relevantWordAlt += wordInfo.Repeat
						}
						countWordAlt += wordInfo.Repeat
					}
				}
			}

			if countWordTitle > 0 {
				coverage := float64(relevantWordTitle / countWordTitle)
				weight += coverage * WeightResponseRelevantCoverage * WeightResponseRelevantCoverageTitle
			}
			if countWordTranscript > 0 {
				coverage := float64(relevantWordTranscript / countWordTranscript)
				weight += coverage * WeightResponseRelevantCoverage * WeightResponseRelevantCoverageTranscript
			}
			if countWordAlt > 0 {
				coverage := float64(relevantWordAlt / countWordAlt)
				weight += coverage * WeightResponseRelevantCoverage * WeightResponseRelevantCoverageAlt
			}

		}
		result = append(result, IDWeight{ID: ID, Weight: weight})
	}

	//sort slice by weight
	sort.Slice(result, func(i, j int) bool {
		return result[i].Weight > result[j].Weight
	})

	return result
}

func FindComics(requestString string, indexDB map[string][]int, dbData map[int]scraper.ScrapedData) {

	//create indexedRequest
	requestWords := words.StringNormalization(requestString)
	weightRequest := createWeightRequest(requestWords)

	//get response
	response := createWeightComics(indexDB, weightRequest, dbData)

	fmt.Printf("\nComics for you:\n\n")

	limitedResponse := response
	if len(response) > 10 {
		limitedResponse = response[:10]
	}
	for _, responseData := range limitedResponse {
		fmt.Printf("https://xkcd.com/%d/\n", responseData.ID)

	}

}
