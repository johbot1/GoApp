package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	data, err := os.ReadFile(filename) // Read the file as a byte slice
	if err != nil {
		log.Printf("ERROR [loadWordList] Failed to read file '%s': %v", filename, err)
		return fmt.Errorf("error reading word list file '%s': %w", filename, err)
	}
	//Unmarshal directly into the global wordList slice
	err = json.Unmarshal(data, &wordList)
	if err != nil {
		log.Printf("ERROR [loadWordList] Failed to parse JSON array from '%s': %v", filename, err)
		return fmt.Errorf("error unmarshalling word list from '%s': %w", filename, err)
	}

	// Check if the list is empty after successful unmarshalling
	if len(wordList) == 0 {
		wordListLoaded = false // Treat empty list as not successfully loaded for practical purposes
		log.Printf("Warning: Word list '%s' loaded but is empty.", filename)
		return fmt.Errorf("word list '%s' is empty", filename)
	} else {
		wordListLoaded = true
		log.Printf("Word list '%s' loaded successfully with %d words.", filename, len(wordList))
	}
	return nil
}

// main starts the web server and handles static file routing
func main() {
	// Route for the homepage and password form passwordFormHandler
	http.HandleFunc("/", passwordFormHandler)
	log.Println("Listening on port 8080")

	// Serve static assets (JS, CSS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))

}

// passwordFormHandler parses form submissions and injects data into the template
func passwordFormHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[passwordFormHandler] Request received: method=%s", r.Method)
	defer func() {
		wordListError = "" // Clear after handling this request
	}()

	// Template data container
	data := make(map[string]string)

	if r.Method == http.MethodPost {
		// Parse the form submission
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			log.Println("ERROR [passwordFormHandler] Password generation failed or returned empty result")
			return
		}

		// Extract and parse user selections
		lengthStr := r.FormValue("length")
		symbols := r.FormValue("symbols") == "true"
		words := r.FormValue("words") == "true"
		casePref := r.FormValue("case")         // Get case setting
		includeUppercase := casePref == "upper" // Convert to boolean

		log.Printf("[passwordFormHandler] Form values - length=%s | case=%s | symbols=%t | words=%t",
			lengthStr, r.FormValue("case"), r.FormValue("symbols") == "true", r.FormValue("words") == "true")

		// Validate length
		length, err := strconv.Atoi(lengthStr)
		if err != nil || length <= 0 {
			http.Error(w, "Invalid length", http.StatusBadRequest)
			log.Println("ERROR [passwordFormHandler] Password generation failed or returned empty result")
			return
		}

		// Generate the password and populate the template data
		log.Println("[passwordFormHandler] Calling generatePassword with parsed options")
		password := generatePassword(length, includeUppercase, symbols, words)
		if casePref == "upper" && !words {
			password = strings.ToUpper(password)
		}
		data["Password"] = password
		if password != "" {
			log.Printf("[passwordFormHandler] Password generated successfully!")
		} else {
			log.Println("ERROR [passwordFormHandler] Password generation failed or returned empty result")
		}

		// Always pass the latest error (it may have changed in generatePassword)
		data["WordListError"] = wordListError
		data["Length"] = lengthStr
		data["Case"] = casePref
		data["Symbols"] = strconv.FormatBool(symbols)
		data["Words"] = strconv.FormatBool(words)
	} else {
		// On initial page load or GET, pass in error if one exists
		data["WordListError"] = wordListError
		data["Length"] = "8"
		data["Case"] = "lower"
	}

	// Parse and render the HTML template
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		log.Println("ERROR [passwordFormHandler] Password generation failed or returned empty result")
		return
	}
}
