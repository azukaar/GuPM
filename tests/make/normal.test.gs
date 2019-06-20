writeJsonFile("gupm.json", {
    "dependencies": {
        "default": {
            "git://github.com/Masterminds/semver": "master"
        }
    }
})

var err = exec("../build/dg", ["make"])

if(err) {
    throw new Error(JSON.stringify(err))
}

if (!fileExists('gupm_modules/github.com/Masterminds/semver')) {
    throw new Error(1)
}