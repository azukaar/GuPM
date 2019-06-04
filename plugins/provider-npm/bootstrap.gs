console.log("Warning: you used npm as a provider meaning GuPM is going to use package.json instead of gupm.json which is not recommmended.")
console.log("Press any key if that's what you want or CTRL+C to cancel.")
waitForKey()

if(fileExists("package.json") || fileExists("gupm.json")) {
    console.error("A project already exist in this project. Aborting.")
    exit()
}

var name = waitForInput("Please enter the name of the project: ")
var description = waitForInput("Enter a description: ")
var author = waitForInput("Enter the author: ")
var licence = waitForInput("Enter the licence (ISC): ")

if(name == "") {
    console.error("Name cannot be empty. Try again.")
    exit()
}

var result = {
    name: name,
    description: description,
    author: author,
    licence: licence || "ISC"
}


writeJsonFile("package.json", result)