package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	// declare test cases
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  Hello  World  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "Joshua Serna & ",
			expected: []string{"joshua", "serna", "&"},
		},
		// add more cases here
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		if len(actual) != len(c.expected) {
			t.Errorf("Actual list length does not match Test list length\n")
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Actual word %s != %s\n", word, expectedWord)
			}
		}
	}
}
