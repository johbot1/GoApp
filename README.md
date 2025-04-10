# GoApp
One of three Go! related projects for my final semester, this one a password generator web app!
I understand the point of making students use a tech they are unfamiliar with, however Go has made
this seemingly simple project a nightmare, due to how it's initialized and how it is setup to run in
IntelliJ. After a few rough starts, I was able to finally begin work, which while challenging, was not
_too_ bad. Figuring out Go was convoluted and not very fun. I am glad to be done with the language.

# Flow:
1) When teh server starts, wordListError is empty
2) If the wordbank can't be loaded, the loadWordList will set an error with an informative message
3) When the user submits the form with "Include Words" checked, generatePassword is called
4) So long as wordListError is empty, the password is returned. If it isn't empty, the password returned is the error message
5) The handler will pass this message to the template
6) The index.html will check if WordListError exists, and display it above password display area


# Sources:
- https://go.dev/doc/articles/wiki/
- https://medium.com/@harshvardhancc/building-web-apps-with-golang-a-step-by-step-guide-3754a25dcf47
- https://www.sohamkamani.com/golang/how-to-build-a-web-application/
- https://pkg.go.dev/sync