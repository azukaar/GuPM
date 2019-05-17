removeFiles(["build/dg", "build/plugins", "build/gupm.json"]) 
exec("go", ["build", "-o", "build/dg", "src/index.go", "src/addDependency.go", "src/installProject.go"])
copyFiles("plugins", "build/plugins")
copyFiles("gupm.json", "build/gupm.json")
