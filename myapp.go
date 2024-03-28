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

func _cleanWord(uncleanedWord string) string {

	var regex = regexp.MustCompile(`[^a-zA-Z' ]+`)
	cleanedWord := regex.ReplaceAllString(uncleanedWord, " ")
	return cleanedWord

}

func _stemString(notNormalizedString string) string {

	var stemmedString string

	for _, word := range strings.Fields(notNormalizedString) {
		var stemmedWord, err = snowball.Stem(_cleanWord(word), "english", true)
		if err == nil {
			stemmedString += " " + stemmedWord
		}
	}

	return stemmedString
}

func _siftingString(stemString string, siftingTags map[string]bool, siftingWords map[string]bool) string {

	var siftedString string

	doc, _ := prose.NewDocument(stemString)
	for _, tok := range doc.Tokens() {

		if siftingTags[tok.Tag] && !siftingWords[tok.Text] && len(tok.Text) > 1 {
			siftingWords[tok.Text] = true
			siftedString += tok.Text + " "
		}
		//fmt.Println(tok.Text, tok.Tag)
	}

	if len(siftedString) < 1 {
		return ""
	} else {
		return siftedString[:len(siftedString)-1]
	}
}

func stringNormalization(notNormalizedString string) string {

	siftingTags := map[string]bool{
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

	siftingWords := map[string]bool{
		"be": true,
	}

	normalizedString := _stemString(notNormalizedString)
	siftedString := _siftingString(normalizedString, siftingTags, siftingWords)

	return siftedString
}

func init() {
	flag.StringVar(&inputString, "s", "good a with you need to will be have", "string to words")
}

func main() {
	flag.Parse()
	fmt.Println(stringNormalization(inputString))
}
