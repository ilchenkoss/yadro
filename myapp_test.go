package main

import (
	"github.com/tjarratt/babble"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestExample1(t *testing.T) {
	notNormalizedString := "follower brings bunch of questions"
	expected := "follow bring bunch question"
	actual := stringNormalization(notNormalizedString)

	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestExample2(t *testing.T) {
	notNormalizedString := "i'll follow you as long as you are following me"
	expected := "follow long"
	actual := stringNormalization(notNormalizedString)
	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestEmpty(t *testing.T) {
	notNormalizedString := ""
	expected := ""
	actual := stringNormalization(notNormalizedString)
	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	}
}

func TestAutoGen(t *testing.T) {

	babbler := babble.NewBabbler()
	babbler.Count = 1

	punctuationChance := 40
	punctuation := []string{
		".",
		",",
		" -",
		"?",
		"!",
		"'s",
	}

	trashWordsChance := 20
	trashWords := []string{
		"to",
		"be",
		"will",
		"she",
		"he",
		"we",
		"a",
		"with",
		"where",
	}

	wordsCount := 15
	var words string

	i := 0
	for i < wordsCount {

		words += babbler.Babble()
		//add test words
		i++

		//punctuation intervention
		if rand.Intn(100) < punctuationChance {
			randomIndex := rand.Intn(len(punctuation))
			pick := punctuation[randomIndex]
			words += pick
		}

		//trashWord intervention
		if rand.Intn(100) < trashWordsChance {
			randomIndex := rand.Intn(len(trashWords))
			pick := trashWords[randomIndex]
			words += " " + pick
		}

		if i != wordsCount-1 {
			words += " "
		}

	}

	expected := len(strings.Split(stringNormalization(words), " "))

	if expected != wordsCount {
		t.Errorf("Result was incorrect, got: %s, want: %s.", strconv.Itoa(expected), strconv.Itoa(wordsCount))
	}

}
