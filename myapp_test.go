package main

import (
	"bytes"
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
	expected := []string{"follow", "bunch", "bring", "question"}
	actual := stringNormalization(notNormalizedString)

	if equal, errDetails := comparisonSlices(expected, actual); equal == false {
		t.Errorf("\nResult was incorrect. \n got: %s, \n want: %s.", errDetails["actual"], errDetails["expected"])
	}
}

func TestExample2(t *testing.T) {

	notNormalizedString := "i'll follow you as long as you are following me"
	expected := []string{"follow", "long"}
	actual := stringNormalization(notNormalizedString)

	if equal, errDetails := comparisonSlices(expected, actual); equal == false {
		t.Errorf("\nResult was incorrect. \n got: %s, \n want: %s.", errDetails["actual"], errDetails["expected"])
	}
}

func TestEmpty(t *testing.T) {

	notNormalizedString := ""
	var expected []string
	actual := stringNormalization(notNormalizedString)

	if equal, errDetails := comparisonSlices(expected, actual); equal == false {
		t.Errorf("\nResult was incorrect. \n want: %s, \n got: %s.", errDetails["actual"], errDetails["expected"])
	}
}

func TestDuplicate(t *testing.T) {

	notNormalizedString := "follow following, follower with followers"
	expected := []string{"follow"}
	actual := stringNormalization(notNormalizedString)

	if equal, errDetails := comparisonSlices(expected, actual); equal == false {
		t.Errorf("\nResult was incorrect. \n got: %s, \n want: %s.", errDetails["actual"], errDetails["expected"])
	}
}

func TestSifting(t *testing.T) {

	tests := 100

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

		generatedWords := generateUniqueWords(keyWordsCount, trashWords)
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

func TestSynth(t *testing.T) {

	var tests = 10
	wordsCount := 15

	punctuationChance := 40
	punctuation := []string{
		".",
		",",
		" -",
		"?",
		"!",
	}

	trashWordsChance := 20
	trashWords := loadStopWords()

	//create slice with keys
	trashWordKeys := make([]string, 0, len(trashWords))
	for key := range trashWords {
		trashWordKeys = append(trashWordKeys, key)
	}

	for i := 0; i < tests; i++ {

		//buffer for synth string
		var synthStringBuffer bytes.Buffer

		generatedWords := generateUniqueWords(wordsCount, trashWords)

		for index, word := range generatedWords {

			synthStringBuffer.WriteString(strings.ToLower(word))

			//add punctuation
			if rand.Intn(100) < punctuationChance {

				//pick random punctuation
				randomIndex := rand.Intn(len(punctuation))
				pick := punctuation[randomIndex]

				synthStringBuffer.WriteString(pick)
			}

			//add trashWord

			if rand.Intn(100) < trashWordsChance {
				//pick random trashWord
				randomIndex := rand.Intn(len(trashWordKeys))
				pick := trashWordKeys[randomIndex]

				synthStringBuffer.WriteString(" " + pick)
			}

			if index != wordsCount-1 {
				synthStringBuffer.WriteString(" ")
			}
		}

		finalString := synthStringBuffer.String()
		actual := stringNormalization(finalString)

		if len(generatedWords) != len(actual) {

			_, errDetails := comparisonSlices(actual, generatedWords)

			t.Errorf("\nResult was incorrect. \n got: %s, \n want: %s.", errDetails["actual"], errDetails["expected"])
		}
	}
}

func generateUniqueWords(uniqueWordsCount int, trashWords map[string]bool) []string {

	duplicateContainer := make(map[string]bool)
	generatedWords := make([]string, uniqueWordsCount)

	//word generator
	babbler := babble.NewBabbler()
	babbler.Count = 1

	for i := 0; i < uniqueWordsCount; i++ {

		var successGen bool

		successGen, word := generateUniqueWord(duplicateContainer, trashWords, 5, babbler)

		if !successGen {
			word = "word"
		}

		generatedWords[i] = word
		duplicateContainer[word] = true
	}

	return generatedWords
}

func generateUniqueWord(duplicateContainer map[string]bool, trashWords map[string]bool, retry int, babbler babble.Babbler) (bool, string) {

	word := babbler.Babble()
	word = strings.ToLower(word)

	if trashWords[word] || len(word) <= 2 || duplicateContainer[word] {
		return generateUniqueWord(duplicateContainer, trashWords, retry-1, babbler)
	}

	return true, word
}
