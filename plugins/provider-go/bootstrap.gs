if(fileExists("gupm.json")) {
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


writeJsonFile("gupm.json", result)
writeFile(".gupm_rc.gs", 'env("GOPATH", run("go", ["env", "GOROOT"]) + ":" + pwd() + "/go_modules")')