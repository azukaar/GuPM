removeFiles("gupm")
removeFiles("docs/install.sh")

var goArgs = ["build", "-o", "gupm/g"]
goArgs = goArgs.concat(dir("src/*.go"))
exec("go", goArgs)

copyFiles("plugins", "gupm/plugins")
copyFiles("gupm.json", "gupm/gupm.json")
copyFiles("install.sh", "docs/install.sh")

var arch = tar("gupm")

if(typeof $1 != "undefined" && $1 == "mac") {
    removeFiles("docs/gupm_mac.tar.gz")
    saveFileAt(arch, "docs/gupm_mac.tar.gz")
} else {
    removeFiles("docs/gupm.tar.gz")
    saveFileAt(arch, "docs/gupm.tar.gz")
}

removeFiles("gupm")