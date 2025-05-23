package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"
)

const (
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	symbolChars    = "!@#$%^&*()-_=+"
)

// generatePassword constructs a password of exact 'length'.
// If includeWords is true, it assembles the password from whole words whose combined length matches the target.
// Symbols are inserted at the beginning, middle, and end if enabled. Letter casing is applied after construction.
func generatePassword(length int, casePref string, includeSymbols bool, includeWords bool) string {
	// Validate input bounds (had an issue with it flying out of bounds, so this will ensure it doesn't)
	if length <= 7 || length > 64 {
		log.Printf("[generatePassword]: Invalid password length: %d", length)
		wordListError = "Please specify a password length between 8 and 64."
		return "Please specify a password length between 8 and 64."
	}

	// Word-based password construction
	if includeWords {
		// Ensure the word list is loaded only once
		loadWordListMutex.Do(func() {
			err := loadWordList("./static/words.json")
			if err != nil {
				log.Printf("[generatePassword]: Failed loading word list")
				wordListError = "Failed to load word list"
				wordListLoaded = false
				wordList = nil
			}
		})

		if !wordListLoaded || len(wordList) == 0 {
			wordListError = "Word list unavailable"
			return ""
		}

		// Group words by their length for efficient lookup
		lengthMap := map[int][]string{}
		minLen := 100
		for _, word := range wordList {
			l := len(word)
			if l > length {
				continue
			}
			lengthMap[l] = append(lengthMap[l], word)
			if l < minLen {
				minLen = l
			}
		}

		var result []string

		// Attempt to find a valid combination of words whose lengths sum to the requested password length
		if findWordCombo(length, []int{}, &result, lengthMap) {
			password := strings.Join(result, "")

			// Apply casing if enabled
			if casePref == "upper" {
				password = strings.ToUpper(password)
			} else if casePref == "mixed" {
				password = applyMixedCase(password)
			}

			// Apply symbols if enabled
			if includeSymbols && len(password) >= 3 {
				passwordRunes := []rune(password)

				startIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(symbolChars))))
				passwordRunes[0] = rune(symbolChars[startIdx.Int64()])

				endIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(symbolChars))))
				passwordRunes[len(passwordRunes)-1] = rune(symbolChars[endIdx.Int64()])

				midPos := len(passwordRunes) / 2
				midIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(symbolChars))))
				passwordRunes[midPos] = rune(symbolChars[midIdx.Int64()])

				password = string(passwordRunes)
			}

			return password
		}

		wordListError = "No combination of words could be found for the selected length"
		return ""
	}

	// Standard character-based password construction
	chars := lowercaseChars
	if casePref == "upper" {
		chars = uppercaseChars
	}
	if includeSymbols {
		chars += symbolChars
	}

	allChars := []rune(chars)

	if len(allChars) == 0 {
		wordListError = "Please select at least one character set"
		return ""
	}

	password := make([]rune, length)
	for i := 0; i < length; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			wordListError = "Unexpected error while generating password"
			return ""
		}
		password[i] = allChars[randIndex.Int64()]
	}

	return string(password)
}

// findWordCombo recursively searches for a combination of word lengths that exactly sum to the target password length.
// When a valid combination is found, it appends randomly selected words matching those lengths to the result.
func findWordCombo(remaining int, path []int, result *[]string, lengthMap map[int][]string) bool {
	if remaining == 0 {
		for _, l := range path {
			words := lengthMap[l]
			if len(words) == 0 {
				return false
			}
			rIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
			if err != nil {
				return false
			}
			*result = append(*result, words[rIdx.Int64()])
		}
		return true
	}

	for l := range lengthMap {
		if l <= remaining {
			if findWordCombo(remaining-l, append(path, l), result, lengthMap) {
				return true
			}
		}
	}
	return false
}
