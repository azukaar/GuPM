mkdir('node_modules/.bin')

function installDir(files, parent) {
    if(!parent) parent = '';

    for(f in files) {
        var dirName = files[f];
        if(dirName.match(/^@/)) {
            var subfiles = readDir('node_modules/' + parent + dirName)
            installDir(subfiles, dirName + '/')
        }
        else {
            // console.log('node_modules/' + parent + dirName + '/package.json')
            var package = readJsonFile('node_modules/' + parent + dirName + '/package.json')
            if(package.bin) {
                if(typeof package.bin == 'string') {
                    var relPath = '../' + parent + dirName + '/' +package.bin.replace(/^\.\//, '')
                    createSymLink(relPath, 'node_modules/.bin/' + package.name)
                } else {
                    for(b in package.bin) {
                        var bin = package.bin[b]
                        var relPath = '../' + parent + dirName + '/' +bin.replace(/^\.\//, '')
                        createSymLink(relPath, 'node_modules/.bin/' + b)
                    }
                }
            }
        }
    }
}

var files = readDir(Source)
installDir(files)