writeJsonFile("gupm.json", {
    "dependencies": {
        "default": {
            "npm://webpack": "latest"
        }
    }
})

var err = exec("../build/dg", ["pl", "install", "https://azukaar.github.io/GuPM-official/repo:provider-npm"])
var err2 = exec("../build/dg", ["make"])

if(err || err2) {
    throw new Error(JSON.stringify(err))
}

if (!fileExists('node_modules/webpack')) {
    throw new Error(1)
}