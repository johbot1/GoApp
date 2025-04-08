package main

import (
	"crypto/rand"
	"encoding/json"
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

// Function to load the word list from the JSON file
func loadWordList(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	var wordsMap map[string]int
	err = json.Unmarshal(data, &wordsMap)
	if err != nil {
		return err
	}

	// Extract the keys (the words) from the map
	for word := range wordsMap {
		wordList = append(wordList, word)
	}
	wordListLoaded = true
	return nil
}

func generatePassword(length int, includeUppercase bool, includeSymbols bool, includeWords bool) string {
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

	if includeWords {
		// Load the word list only if it hasn't been loaded yet
		loadWordListMutex.Do(func() {
			err := loadWordList("words.json")
			if err != nil {
				log.Printf("Error loading word list: %v", err)
				wordListLoaded = false
				// Clear any potentially partially loaded list
				wordList = nil
			}
		})

		if wordListLoaded && len(wordList) > 0 {
			// Add words to the pool of characters
			for _, word := range wordList {
				allChars = append(allChars, []rune(word)...)
			}
		}
	}
	allChars = append(allChars, []rune(chars)...)

	if len(allChars) == 0 {
		return "Please select at least one character set."
	}

	password := make([]rune, length)
	for i := 0; i < length; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		if err != nil {
			log.Println("Error generating random index:", err)
			return "Error generating password."
		}
		password[i] = allChars[randIndex.Int64()]
	}
	return string(password)
}

func main() {
	// Load the word list when the application starts
	err := loadWordList("./static/words.json")
	// If it can't find it, error out
	if err != nil {
		log.Fatalf("Error loading word list: %v", err)
	}

	http.HandleFunc("/", handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
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

		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		data := map[string]string{
			"Password": password,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	} else {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}
	}
}
