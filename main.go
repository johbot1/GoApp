package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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

// main starts the web server and handles static file routing
func main() {
	// Route for the homepage and password form handler
	http.HandleFunc("/", handler)
	log.Println("Listening on port 8080")

	// Serve static assets (JS, CSS, etc.)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))

}

// handler parses form submissions and injects data into the template
func handler(w http.ResponseWriter, r *http.Request) {
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
			return
		}

		// Extract and parse user selections
		lengthStr := r.FormValue("length")
		symbols := r.FormValue("symbols") == "true"
		words := r.FormValue("words") == "true"
		casePref := r.FormValue("case")         // Get case setting
		includeUppercase := casePref == "upper" // Convert to boolean

		// Validate length
		length, err := strconv.Atoi(lengthStr)
		if err != nil || length <= 0 {
			http.Error(w, "Invalid length", http.StatusBadRequest)
			return
		}

		// Generate the password and populate the template data
		password := generatePassword(length, includeUppercase, symbols, words)
		if casePref == "upper" && !words {
			password = strings.ToUpper(password)
		}
		data["Password"] = password

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
		return
	}
}
