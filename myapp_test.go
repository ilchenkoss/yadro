package main

import (
	"testing"
)

func TestExample1(t *testing.T) {
	notNormalizedString := "follower brings bunch of questions!"
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
