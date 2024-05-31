package words

import (
	"github.com/kljensen/snowball"
	"github.com/tjarratt/babble"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func comparisonSlices(expected []string, actual []string) (bool, map[string][]string) {

	var errDetails = make(map[string][]string)

	//sorting for stability result
	sort.Strings(expected)
	sort.Strings(actual)

	if reflect.DeepEqual(expected, actual) {
		return true, errDetails
	} else {

		//in expected, but not in actual
		for _, val := range expected {
			if !contains(actual, val) {
				errDetails["expected"] = append(errDetails["expected"], val)
			}
		}

		//in actual, but not in expected
		for _, val := range actual {
			if !contains(expected, val) {
				errDetails["actual"] = append(errDetails["actual"], val)
			}
		}

		return false, errDetails
	}
}

func TestExample1(t *testing.T) {

	notNormalizedString := "follower brings bunch of questions!"
	expected := map[string]KeywordsInfo{"follow": {1, 0}, "bunch": {1, 2}, "bring": {1, 1}, "question": {1, 3}}
	actual := StringNormalization(notNormalizedString)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nResult was incorrect. \n got: %v, \n want: %v.", actual, expected)
	}
}

func TestExample2(t *testing.T) {

	notNormalizedString := "i'll follow you as long as you are following me"
	expected := map[string]KeywordsInfo{"follow": {2, 0}, "long": {1, 1}}
	actual := StringNormalization(notNormalizedString)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nResult was incorrect. \n got: %v, \n want: %v.", actual, expected)
	}
}

func TestEmpty(t *testing.T) {

	notNormalizedString := ""
	expected := make(map[string]KeywordsInfo)
	actual := StringNormalization(notNormalizedString)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nResult was incorrect. \n got: %v, \n want: %v.", actual, expected)
	}
}

func TestDuplicate(t *testing.T) {

	notNormalizedString := "follow following, follower with followers"
	expected := map[string]KeywordsInfo{"follow": {4, 0}}
	actual := StringNormalization(notNormalizedString)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nResult was incorrect. \n got: %v, \n want: %v.", actual, expected)
	}
}

func TestSifting(t *testing.T) {

	tests := 10

	//words to generate
	keyWordsCount := 30
	trashWordsChance := 30

	trashWords := loadStopWords()

	//create slice with keys
	trashWordsKeys := make([]string, 0, len(trashWords))
	for key := range trashWords {
		trashWordsKeys = append(trashWordsKeys, key)
	}

	for i := 0; i < tests; i++ {

		generatedWords := generateWords(keyWordsCount, trashWords)
		resultSlice := append([]string(nil), generatedWords...)

		trashWordsCount := (keyWordsCount * trashWordsChance) / 100
		for i := 0; i < trashWordsCount; i++ {

			//pick random trashWord
			randomIndex := rand.Intn(len(trashWordsKeys))
			pick := trashWordsKeys[randomIndex]

			resultSlice = append(resultSlice, pick)
		}

		//result from sifting func
		actual := sifting(resultSlice, trashWords)

		if len(generatedWords) != len(actual) {

			_, errDetails := comparisonSlices(actual, generatedWords)

			t.Errorf("\nResult was incorrect. \n want: %s, \n got: %s.", errDetails["actual"], errDetails["expected"])
		}
	}
}

func generateWords(uniqueWords int, trashWords map[string]bool) []string {

	wordsCount := uniqueWords

	//duplicate stemmed words
	duplicateContainer := make(map[string]bool)

	generatedWords := make([]string, wordsCount)

	//word generator
	babbler := babble.NewBabbler()
	babbler.Count = 1

	for i := 0; i < uniqueWords; i++ {

		var successGen bool

		successGen, word, stemmedWord := generateUniqueWord(duplicateContainer, trashWords, 5, babbler)

		if !successGen {
			word = "word"
		}

		generatedWords[i] = word
		duplicateContainer[stemmedWord] = true
	}

	return generatedWords
}

func generateUniqueWord(duplicateContainer map[string]bool, trashWords map[string]bool, retry int, babbler babble.Babbler) (bool, string, string) {

	word := babbler.Babble()
	word = strings.ToLower(word)

	//uniqueness
	stemmedWord, err := snowball.Stem(word, "english", true)
	if err == nil {

		if trashWords[word] || len(word) <= 2 || duplicateContainer[stemmedWord] {
			return generateUniqueWord(duplicateContainer, trashWords, retry-1, babbler)
		}

		return true, word, stemmedWord

	} else {
		//if stem error
		return generateUniqueWord(duplicateContainer, trashWords, retry-1, babbler)
	}

}
