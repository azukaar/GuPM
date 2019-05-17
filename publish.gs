removeFiles("gupm")
removeFiles("docs")
mkdir("docs")

var goArgs = ["build", "-o", "gupm/g"]
goArgs = goArgs.concat(dir("src/*.go"))
exec("go", goArgs)

copyFiles("plugins", "gupm/plugins")
copyFiles("gupm.json", "gupm/gupm.json")
copyFiles("install.sh", "docs/install.sh")

saveFileAt(tar("gupm"), "docs/gupm.tar.gz")

removeFiles("gupm")
