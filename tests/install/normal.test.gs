writeJsonFile("gupm.json", {
    "dependencies": {
        "default": {
        }
    }
})


var err = exec("../build/dg", ["i", "git://github.com/Masterminds/semver"])

if(err) {
    throw new Error(JSON.stringify(err))
}

if (!fileExists('gupm_modules/github.com/Masterminds/semver')) {
    throw new Error(1)
}