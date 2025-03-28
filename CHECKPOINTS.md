# Password Generator Web App - Checkpoints

This document outlines the checkpoints for building a dynamic password generator web application with enhanced features using Go.

## Phase 1: Project Setup and Basic Structure

- [ ] **1.1** Set up Go development environment.
- [ ] **1.2** Create project directory and `main.go` file.
- [ ] **1.3** Create `templates` directory for HTML templates.
- [ ] **1.4** Define basic project structure.

## Phase 2: Enhanced Password Generation Logic (Go)

- [ ] **2.1** Implement `generatePassword(length int, includeUppercase bool, includeSymbols bool, includeWords bool)` function.
- [ ] **2.2** Define character sets for lowercase, uppercase, symbols, and words (from a simple word list).
- [ ] **2.3** Dynamically build the character set based on user selections.
- [ ] **2.4** Handle potential errors during random number generation.
- [ ] **2.5** Include `math/big` for proper random number generation.

## Phase 3: Web Server Setup (Go)

- [ ] **3.1** Set up basic HTTP server using `net/http`.
- [ ] **3.2** Define a handler function for the root URL (`/`).
- [ ] **3.3** Implement basic routing for GET and POST requests.
- [ ] **3.4** Start server on port 8080.

## Phase 4: Enhanced HTML Template Creation

- [ ] **4.1** Create `index.html` template in the `templates` directory.
- [ ] **4.2** Design a form with:
    - [ ] A slider for password length.
    - [ ] Checkboxes for uppercase, symbols, and words.
    - [ ] A large text area to display the generated password.
    - [ ] A "Copy to Clipboard" button.
- [ ] **4.3** Use template actions (`{{if}}`, `{{end}}`, `{{.Variable}}`) to dynamically display the generated password.
- [ ] **4.4** Implement basic HTML structure.
- [ ] **4.5** Implement javascript for the copy to clipboard button.

## Phase 5: Integrating Go and HTML

- [ ] **5.1** Parse HTML template using `html/template`.
- [ ] **5.2** Handle form submission (POST request) in Go.
- [ ] **5.3** Extract password length, uppercase, symbols, and words selections from form data.
- [ ] **5.4** Call `generatePassword()` with the extracted data.
- [ ] **5.5** Pass the generated password to the template for rendering.
- [ ] **5.6** Render the template with the generated password.
- [ ] **5.7** Handle GET requests for initial page load (no password).
- [ ] **5.8** Handle invalid input lengths.

## Phase 6: Testing and Refinement

- [ ] **6.1** Test password generation with various lengths and combinations of options.
- [ ] **6.2** Verify secure randomness of generated passwords.
- [ ] **6.3** Test form submission and template rendering.
- [ ] **6.4** Test "Copy to Clipboard" functionality.
- [ ] **6.5** Improve error handling and user feedback.
- [ ] **6.6** Refactor code for readability and maintainability.

## Phase 7: Feedback
- To be completed upon project submission
