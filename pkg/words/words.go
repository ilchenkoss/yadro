package words

import (
	"bufio"
	"github.com/kljensen/snowball"
	"os"
	"regexp"
	"strings"
)

func CleanWord(uncleanedWord string) string {
	//clearing a word from non-word characters

	var regex = regexp.MustCompile(`[^a-zA-Z']+`)

	cleanedWord := regex.ReplaceAllString(uncleanedWord, "")

	return cleanedWord

}

func stemming(notNormalizedString []string) []string {

	duplicateContainer := make(map[string]bool)
	var stemmedWords []string

	for _, word := range notNormalizedString {
		var stemmedWord, err = snowball.Stem(word, "english", true)

		//uniqueness
		if err == nil && !duplicateContainer[stemmedWord] {
			duplicateContainer[stemmedWord] = true
			stemmedWords = append(stemmedWords, stemmedWord)
		}
	}

	return stemmedWords
}

func loadStopWords() map[string]bool {

	stopWords := make(map[string]bool)

	file, err := os.Open("./pkg/words/wordsToRemove.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := scanner.Text()
		word := strings.TrimSpace(line)

		if word != "" {
			stopWords[strings.ToLower(word)] = true
		}
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

func StringNormalization(inputString string) []string {

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
