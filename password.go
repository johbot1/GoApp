package main

import (
	"crypto/rand"
	"log"
	"math/big"
)

// generatePassword dynamically builds a password based on user input flags.
// It supports lowercase, uppercase, symbols, and word-based construction.
func generatePassword(length int, includeUppercase bool, includeSymbols bool, includeWords bool) string {
	// Edge case: prevent broken input (e.g., slider malfunction)
	if length <= 7 || length > 64 {
		log.Printf("[generatePassword]: Invalid password length provided: %d", length)
		wordListError = "Please specify a password length between 8 and 64."
		return "Please specify a password length between 8 and 64."
	}

	// Start with lowercase alphabet
	chars := "abcdefghijklmnopqrstuvwxyz"

	// Append optional character sets
	if includeUppercase {
		chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	if includeSymbols {
		chars += "!@#$%^&*()-_=+"
	}

	// A rune is a Go Type that supports unicode characters
	// This was the recommended to use over a standard array, due to how
	// Go handles unicode characters, and iterating over them
	var allChars []rune

	if includeWords {
		// Load the word list once if it's never been loaded
		loadWordListMutex.Do(func() {
			err := loadWordList("./static/words.json")
			if err != nil {
				log.Printf("[GeneratePassword]: Failed loading word list!")
				wordListError = "Failed to load word list!"
				wordListLoaded = false
				wordList = nil // Reset in case of partial load
			}
		})

		// Return error if loading failed
		if !wordListLoaded && wordListError != "" {
			return wordListError
		}

		// Convert words into characters and append to pool
		if wordListLoaded && len(wordList) > 0 {
			for _, word := range wordList {
				allChars = append(allChars, []rune(word)...)
			}
		}
	}

	// Add any additional characters gathered from toggles
	allChars = append(allChars, []rune(chars)...)

	// If no characters are available, return an error
	if len(allChars) == 0 {
		log.Printf("[generatePassword]: No character sets selected.")
		wordListError = "Please select at least one character set."
		return ""
	}

	// Warn user if password might lack randomness due to low character diversity
	if length > len(allChars) && includeWords && wordListLoaded && len(wordList) > 0 {
		log.Printf("[generatePassword]: Length exceeds available characters.")
		wordListError = "The requested length might result in a less random password..."
		return ""
	}

	// Password generation loop
	password := make([]rune, length)
	for i := 0; i < length; i++ {
		// Cryptographically secure random selection
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			log.Printf("[generatePassword]: Error during password generation.")
			wordListError = "Unexpected error while generating password."
			return ""
		}
		password[i] = allChars[randIndex.Int64()]
	}
	return string(password)
}
