removeFiles(["build/dg", "build/plugins", "build/gupm.json"]) 

var goArgs = ["build", "-o", "build/dg"]
goArgs = goArgs.concat(dir("src/*.go"))
exec("go", goArgs)

copyFiles("plugins", "build/plugins")
copyFiles("gupm.json", "build/gupm.json")
