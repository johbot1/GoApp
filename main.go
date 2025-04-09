package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"
)

// Define a global variable to hold the word list and a flag for loading status
var wordList []string
var wordListLoaded bool
var loadWordListMutex sync.Once // Ensures loadWordList is executed once AND ONLY ONCE!
var wordListError string        // Storing error messages

// Function to load the word list from the JSON file
func loadWordList(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading word list file '%s': %w", filename, err)
	}

	var wordsMap map[string]int
	err = json.Unmarshal(data, &wordsMap)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON from '%s': %w", filename, err)
	}

	// Extract the keys (the words) from the map
	for word := range wordsMap {
		wordList = append(wordList, word)
	}
	wordListLoaded = true
	return nil
}

func generatePassword(length int, includeUppercase bool, includeSymbols bool, includeWords bool) string {
	// Error Handling: Handles the case when the slider breaks like it had been.
	if length <= 7 || length > 64 {
		log.Printf("[generatePassword]: Invalid password length provided: %d", length)
		return "Please specify a password length between 8 and 64."
	}

	// Begins a list with lowercase Letters
	chars := "abcdefghijklmnopqrstuvwxyz"
	// Adds in uppercase Letters
	if includeUppercase {
		chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
	// Adds in symbols if selected
	if includeSymbols {
		chars += "!@#$%^&*()-_=+"
	}

	// A rune is a Go Type that supports unicode characters
	// This was the recommended to use over a standard array, due to how
	// Go handles unicode characters, and iterating over them
	var allChars []rune

	// Error Handling: If allChars is empty, and if the
	if includeWords {
		// Load the word list only if it hasn't been loaded yet
		loadWordListMutex.Do(func() {
			err := loadWordList("./static/words.json")
			if err != nil {
				log.Printf("[GeneratePassword]: Failed loading word list!")
				wordListError = "Failed to load word list!"
				wordListLoaded = false
				// Clear any potentially partially loaded list
				wordList = nil
			}
		})

		if !wordListLoaded && wordListError != "" {
			return wordListError
		}

		if wordListLoaded && len(wordList) > 0 {
			// Add words to the pool of characters
			for _, word := range wordList {
				allChars = append(allChars, []rune(word)...)
			}
		}
	}
	allChars = append(allChars, []rune(chars)...)

	// Error Handling: If the user includes words, but the length > the number of unique characters in allChars
	// the password would probably be less random.
	if len(allChars) == 0 {
		log.Printf("[generatePassword]: No character sets selected for password generation.")
		return "Please select at least one character set."
	}

	// Error Handling: In case the random integer function doesn't generate a random integer
	if length > len(allChars) && includeWords && wordListLoaded && len(wordList) > 0 {
		log.Printf("[generatePassword]: Requested length (%d) exceeds the available "+
			"unique characters when words are included. Consider a shorter length or fewer options.", length)
		return "The requested length might result in a less random password with the current options. " +
			"Consider adjusting the length or character sets."
	}

	password := make([]rune, length)
	for i := 0; i < length; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			log.Printf("[generatePassword]: Error generating random index during password creation.")
			return "Error generating password."
		}
		password[i] = allChars[randIndex.Int64()]
	}
	return string(password)
}

func main() {
	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string) // Create a map to hold data for the template

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		lengthStr := r.FormValue("length")
		uppercase := r.FormValue("uppercase") == "true"
		symbols := r.FormValue("symbols") == "true"
		words := r.FormValue("words") == "true"

		length, err := strconv.Atoi(lengthStr)
		if err != nil || length <= 0 {
			http.Error(w, "Invalid length", http.StatusBadRequest)
			return
		}

		password := generatePassword(length, uppercase, symbols, words)
		data["Password"] = password
		data["WordListError"] = wordListError // Pass the word list error to the template

	} else {
		data["WordListError"] = wordListError // Pass the word list error even on initial load (might be empty)
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}
