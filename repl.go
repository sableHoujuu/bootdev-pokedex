package main

import "strings"

// Handles everything to do with the REPL specific functions of the program

func cleanInput(text string) []string {
	if text == "" {
		return nil
	}
	text = strings.TrimSpace(text)
	var words []string
	for part := range strings.SplitSeq(text, " ") {
		if part == " " || part == "" {
			continue
		}
		words = append(words, strings.ToLower(part))
	}
	return words
}
