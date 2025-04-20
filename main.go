package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Global state for the word list and error tracking
var wordList []string
var wordListLoaded bool
var loadWordListMutex sync.Once // Ensures loadWordList is executed only once
var wordListError string        // Stores error messages

// loadWordList reads a JSON file containing a list of words.
// It populates the global wordList slice and is protected by sync.Once.
func loadWordList(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("ERROR [loadWordList] Failed to read file '%s': %v", filename, err)
		return fmt.Errorf("error reading word list file '%s': %w", filename, err)
	}

	err = json.Unmarshal(data, &wordList)
	if err != nil {
		log.Printf("ERROR [loadWordList] Failed to parse JSON array from '%s': %v", filename, err)
		return fmt.Errorf("error unmarshalling word list from '%s': %w", filename, err)
	}

	if len(wordList) == 0 {
		wordListLoaded = false
		log.Printf("Warning: Word list '%s' loaded but is empty.", filename)
		return fmt.Errorf("word list '%s' is empty", filename)
	}

	wordListLoaded = true
	log.Printf("Word list '%s' loaded successfully with %d words.", filename, len(wordList))
	return nil
}

// main starts the web server and configures routes and static file serving
func main() {
	http.HandleFunc("/", passwordFormHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// passwordFormHandler parses form submissions, generates a password, and renders the template
func passwordFormHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[passwordFormHandler] Request received: method=%s", r.Method)
	defer func() { wordListError = "" }() // Clear global error state after each request

	data := make(map[string]string)

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			log.Println("ERROR [passwordFormHandler] Failed to parse form")
			return
		}

		// Extract form values
		lengthStr := r.FormValue("length")
		symbols := r.FormValue("symbols") == "true"
		words := r.FormValue("words") == "true"
		casePref := r.FormValue("case")

		log.Printf("[passwordFormHandler] Form values - length=%s | case=%s | symbols=%t | words=%t",
			lengthStr, casePref, symbols, words)

		// Convert length to integer
		length, err := strconv.Atoi(lengthStr)
		if err != nil || length <= 0 {
			http.Error(w, "Invalid length", http.StatusBadRequest)
			log.Println("ERROR [passwordFormHandler] Invalid length provided")
			return
		}

		// Generate the password using user preferences
		log.Println("[passwordFormHandler] Calling generatePassword with parsed options")
		password := generatePassword(length, casePref, symbols, words)

		data["Password"] = password
		data["WordListError"] = wordListError
		data["Length"] = lengthStr
		data["Case"] = casePref
		data["Symbols"] = strconv.FormatBool(symbols)
		data["Words"] = strconv.FormatBool(words)

		if password != "" {
			log.Println("[passwordFormHandler] Password generated successfully")
		} else {
			log.Println("ERROR [passwordFormHandler] Password generation failed or returned empty result")
		}
	} else {
		// Handle initial page load (GET)
		data["WordListError"] = wordListError
		data["Length"] = "8"
		data["Case"] = "lower"
	}

	// Load and render HTML template
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		log.Println("ERROR [passwordFormHandler] Failed to execute template")
	}
}
