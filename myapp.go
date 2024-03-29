package main

import (
	"flag"
	"fmt"
	"github.com/jdkato/prose/v2"
	"github.com/kljensen/snowball"
	"regexp"
	"strings"
)

var (
	inputString string
)

func cleanWord(uncleanedWord string) string {

	var regex = regexp.MustCompile(`[^a-zA-Z' ]+`)
	cleanedWord := regex.ReplaceAllString(uncleanedWord, " ")
	return cleanedWord

}

func stemString(notNormalizedString string) string {

	var stemmedString string

	for _, word := range strings.Fields(notNormalizedString) {
		var stemmedWord, err = snowball.Stem(cleanWord(word), "english", true)
		if err == nil {
			stemmedString += " " + stemmedWord
		}
	}

	return stemmedString
}

func siftingString(stemString string) string {

	var siftedString string

	duplicateWords := map[string]bool{}

	doc, _ := prose.NewDocument(stemString)
	for _, tok := range doc.Tokens() {

		//siftingTags and siftingWords in config.go
		if siftingTags[tok.Tag] && !siftingWords[tok.Text] && !duplicateWords[tok.Text] && len(tok.Text) > 1 {
			duplicateWords[tok.Text] = true
			siftedString += tok.Text + " "
		}
	}

	if len(siftedString) > 1 && siftedString[len(siftedString)-1:] == " " {
		return siftedString[:len(siftedString)-1]
	} else {
		return siftedString
	}
}

func stringNormalization(notNormalizedString string) string {

	normalizedString := stemString(notNormalizedString)
	siftedString := siftingString(normalizedString)

	return siftedString
}

func main() {
	flag.StringVar(&inputString, "s", "string to normalize", "string to normalize")
	flag.Parse()
	fmt.Println(stringNormalization(inputString))
}
