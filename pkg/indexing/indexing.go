package indexing

import (
	"encoding/json"
	"fmt"
	"myapp/pkg/database"
	"myapp/pkg/scraper"
	"myapp/pkg/words"
	"sort"
)

const (
	WeightStandard = 10.0

	//ВЕСА ДЛЯ ЗАПРОСА
	WeightRequestWordIndex     = -0.3 //первое слово имеет больший вес к последующему
	WeightRequestWordDuplicate = 2    //слово повторяется чаще в запросе

	//ВЕСА ДЛЯ ДАННЫХ
	WeightComicsWordIndex     = -0.1   //слово используется в начале
	WeightComicsWordDuplicate = 2.0    //слово повторяется чаще в комиксе
	WeightComicsActual        = -0.001 //актуальность комикса измеряется по ID

	//ВЕСА ДЛЯ ВЫДАЧИ
	WeightResponseRelevantWord     = 5
	WeightResponseRelevantCoverage = 10
	WeightResponseIDIndexFromDB    = -0.2
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
			weight += float64(wordInfo.EntryIndex+1) * WeightComicsWordIndex
			weight += float64(wordInfo.Repeat) * WeightComicsWordDuplicate
			weight += float64(ID) * WeightComicsActual

			weightData[word] = append(weightData[word], IDWeight{Weight: weight, ID: ID})
		}
	}
	fmt.Println(weightData)
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
		weight += float64(wordInfo.EntryIndex+1) * WeightRequestWordIndex
		weight += float64(wordInfo.Repeat) * WeightRequestWordDuplicate

		resultWeightSlice = append(resultWeightSlice, WordsWeight{Word: word, Weight: weight})
	}

	//sort slice by weight
	sort.Slice(resultWeightSlice, func(i, j int) bool {
		return resultWeightSlice[i].Weight > resultWeightSlice[j].Weight
	})

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

	//определяем все ID, которые включают в себя запрошенные слова
	competitiveIDs := make(map[int]float64) //map[id]word
	//определяем все Words, которые участвуют в конкурсе :)
	competitiveWords := make(map[string]bool)

	for _, indexedRequestWord := range indexRequest {
		competitiveWords[indexedRequestWord.Word] = true

		idsWithWord := indexData[indexedRequestWord.Word]
		//определяем первичный вес для каждого ID на основе сортированного ранее слайса из бд и веса слова запроса
		for index, ID := range idsWithWord {
			if competitiveIDs[ID] == 0 {
				competitiveIDs[ID] = WeightStandard
			}
			competitiveIDs[ID] += indexedRequestWord.Weight + float64(index)*WeightResponseIDIndexFromDB
		}
	}

	//производим коррекцию на нерелвантные слова слов
	for ID, weight := range competitiveIDs {

		for word := range dbData[ID].Keywords {
			_, ok := competitiveWords[word]
			if !ok {
				weight += float64(len(dbData[ID].Keywords))
				continue
			}
			weight += WeightResponseRelevantWord
		}

		result = append(result, IDWeight{ID: ID, Weight: weight})
	}

	//sort slice by weight
	sort.Slice(result, func(i, j int) bool {
		return result[i].Weight > result[j].Weight
	})

	return result
}

func MainIndexing(requestString string) {

	//чем меньше нерелевантных слов, тем лучше??

	//load db
	dbDataBytes := database.ReadBytesFromFile("pkg/database/database.json")
	dbData := scraper.DecodeFileData(dbDataBytes)

	//create indexedDB
	indexDB := createIndexingDB(dbData)

	//save indexedDB
	data, _ := json.MarshalIndent(indexDB, "", "\t")
	database.WriteToFile("pkg/database/index.json", data, false)

	//create indexedRequest
	requestWords := words.StringNormalization(requestString)
	indexRequest := createWeightRequest(requestWords)

	//get response
	response := createWeightComics(indexDB, indexRequest, dbData)

	fmt.Printf("\n\n\n")

	if len(response) >= 10 {
		fmt.Println(response[:10])
	} else {
		fmt.Println(response)
	}

}
