package main

import (
	"flag"
	"fmt"
	"github.com/jdkato/prose/v2"
	"github.com/kljensen/snowball"
	"regexp"
	"strings"
)

func cleanWord(uncleanedWord string) string {

	var regex = regexp.MustCompile(`[^a-zA-Z' ]+`)
	cleanedWord := regex.ReplaceAllString(uncleanedWord, " ")
	return cleanedWord

}

func stringNormalization(notNormalizedString string) string {

	var resultedWords = make(map[string]bool)
	var resultedString string
	var cleanString string

	tags := map[string]bool{
		"(":    false, //left round bracket
		")":    false, //right round bracket
		",":    false, //comma
		":":    false, //colon
		".":    false, //period
		"''":   false, //closing quotation mark
		"``":   false, //opening quotation mark
		"#":    false, //number sign
		"$":    false, //currency
		"CC":   false, //conjunction, coordinating
		"CD":   false, //cardinal number
		"DT":   false, //determiner
		"EX":   false, //existential there
		"FW":   false, //foreign word
		"IN":   false, //conjunction, subordinating or preposition
		"JJ":   true,  //adjective
		"JJR":  false, //adjective, comparative
		"JJS":  false, //adjective, superlative
		"LS":   false, //list item marker
		"MD":   false, //verb, modal auxiliary
		"NN":   true,  //noun, singular or mass
		"NNP":  false, //noun, proper singular
		"NNPS": false, //noun, proper plural
		"NNS":  false, //noun, plural
		"PDT":  false, //predeterminer
		"POS":  false, //possessive ending
		"PRP":  false, //pronoun, personal
		"PRP$": false, //pronoun, possessive
		"RB":   true,  //adverb
		"RBR":  false, //adverb, comparative
		"RBS":  false, //adverb, superlative
		"RP":   false, //adverb, particle
		"SYM":  false, //symbol
		"TO":   false, //infinitival to
		"UH":   false, //interjection
		"VB":   true,  //verb, base form
		"VBD":  false, //verb, past tense
		"VBG":  false, //verb, gerund or present participle
		"VBN":  false, //verb, past participle
		"VBP":  false, //verb, non-3rd person singular present
		"VBZ":  false, //verb, 3rd person singular present
		"WDT":  false, //wh-determiner
		"WP":   false, //wh-pronoun, personal
		"WP$":  false, //wh-pronoun, possessive
		"WRB":  false, //wh-adverb
	}

	for _, word := range strings.Fields(notNormalizedString) {
		var stemmedWord, err = snowball.Stem(cleanWord(word), "english", true)
		if err == nil {
			cleanString = cleanString + " " + stemmedWord
		}
	}

	doc, _ := prose.NewDocument(cleanString)
	for _, tok := range doc.Tokens() {

		if tags[tok.Tag] && !resultedWords[tok.Text] && len(tok.Text) > 1 {
			resultedWords[tok.Text] = true
			resultedString += tok.Text + " "
		}
		if tok.Text != " " {
			println(tok.Tag, tok.Text)
		}

	}

	return resultedString[:len(resultedString)-1]

}

func main() {

	uncleanedWords := flag.String("s", "i you s f long sugar i'll follower follow brings bunch of questions", "string to words")
	flag.Parse()

	fmt.Println(stringNormalization(*uncleanedWords))
}
