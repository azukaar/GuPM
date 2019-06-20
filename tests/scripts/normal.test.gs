writeFile("testscript.gs","writeFile('test','it works');")

var err = exec("../build/dg", ["testscript"])

if(err) {
    throw new Error(JSON.stringify(err))
}

if (!fileExists('test')) {
    throw new Error(1)
}