writeJsonFile("gupm.json", {
    "dependencies": {
        "default": {
        }
    }
})

var err = exec("../build/dg", ["make"])

if(err) {
    throw new Error(JSON.stringify(err))
}