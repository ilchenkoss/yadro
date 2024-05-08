package words

import (
	_ "embed"
	"github.com/kljensen/snowball"
	"regexp"
	"strings"
	"unicode"
)

//go:embed wordsToRemove.txt
var WordsFile string

func CleanWord(uncleanedWord string) string {
	//clearing a word from non-word characters

	var regex = regexp.MustCompile(`[^a-zA-Z']+`)

	cleanedWord := regex.ReplaceAllString(uncleanedWord, "")

	return cleanedWord

}

type KeywordsInfo struct {
	Repeat     int
	EntryIndex int
}

func stemming(notNormalizedString []string) map[string]KeywordsInfo {

	duplicateContainer := make(map[string]bool)
	stemmedWords := make(map[string]KeywordsInfo)

	for wordIndex, word := range notNormalizedString {

		var stemmedWord, err = snowball.Stem(word, "english", true)

		if stemmedWord == "" {
			continue
		}

		//uniqueness
		if err == nil {
			if duplicateContainer[stemmedWord] {
				repeat := stemmedWords[stemmedWord].Repeat
				entryIndex := stemmedWords[stemmedWord].EntryIndex
				stemmedWords[stemmedWord] = KeywordsInfo{repeat + 1, entryIndex}
				continue
			}
			duplicateContainer[stemmedWord] = true
			stemmedWords[stemmedWord] = KeywordsInfo{1, wordIndex}
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

func StringNormalization(inputString string) map[string]KeywordsInfo {

	//parse string
	stringFields := strings.FieldsFunc(inputString, func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSpace(r)
	})
	//load stop words
	stopWords := loadStopWords()
	//sifting words from garbage
	siftingWords := sifting(stringFields, stopWords)
	//stemmed input words in string and uniqueness
	stemmedWords := stemming(siftingWords)

	return stemmedWords
}
