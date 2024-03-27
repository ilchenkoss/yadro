package main

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	notNormalizedString := "i'll follow you as long as you are following me"
	expected := "follow long"
	actual := stringNormalization(notNormalizedString)

	if expected != actual {
		t.Errorf("Result was incorrect, got: %s, want: %s.", actual, expected)
	} else {
		fmt.Printf("Good: %s, want: %s.", actual, expected)
	}
}
