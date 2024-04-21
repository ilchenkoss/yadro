package indexing

import (
	_ "embed"
	"fmt"
	"myapp/pkg/database"
	"myapp/pkg/scraper"
	"myapp/pkg/words"
	"testing"
)

func BenchmarkFindComics(b *testing.B) {

	//weight requested words
	requestWords := words.StringNormalization("follow with me, my dear friend")
	weightRequest := createWeightRequest(requestWords)

	//data from db
	dbBytes := database.ReadBytesFromFile("database_test.json")
	dbData := scraper.DecodeFileData(dbBytes)

	//create index table
	weightData := createWeightData(dbData)
	indexDB := createIndexData(weightData)

	b.Run(
		fmt.Sprintf("WithoutIndexTable"), func(b *testing.B) {
			b.ResetTimer()
			wordsWeight := createWordsWeightWithoutIndex(weightRequest, dbData)
			createWeightComics(wordsWeight, weightRequest, dbData)
		})

	b.Run(
		fmt.Sprintf("WithIndexTable"), func(b *testing.B) {
			b.ResetTimer()
			createWeightComics(indexDB, weightRequest, dbData)
		})

}
