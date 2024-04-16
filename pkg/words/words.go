package words

import (
	_ "embed"
	"github.com/kljensen/snowball"
	"regexp"
	"strings"
)

//go:embed wordsToRemove.txt
var WordsFile string

func CleanWord(uncleanedWord string) string {
	//clearing a word from non-word characters

	var regex = regexp.MustCompile(`[^a-zA-Z']+`)

	cleanedWord := regex.ReplaceAllString(uncleanedWord, "")

	return cleanedWord

}

func stemming(notNormalizedString []string) map[string]int {

	duplicateContainer := make(map[string]bool)
	stemmedWords := make(map[string]int)

	for _, word := range notNormalizedString {
		var stemmedWord, err = snowball.Stem(word, "english", true)

		//uniqueness
		if err == nil {
			if duplicateContainer[stemmedWord] {
				stemmedWords[stemmedWord]++
				continue
			}
			duplicateContainer[stemmedWord] = true
			stemmedWords[stemmedWord] = 1
		}
	}

	return stemmedWords
}

func loadStopWords() map[string]bool {

	stopWords := make(map[string]bool)

	words := strings.Fields(WordsFile)

	for _, word := range words {
		stopWords[strings.ToLower(word)] = true
	}

	return stopWords
}

func sifting(sliceWords []string, stopWords map[string]bool) []string {

	var keywords []string

	for _, word := range sliceWords {

		word = strings.ToLower(CleanWord(word))

		if !stopWords[word] && len(word) > 2 {
			keywords = append(keywords, word)
		}
	}

	return keywords

}

func StringNormalization(inputString string) map[string]int {

	//parse string
	stringFields := strings.Fields(inputString)
	//load stop words
	stopWords := loadStopWords()
	//sifting words from garbage
	siftingWords := sifting(stringFields, stopWords)
	//stemmed input words in string and uniqueness
	stemmedWords := stemming(siftingWords)

	return stemmedWords
}
