# Password Generator Web App - Checkpoints

This document outlines the checkpoints for building a dynamic password generator web application with enhanced features using Go.

## Phase 1: Project Setup and Basic Structure

- [x] **1.1** Set up Go development environment.
- [x] **1.2** Create project directory and `main.go` file.
- [x] **1.3** Create `templates` directory for HTML templates.
- [x] **1.4** Define basic project structure.

## Phase 2: Enhanced Password Generation Logic (Go)

- [X] **2.1** Implement `generatePassword(length int, includeUppercase bool, includeSymbols bool, includeWords bool)` function.
- [X] **2.2** Define character sets for lowercase, uppercase, symbols, and words (from a simple word list).
- [X] **2.3** Dynamically build the character set based on user selections.
- [X] **2.4** Handle potential errors during random number generation.
- [X] **2.5** Include `math/big` for proper random number generation.

## Phase 3: Web Server Setup (Go)

- [X] **3.1** Set up basic HTTP server using `net/http`.
- [X] **3.2** Define a handler function for the root URL (`/`).
- [X] **3.3** Implement basic routing for GET and POST requests.
- [X] **3.4** Start server on port 8080.

## Phase 4: Enhanced HTML Template Creation
- [X] **4.1** Create `index.html` template in the `templates` directory.
- [X] **4.2** Design a form with:
    - [X] A slider for password length.
    - [X] Checkboxes for uppercase, symbols, and words.
    - [X] A large text area to display the generated password.
    - [X] A "Copy to Clipboard" button.
- [X] **4.3** Use template actions (`{{if.Password}}` and `{{.Password}}`) to dynamically display the generated password.
- [X] **4.4** Implement basic HTML structure.
- [X] **4.5** Implement javascript for the copy to clipboard button.

## Phase 5: Integrating Go and HTML

- [X] **5.1** Parse HTML template using `html/template`.
- [X] **5.2** Handle form submission (POST request) in Go.
- [X] **5.3** Extract password length, uppercase, symbols, and words selections from form data.
- [X] **5.4** Call `generatePassword()` with the extracted data.
- [X] **5.5** Pass the generated password to the template for rendering.
- [X] **5.6** Render the template with the generated password.
- [X] **5.7** Handle GET requests for initial page load (no password).
- [X] **5.8** Handle invalid input lengths.

## Phase 6: Testing and Refinement

- [X] **6.1** Test password generation with various lengths and combinations of options.
- [X] **6.2** Verify secure randomness of generated passwords.
- [X] **6.3** Test form submission and template rendering.
- [X] **6.4** Test "Copy to Clipboard" functionality.
- [X] **6.5** Improve error handling and user feedback.
- [X] **6.6** Refactor code for readability and maintainability.

## Phase 7: Feedback
- [X] Start at 8
- [X] Change "Include Uppercase" to be default 
  - [X] Add "Only Upper Case" / "only lowercase" instead
- [X] Change include words to ONLY words instead
  - [X] This means that it will be ONLY words squished together
- [] Add minor debug statements
- [] Change UPPERCASE to constants, same w/ Symbols and lowercase 
- [] Rename Handler method to more descriptive name
