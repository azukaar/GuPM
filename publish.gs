removeFiles("gupm")
removeFiles("docs/install.sh")

var goArgs = ["build", "-o"]

if(typeof $1 != "undefined" && $1 == "windows") {
    goArgs.push("gupm/g.exe")
} else {
    goArgs.push("gupm/g")
}

goArgs = goArgs.concat(dir("src/*.go"))

copyFiles("plugins", "gupm/plugins")
copyFiles("gupm.json", "gupm/gupm.json")
copyFiles("install.sh", "docs/install.sh")

var arch = tar("gupm")

if(typeof $1 != "undefined" && $1 == "mac") {
    env("GOOS", "darwin")
    env("go version", "amd64")
    exec("go", goArgs)
    removeFiles("docs/gupm_mac.tar.gz")
    saveFileAt(arch, "docs/gupm_mac.tar.gz")
} 
if(typeof $1 != "undefined" && $1 == "windows") {
    env("GOOS", "windows")
    env("GOARCH", "amd64")
    exec("go", goArgs)
    removeFiles("docs/gupm_windows.tar.gz")
    saveFileAt(arch, "docs/gupm_windows.tar.gz")
} else {
    env("GOOS", "linux")
    env("GOARCH", "amd64")
    exec("go", goArgs)
    removeFiles("docs/gupm.tar.gz")
    saveFileAt(arch, "docs/gupm.tar.gz")
}

removeFiles("gupm")