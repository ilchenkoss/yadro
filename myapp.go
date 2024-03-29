package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/kljensen/snowball"
	"os"
	"regexp"
	"strings"
)

var (
	inputString string
)

func cleanWord(uncleanedWord string) string {
	//clearing a word from non-word characters

	var regex = regexp.MustCompile(`[^a-zA-Z' ]+`)

	cleanedWord := regex.ReplaceAllString(uncleanedWord, " ")

	return cleanedWord

}

func stemming(notNormalizedString string) []string {

	var stemmedWords []string

	for _, word := range strings.Fields(notNormalizedString) {
		var stemmedWord, err = snowball.Stem(cleanWord(word), "english", true)
		if err == nil {
			stemmedWords = append(stemmedWords, stemmedWord)
		}
	}

	return stemmedWords
}

func loadStopWords() map[string]bool {

	stopWords := make(map[string]bool)

	file, err := os.Open("wordsToRemove.txt")
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

func sifting(stemmedWords []string, stopWords map[string]bool) string {

	duplicateContainer := make(map[string]bool)
	var output string

	for _, word := range stemmedWords {

		if !duplicateContainer[word] && !stopWords[word] && len(word) > 1 {
			duplicateContainer[word] = true
			output += strings.ToLower(word) + " "
		}
	}

	return strings.TrimSpace(output)

}

func stringNormalization(inputString string) string {

	//stemmed input words in string
	stemmedWords := stemming(inputString)
	//load stop words
	stopWords := loadStopWords()
	//sifting string from garbage
	outputString := sifting(stemmedWords, stopWords)

	return outputString
}

func main() {

	//приложение, которое нормализует перечисленные в виде аргументов слова (на английском).
	//Приложение должно отсеивать часто употребляемые слова
	//типа of/a/the/, местоимения и глагольные частицы (will)

	flag.StringVar(&inputString, "s", "string to normalize", "string to normalize")
	flag.Parse()

	//processing the separation words without a space
	replacer := strings.NewReplacer(",", " ", ".", " ", "?", " ", "!", " ")
	inputString = replacer.Replace(inputString)

	//result
	fmt.Println(stringNormalization(inputString))
}
