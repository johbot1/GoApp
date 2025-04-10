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

// Global state for the word list and error tracking
var wordList []string
var wordListLoaded bool
var loadWordListMutex sync.Once // Ensures loadWordList is executed once AND ONLY ONCE!
var wordListError string        // Storing error messages

// loadWordList reads a JSON file containing a map of words and extracts just the keys (the words).
// This function is only called once, even if triggered multiple times (guarded by sync.Once).
func loadWordList(filename string) error {
	data, err := ioutil.ReadFile(filename) // Read the file as a byte slice
	if err != nil {
		return fmt.Errorf("error reading word list file '%s': %w", filename, err)
	}

	var wordsMap map[string]int
	err = json.Unmarshal(data, &wordsMap) // Parse the JSON into a map
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON from '%s': %w", filename, err)
	}

	// Convert map keys to a slice of words
	for word := range wordsMap {
		wordList = append(wordList, word)
	}
	wordListLoaded = true
	return nil
}

// generatePassword dynamically builds a password based on user input flags.
// It supports lowercase, uppercase, symbols, and word-based construction.
func generatePassword(length int, includeUppercase bool, includeSymbols bool, includeWords bool) string {
	// Edge case: prevent broken input (e.g., slider malfunction)
	if length <= 7 || length > 64 {
		log.Printf("[generatePassword]: Invalid password length provided: %d", length)
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
		log.Printf("[generatePassword]: No character sets selected for password generation.")
		return "Please select at least one character set."
	}

	// Warn user if password might lack randomness due to low character diversity
	if length > len(allChars) && includeWords && wordListLoaded && len(wordList) > 0 {
		log.Printf("[generatePassword]: Requested length (%d) exceeds the available "+
			"unique characters when words are included. Consider a shorter length or fewer options.", length)
		return "The requested length might result in a less random password with the current options. " +
			"Consider adjusting the length or character sets."
	}

	// Password generation loop
	password := make([]rune, length)
	for i := 0; i < length; i++ {
		// Cryptographically secure random selection
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			log.Printf("[generatePassword]: Error generating random index during password creation.")
			return "Error generating password."
		}
		password[i] = allChars[randIndex.Int64()]
	}
	return string(password)
}

// main starts the web server and handles static file routing
func main() {
	// Route for the homepage and password form handler
	http.HandleFunc("/", handler)
	fmt.Println("Beginning web app!")

	// Serve static assets (JS, CSS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))

}

// handler parses form submissions and injects data into the template
func handler(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string) // Template data container

	if r.Method == http.MethodPost {
		// Parse the form submission
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Extract and parse user selections
		lengthStr := r.FormValue("length")
		uppercase := r.FormValue("uppercase") == "true"
		symbols := r.FormValue("symbols") == "true"
		words := r.FormValue("words") == "true"

		// Validate length
		length, err := strconv.Atoi(lengthStr)
		if err != nil || length <= 0 {
			http.Error(w, "Invalid length", http.StatusBadRequest)
			return
		}

		// Generate the password and populate the template data
		password := generatePassword(length, uppercase, symbols, words)
		data["Password"] = password
		data["WordListError"] = wordListError // Show any word list errors in UI
		data["Length"] = lengthStr
		data["Uppercase"] = strconv.FormatBool(uppercase)
		data["Symbols"] = strconv.FormatBool(symbols)
		data["Words"] = strconv.FormatBool(words)
	} else {
		// On initial page load or GET, pass in error if one exists
		data["WordListError"] = wordListError
	}

	// Parse and render the HTML template
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}
