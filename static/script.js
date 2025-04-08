// Javascript Flow
// 1. Get references to the length slider, its display, and the password display area.
// 2. Add an input listener to the length slider to dynamically update the displayed value.
// 3. Define a function to copy the generated password to the clipboard.
// 4. Inside the copy function, get the textarea element containing the password.
// 5. If the textarea exists, select its content for copying.
// 6. Execute the copy command to place the selected text in the clipboard.
// 7. Show a visual alert to the user confirming the copy action.
const lengthSlider = document.getElementById("length");
const lengthValueSpan = document.getElementById("length-value");
const passwordDisplay = document.getElementById("password-display");

// Ensure the length slider and its value display element exist before adding the listener.
if (lengthSlider && lengthValueSpan) {
    lengthSlider.oninput = function () {
        lengthValueSpan.textContent = this.value; // Update the displayed length as the slider moves.
    };
}

// Function to copy the text content of the password display textarea to the clipboard.
function copyPassword() {
    // Get the text area where the generated password is displayed.
    const textarea = passwordDisplay ? passwordDisplay.querySelector("textarea") : null;
    // Check if the text area exists before attempting to copy.
    if (textarea) {
        textarea.select(); // Select the entire text content within the textarea.
        document.execCommand("copy"); // Execute the browser's copy command.
        alert("Password copied to clipboard!"); // Inform the user that the password has been copied.
    }
}

if (passwordForm) {
    // Event listener for the submission event
    passwordForm.addEventListener('submit', function (event) {
        // Prevents the default page reload after hitting "generate Password"
        event.preventDefault();

        // Creates a formData object
        const formData = new FormData(passwordForm);

        // Sends a POST request to the app
        fetch('/', {
            method: 'POST',
            body: formData,
        })
            // The app will send the password as plain text as its response
            .then(response => response.text())
            // When successful, it updates the passwordDisplay div with the response and the clipboard button
            .then(data => {
                // Update the password display area with the received data
                passwordDisplay.innerHTML = `<textarea rows="5" readonly>${data}</textarea>
                                         <button type="button" id="copy-button" onclick="copyPassword()">Copy to Clipboard</button>`;
            })
            //If there is anything that goes wrong, log it in the console and display it in the passwordDisplay div
            .catch(error => {
                console.error('Error generating password:', error);
                passwordDisplay.innerHTML = `<div class="error-message">Error generating password. Please try again.</div>`;
            });
    });
}