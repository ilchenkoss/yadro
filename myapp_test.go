package main

import (
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

//
//func TestGenerateWords(t *testing.T) {
//
//	words := 10
//	//punctuation := map[string]bool{
//	//	".": true,
//	//	",": true,
//	//	"-": true,
//	//	"?": true,
//	//	"1": true,
//	//}
//
//	babbler := babble.NewBabbler()
//	babbler.Separator = " "
//	babbler.Count = words
//	fakeString := babbler.Babble()
//	fmt.Println(fakeString)
//	fmt.Println(stringNormalization(fakeString))
//
//	return
//}
