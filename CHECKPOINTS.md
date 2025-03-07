# Password Generator Web App - Checkpoints

This document outlines the checkpoints for building a dynamic password generator web application with enhanced features using Go.

## Phase 1: Project Setup and Basic Structure

- [x] **1.1** Set up Go development environment.
- [x] **1.2** Create project directory and `main.go` file.
- [x] **1.3** Create `templates` directory for HTML templates.
- [x] **1.4** Define basic project structure.

## Phase 2: Enhanced Password Generation Logic (Go)

- [x] **2.1** Implement `generatePassword(length int, includeUppercase bool, includeSymbols bool, includeWords bool)` function.
- [x] **2.2** Define character sets for lowercase, uppercase, symbols, and words (from a simple word list).
- [x] **2.3** Dynamically build the character set based on user selections.
- [x] **2.4** Handle potential errors during random number generation.
- [x] **2.5** Include `math/big` for proper random number generation.

## Phase 3: Web Server Setup (Go)

- [x] **3.1** Set up basic HTTP server using `net/http`.
- [x] **3.2** Define a handler function for the root URL (`/`).
- [x] **3.3** Implement basic routing for GET and POST requests.
- [x] **3.4** Start server on port 8080.

## Phase 4: Enhanced HTML Template Creation

- [x] **4.1** Create `index.html` template in the `templates` directory.
- [x] **4.2** Design a form with:
    - [x] A slider for password length.
    - [x] Checkboxes for uppercase, symbols, and words.
    - [x] A large text area to display the generated password.
    - [x] A "Copy to Clipboard" button.
- [x] **4.3** Use template actions (`{{if}}`, `{{end}}`, `{{.Variable}}`) to dynamically display the generated password.
- [x] **4.4** Implement basic HTML structure.
- [x] **4.5** Implement javascript for the copy to clipboard button.

## Phase 5: Integrating Go and HTML

- [x] **5.1** Parse HTML template using `html/template`.
- [x] **5.2** Handle form submission (POST request) in Go.
- [x] **5.3** Extract password length, uppercase, symbols, and words selections from form data.
- [x] **5.4** Call `generatePassword()` with the extracted data.
- [x] **5.5** Pass the generated password to the template for rendering.
- [x] **5.6** Render the template with the generated password.
- [x] **5.7** Handle GET requests for initial page load (no password).
- [x] **5.8** Handle invalid input lengths.

## Phase 6: Testing and Refinement

- [ ] **6.1** Test password generation with various lengths and combinations of options.
- [ ] **6.2** Verify secure randomness of generated passwords.
- [ ] **6.3** Test form submission and template rendering.
- [ ] **6.4** Test "Copy to Clipboard" functionality.
- [ ] **6.5** Improve error handling and user feedback.
- [ ] **6.6** Refactor code for readability and maintainability.

## Phase 7: Feedback
- To be completed upon project submission
